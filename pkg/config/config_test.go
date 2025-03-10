package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestConfig_Load(t *testing.T) {
	type fields struct {
		Port          int
		Dsn           string
		TgSecret      string
		InitBotSecret string
	}
	type args struct {
		v *viper.Viper
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test_Load",
			fields: fields{
				Port:          8080,
				Dsn:           "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
				TgSecret:      "secret",
				InitBotSecret: "secret",
			},
			args: args{
				v: viper.New(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Port:          tt.fields.Port,
				Dsn:           tt.fields.Dsn,
				TgSecret:      tt.fields.TgSecret,
				InitBotSecret: tt.fields.InitBotSecret,
			}
			if err := c.Load(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "TestNewConfig",
			want: &Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
