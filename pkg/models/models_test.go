package models

import "testing"

func TestRole_String(t *testing.T) {
	tests := []struct {
		name string
		r    Role
		want string
	}{
		{
			name: "test regular user role",
			r:    RoleRegularUser,
			want: "regular user",
		},
		{
			name: "test admin role",
			r:    RoleAdmin,
			want: "admins",
		},
		{
			name: "test unknown role",
			r:    RoleUnknown,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
