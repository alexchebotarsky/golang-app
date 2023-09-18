package env

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		envVars           map[string]string
		createTempEnvFile bool
	}
	tests := []struct {
		name       string
		args       args
		wantConfig *Config
		wantErr    bool
	}{
		{
			name: "should load environment variables from env file",
			args: args{
				createTempEnvFile: true,
				envVars: map[string]string{
					"DATABASE_USER":     "test_user",
					"DATABASE_PASSWORD": "test_pass",
					"DATABASE_HOST":     "test_host",
					"DATABASE_PORT":     "2000",
					"DATABASE_NAME":     "test_db",
					"DATABASE_OPTIONS":  "?database_option=test",

					"AUTH_SECRET":     "test_auth_secret",
					"AUTH_ADMIN_KEY":  "test_auth_admin_key",
					"AUTH_EDITOR_KEY": "test_auth_editor_key",
					"AUTH_VIEWER_KEY": "test_auth_viewer_key",

					"EXAMPLE_ENDPOINT": "test_example_endpoint",

					"PORT":         "1000",
					"SERVICE_NAME": "test_service_name",
				},
			},
			wantConfig: &Config{
				DatabaseUser:     "test_user",
				DatabasePassword: "test_pass",
				DatabaseHost:     "test_host",
				DatabasePort:     2000,
				DatabaseName:     "test_db",
				DatabaseOptions:  "?database_option=test",

				AuthSecret:    "test_auth_secret",
				AuthAdminKey:  "test_auth_admin_key",
				AuthEditorKey: "test_auth_editor_key",
				AuthViewerKey: "test_auth_viewer_key",

				ExampleEndpoint: "test_example_endpoint",

				Port:        1000,
				ServiceName: "test_service_name",
			},
			wantErr: false,
		},
		{
			name: "should return an error if required variables are not set",
			args: args{
				createTempEnvFile: true,
				envVars: map[string]string{
					"DATABASE_HOST":    "test_host",
					"DATABASE_PORT":    "2000",
					"DATABASE_OPTIONS": "?database_option=test",

					"EXAMPLE_ENDPOINT": "test_example_endpoint",

					"PORT":         "1000",
					"SERVICE_NAME": "test_service_name",
				},
			},
			wantConfig: nil,
			wantErr:    true,
		},
		{
			name: "should return an error if env file does not exist",
			args: args{
				createTempEnvFile: false,
				envVars:           map[string]string{},
			},
			wantConfig: nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			var envPath string
			if tt.args.createTempEnvFile {
				envFile, err := os.CreateTemp("./", ".test.*.env")
				if err != nil {
					t.Fatalf("error creating temporary env file: %v", err)
				}
				t.Cleanup(func() {
					if err := envFile.Close(); err != nil {
						t.Fatalf("error closing temporary env file: %v", err)
					}

					if err := os.Remove(envFile.Name()); err != nil {
						t.Fatalf("error removing temporary env file: %v", err)
					}
				})
				envPath = envFile.Name()

				for key, value := range tt.args.envVars {
					if _, err := envFile.WriteString(fmt.Sprintf("%s=%v\n", key, value)); err != nil {
						t.Fatalf("error writing to temporary env file: %v", err)
					}
				}
				t.Cleanup(func() {
					for key := range tt.args.envVars {
						if err := os.Unsetenv(key); err != nil {
							t.Fatalf("error unsetting environment variable: %v", err)
						}
					}
				})
			}

			config, err := LoadConfig(ctx, envPath)

			if tt.wantErr != (err != nil) {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(config, tt.wantConfig) {
				t.Fatalf("Load() = %v, want %v", config, tt.wantConfig)
			}
		})
	}
}
