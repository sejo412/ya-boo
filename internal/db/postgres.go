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
	query := `
		INSERT INTO users (id, first_name, last_name, username, role)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(id) DO UPDATE
		SET first_name = $2, last_name = $3, username = $4, role = $5`
	_, err := p.DB.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName, user.Username, user.Role)
	if err != nil {
		log.Printf("failed to upsert user: %v", err)
		return err
	}
	return nil
}

func (p *Postgres) UpdateUserRole(ctx context.Context, user m.User, role m.Role) error {
	query := "UPDATE users SET role = $1 WHERE id = $2"
	_, err := p.DB.ExecContext(ctx, query, role, user.ID)
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
