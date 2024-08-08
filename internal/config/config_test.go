package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfig_UpdateByEnv(t *testing.T) {
	t.Setenv("RUN_ADDRESS", ":8080")

	type envs []struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		envs    envs
		wantErr bool
	}{
		{
			name: "[Positive] Parse env to config",
			envs: envs{
				{"RUN_ADDRESS", ":8080"},
				{"DATABASE_URI", "postgres://test@test"},
				{"ACCRUAL_SYSTEM_ADDRESS", "http://asa.test/"},
			},
			wantErr: false,
		},
		{
			name: "[Negative] Parse empty string from env",
			envs: envs{
				{"RUN_ADDRESS", ""},
				{"DATABASE_URI", ""},
				{"ACCRUAL_SYSTEM_ADDRESS", ""},
			},
			wantErr: true,
		},
		{
			name: "[Negative] Parse empty string from one of the env",
			envs: envs{
				{"RUN_ADDRESS", ""},
				{"DATABASE_URI", "postgres://test@test"},
				{"ACCRUAL_SYSTEM_ADDRESS", "http://asa.test/"},
			},
			wantErr: true,
		},
		{
			name: "[Negative] Parse empty string from one of the env 2",
			envs: envs{
				{"RUN_ADDRESS", ":8080"},
				{"DATABASE_URI", ""},
				{"ACCRUAL_SYSTEM_ADDRESS", "http://asa.test/"},
			},
			wantErr: true,
		},
		{
			name: "[Negative] Parse empty string from one of the env 3",
			envs: envs{
				{"RUN_ADDRESS", ":8080"},
				{"DATABASE_URI", "postgres://test@test"},
				{"ACCRUAL_SYSTEM_ADDRESS", ""},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			for _, env := range test.envs {
				t.Setenv(env.name, env.value)
			}
			tryNextEnv := true

			for _, v := range test.envs {
				if tryNextEnv {
					value := os.Getenv(v.name)
					if test.wantErr {
						if len(value) != 0 {
							continue
						} else {
							tryNextEnv = false
							require.Len(t, value, 0)
						}
					}
					require.Equal(t, v.value, value)
				}

			}

		})
	}
}
