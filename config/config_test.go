package config

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
					"PORT": "1234",
				},
			},
			wantConfig: &Config{
				Port: 1234,
			},
			wantErr: false,
		},
		{
			name: "should return an error if required variables are not set",
			args: args{
				createTempEnvFile: true,
				envVars:           map[string]string{},
			},
			wantConfig: &Config{},
			wantErr:    true,
		},
		{
			name: "should return an error if env file does not exist",
			args: args{
				createTempEnvFile: false,
				envVars:           map[string]string{},
			},
			wantConfig: &Config{},
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
					t.Cleanup(func() {
						if err := os.Unsetenv(key); err != nil {
							t.Fatalf("error unsetting environment variable: %v", err)
						}
					})
				}
			}

			config, err := Load(ctx, envPath)

			if tt.wantErr != (err != nil) {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(config, tt.wantConfig) {
				t.Fatalf("Load() = %v, want %v", config, tt.wantConfig)
			}
		})
	}
}
