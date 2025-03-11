package ai

import (
	"context"
	"errors"
	"testing"
)

const haiku = `Солнце зашумело,
Серые облака в небе,
Зима — покойная.`

type ClientStub struct{}

func (c *ClientStub) ChatCompletion(ctx context.Context, model, req string) (resp string, err error) {
	if req == "error" {
		return "", errors.New("error")
	}
	return haiku, nil
}

func TestClient_ChatCompletion(t *testing.T) {
	type args struct {
		ctx context.Context
		req string
	}
	tests := []struct {
		name     string
		args     args
		wantResp string
		wantErr  bool
	}{
		{
			name: "haiku",
			args: args{
				ctx: context.Background(),
				req: "any",
			},
			wantResp: haiku,
		},
		{
			name: "error",
			args: args{
				ctx: context.Background(),
				req: "error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientStub{}
			gotResp, err := c.ChatCompletion(tt.args.ctx, "local", tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatCompletion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResp != tt.wantResp {
				t.Errorf("ChatCompletion() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
