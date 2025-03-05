package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram/bot/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	m "github.com/sejo412/ya-boo/pkg/models"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Open(dsn string) error {
	var lastErr error
	for attempt := 0; attempt < RetryMaxRetries; attempt++ {
		db, err := sql.Open(driver, dsn)
		if err == nil {
			p.DB = db
			return nil
		}
		lastErr = err
		delay := RetryInitDelay + time.Duration(attempt)*RetryDeltaDelay
		time.Sleep(delay)
		continue
	}
	return fmt.Errorf("failed to open postgres connection: %w", lastErr)
}

func (p *Postgres) Close() {
	_ = p.DB.Close()
}

func (p *Postgres) Ping() error {
	return p.DB.Ping()
}

func (p *Postgres) IsAdminsInitialized(ctx context.Context) (bool, error) {
	return isTrue(p.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE role=$1", m.RoleAdmin))
}

func (p *Postgres) IsUserPresent(ctx context.Context, id int64) (bool, error) {
	return isTrue(p.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", id))
}

func (p *Postgres) IsAdmin(ctx context.Context, id int64) (bool, error) {
	return isTrue(p.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1 AND role = $2", id, m.RoleAdmin))
}

func (p *Postgres) IsWaitingApprove(ctx context.Context, id int64) (bool, error) {
	return p.IsUserInRole(ctx, id, m.RoleUnknown)
}

func (p *Postgres) IsUserInRole(ctx context.Context, id int64, role m.Role) (bool, error) {
	return isTrue(p.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1 AND role = $2", id, role))
}

func (p *Postgres) IsRegisteredUser(ctx context.Context, id int64) (bool, error) {
	return isTrue(p.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1 AND role != $2", id, m.RoleUnknown))
}

func isTrue(row *sql.Row) (bool, error) {
	var res any
	err := row.Scan(&res)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (p *Postgres) UpsertUser(ctx context.Context, user m.User) error {
	updated := time.Now()
	query := `
		INSERT INTO users (id, first_name, last_name, username, role)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(id) DO UPDATE
		SET first_name = $2, last_name = $3, username = $4, role = $5, updated = $6, llm = $7`
	_, err := p.DB.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName, user.Username, user.Role, updated,
		user.LLM.Id)
	if err != nil {
		log.Printf("failed to upsert user: %v", err)
		return err
	}
	return nil
}

func (p *Postgres) UpdateUserRole(ctx context.Context, user m.User, role m.Role) error {
	updated := time.Now()
	query := "UPDATE users SET role = $1, updated = $2 WHERE id = $3"
	_, err := p.DB.ExecContext(ctx, query, role, updated, user.ID)
	if err != nil {
		log.Printf("failed to update user role: %v", err)
		return err
	}
	return nil
}

func (p *Postgres) ListUsers(ctx context.Context) ([]m.User, error) {
	result := make([]m.User, 0)
	rows, err := p.DB.QueryContext(ctx,
		"SELECT id, first_name, last_name, username, role FROM users ORDER BY role DESC")
	if err != nil {
		log.Printf("failed to make query: %v", err)
		return nil, fmt.Errorf("failed to make query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var role m.Role
		var id int64
		var firstName, lastName, username string
		if err = rows.Scan(&id, &firstName, &lastName, &username, &role); err != nil {
			log.Printf("failed to scan users: %v", err)
			return nil, fmt.Errorf("failed to scan users: %w", err)
		}
		result = append(result, m.User{
			User: &models.User{
				ID:        id,
				FirstName: firstName,
				LastName:  lastName,
				Username:  username,
			},
			Role: role,
		})
	}
	if err = rows.Err(); err != nil {
		log.Printf("failed to read rows: %v", err)
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}
	return result, nil
}

func (p *Postgres) GetLLMs(ctx context.Context) ([]m.LLM, error) {
	result := make([]m.LLM, 0)
	rows, err := p.DB.QueryContext(ctx,
		"SELECT * FROM llm ORDER BY id DESC")
	if err != nil {
		log.Printf("failed to make query: %v", err)
		return nil, fmt.Errorf("failed to make query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var id int64
		var name, endpoint, token, description string
		if err = rows.Scan(&id, &name, &endpoint, &token, &description); err != nil {
			log.Printf("failed to scan llms: %v", err)
			return nil, fmt.Errorf("failed to scan llms: %w", err)
		}
		result = append(result, m.LLM{
			Id:          id,
			Name:        name,
			Endpoint:    endpoint,
			Token:       token,
			Description: description,
		})
	}
	if err = rows.Err(); err != nil {
		log.Printf("failed to read rows: %v", err)
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}
	return result, nil
}

func (p *Postgres) GetUserLLM(ctx context.Context, userId int64) (m.LLM, error) {
	var id int64
	var name, endpoint, description string
	query := `
		SELECT llm.id, llm.name, llm.endpoint, llm.description FROM llm 
			JOIN users ON llm.id = users.llm
		WHERE users.id = $1`
	if err := p.DB.QueryRowContext(ctx, query, userId).Scan(&id, &name, &endpoint, &description); err != nil {
		log.Printf("failed to make query: %v", err)
		return m.LLM{}, fmt.Errorf("failed to make query: %w", err)
	}
	return m.LLM{
		Id:          id,
		Name:        name,
		Endpoint:    endpoint,
		Description: description,
	}, nil
}

func (p *Postgres) AddLLM(ctx context.Context, llm m.LLM) error {
	query := `
		INSERT INTO llm
		(name, endpoint, token, description)
		VALUES ($1, $2, $3, $4) ON CONFLICT (name) DO NOTHING`
	_, err := p.DB.ExecContext(ctx, query, llm.Name, llm.Endpoint, llm.Token, llm.Description)
	if err != nil {
		log.Printf("failed add llm: %v", err)
		return err
	}
	return nil
}

func (p *Postgres) RemoveLLM(ctx context.Context, llm m.LLM) error {
	query := "DELETE FROM llm WHERE id = $1"
	_, err := p.DB.ExecContext(ctx, query, llm.Id)
	if err != nil {
		log.Printf("failed remove llm: %v", err)
		return err
	}
	return nil
}

func (p *Postgres) SetUserLLM(ctx context.Context, userId int64, llm m.LLM) error {
	query := "UPDATE users SET llm = $1, updated = $2 WHERE id = $3"
	_, err := p.DB.ExecContext(ctx, query, llm.Id, time.Now(), userId)
	if err != nil {
		log.Printf("failed update user llm: %v", err)
		return err
	}
	return nil
}
