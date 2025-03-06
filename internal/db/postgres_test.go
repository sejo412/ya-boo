//go:build integration

package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/go-telegram/bot/models"
	m "github.com/sejo412/ya-boo/pkg/models"
)

var TestDB = NewPostgres()

func TestMain(m *testing.M) {
	fmt.Println("open db")
	_ = TestDB.Open("postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable")
	defer TestDB.Close()
	exitVal := m.Run()
	_, _ = TestDB.DB.Exec("DELETE FROM users WHERE id < 0")
	_, _ = TestDB.DB.Exec("DELETE FROM llm WHERE id < 0")
	os.Exit(exitVal)
}

func TestPostgres_UpsertUser(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		user m.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Upsert admin user success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				user: m.User{
					User: &models.User{
						ID:       -1,
						Username: "TestAdminUser",
					},
					Role: m.RoleAdmin,
				},
			},
			wantErr: false,
		},
		{
			name: "Upsert regular user success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				user: m.User{
					User: &models.User{
						ID:       -2,
						Username: "TestRegularUser",
					},
					Role: m.RoleRegularUser,
				},
			},
			wantErr: false,
		},
		{
			name: "Upsert unknown user success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				user: m.User{
					User: &models.User{
						ID:       -3,
						Username: "TestUnknownUser",
					},
					Role: m.RoleUnknown,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			if err := p.UpsertUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UpsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgres_ListUsers(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "List Users success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			got, err := p.ListUsers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var adminUser, regularUser, unknownUser = false, false, false
			for _, res := range got {
				if res.ID == -1 {
					adminUser = true
				}
				if res.ID == -2 {
					regularUser = true
				}
				if res.ID == -3 {
					unknownUser = true
				}
			}
			if !adminUser && !regularUser && !unknownUser {
				t.Errorf("ListUsers() error")
			}
		})
	}
}

func TestPostgres_AddLLM(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
		llm m.LLM
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Add LLM success #1",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				llm: m.LLM{
					ID:   -1,
					Name: "test #1",
				},
			},
			wantErr: false,
		},
		{
			name: "Add LLM success #2",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				llm: m.LLM{
					ID:   -2,
					Name: "test #2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			if err := p.AddLLM(tt.args.ctx, tt.args.llm); (err != nil) != tt.wantErr {
				t.Errorf("AddLLM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgres_GetLLMs(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Get LLMs success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			got, err := p.GetLLMs(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLLMs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var testLLM1, testLLM2 = false, false
			for _, res := range got {
				if res.ID == -1 {
					testLLM1 = true
				}
				if res.ID == -2 {
					testLLM2 = true
				}
			}
			if !testLLM1 && !testLLM2 {
				t.Errorf("GetLLMs() error")
			}
		})
	}
}

func TestPostgres_SetUserLLM(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int64
		llm    m.LLM
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Set User LLM success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx:    context.Background(),
				userID: -1,
				llm: m.LLM{
					ID: -1,
				},
			},
			wantErr: false,
		},
		{
			name: "Set User LLM error",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx:    context.Background(),
				userID: -2,
				llm: m.LLM{
					ID: -100,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			if err := p.SetUserLLM(tt.args.ctx, tt.args.userID, tt.args.llm); (err != nil) != tt.wantErr {
				t.Errorf("SetUserLLM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgres_GetUserLLM(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    m.LLM
		wantErr bool
	}{
		{
			name: "Get user LLM success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx:    context.Background(),
				userID: -1,
			},
			want: m.LLM{
				ID:          -1,
				Name:        "test #1",
				Endpoint:    "",
				Description: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			got, err := p.GetUserLLM(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserLLM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserLLM() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgres_IsAdmin(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "IsAdmin true",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				id:  -1,
			},
			want: true,
		},
		{
			name: "IsAdmin false",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				id:  -2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			if got := p.IsAdmin(tt.args.ctx, tt.args.id); got != tt.want {
				t.Errorf("IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgres_IsAdminsInitialized(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "IsAdminsInitialized true",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			got, err := p.IsAdminsInitialized(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsAdminsInitialized() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsAdminsInitialized() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgres_IsRegisteredUser(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "IsRegisteredUser true",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				id:  -1,
			},
			want: true,
		},
		{
			name: "IsRegisteredUser false",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				id:  -200,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			got, err := p.IsRegisteredUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsRegisteredUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsRegisteredUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgres_IsUserInRole(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		id   int64
		role m.Role
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "IsUserInRole true",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx:  context.Background(),
				id:   -1,
				role: m.RoleAdmin,
			},
			want: true,
		},
		{
			name: "IsUserInRole false",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx:  context.Background(),
				id:   -2,
				role: m.RoleAdmin,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			got, err := p.IsUserInRole(tt.args.ctx, tt.args.id, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUserInRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsUserInRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgres_IsUserPresent(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "IsUserPresent true",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				id:  -1,
			},
			want: true,
		},
		{
			name: "IsUserPresent false",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				id:  -300,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			got, err := p.IsUserPresent(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsUserPresent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsUserPresent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgres_IsWaitingApprove(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "IsWaitingApprove true",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				id:  -3,
			},
			want: true,
		},
		{
			name: "IsWaitingApprove false",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				id:  -2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			got, err := p.IsWaitingApprove(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsWaitingApprove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsWaitingApprove() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgres_UpdateUserRole(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx  context.Context
		user m.User
		role m.Role
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UpdateUserRole success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				user: m.User{
					User: &models.User{
						ID: -3,
					},
					Role: m.RoleRegularUser,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			if err := p.UpdateUserRole(tt.args.ctx, tt.args.user, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgres_RemoveLLM(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
		llm m.LLM
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "RemoveLLM success",
			fields: fields{
				DB: TestDB.DB,
			},
			args: args{
				ctx: context.Background(),
				llm: m.LLM{
					ID: -2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Postgres{
				DB: tt.fields.DB,
			}
			if err := p.RemoveLLM(tt.args.ctx, tt.args.llm); (err != nil) != tt.wantErr {
				t.Errorf("RemoveLLM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
