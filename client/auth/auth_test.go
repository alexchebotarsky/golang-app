package auth

import (
	"reflect"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func TestClient_findRole(t *testing.T) {
	c := &Client{
		authSecret: []byte("secret"),
		roles: []Role{
			{
				Name:        "user_role",
				AccessLevel: 10,
				Key:         "user_key",
			},
			{
				Name:        "admin_role",
				AccessLevel: 30,
				Key:         "admin_key",
			},
		},
		TokenTTL:      time.Hour,
		SigningMethod: jwt.SigningMethodHS256,
	}

	type args struct {
		roleName string
		roleKey  string
	}
	tests := []struct {
		name    string
		args    args
		want    Role
		wantErr bool
	}{
		{
			name:    "correct name and key",
			args:    args{roleName: "user_role", roleKey: "user_key"},
			want:    Role{Name: "user_role", AccessLevel: 10, Key: "user_key"},
			wantErr: false,
		},
		{
			name:    "invalid key",
			args:    args{roleName: "user_role", roleKey: "invalid_key"},
			want:    Role{},
			wantErr: true,
		},
		{
			name:    "invalid name",
			args:    args{roleName: "invalid_role", roleKey: "user_key"},
			want:    Role{},
			wantErr: true,
		},
		{
			name:    "empty key",
			args:    args{roleName: "user_role", roleKey: ""},
			want:    Role{},
			wantErr: true,
		},
		{
			name:    "empty name",
			args:    args{roleName: "", roleKey: "user_key"},
			want:    Role{},
			wantErr: true,
		},
		{
			name:    "correct name and key but for different roles",
			args:    args{roleName: "user_name", roleKey: "admin_key"},
			want:    Role{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.findRole(tt.args.roleName, tt.args.roleKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.findRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.findRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
