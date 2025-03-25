package config

import (
	"errors"
	"testing"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
)

type inputData struct {
	configPath     string
	configName     string
	inputEnvValues map[string]string
	inputConfig    models.Config
}

type expectedData struct {
	errReadInConfig error
	expectedConfig  models.Config
}

func TestConfig(t *testing.T) {
	testTable := []struct {
		name     string
		input    inputData
		expected expectedData
	}{
		{
			name: "Init config with test data params",
			input: inputData{
				configPath:  "testdata",
				configName:  "config_test",
				inputConfig: models.Config{},
				inputEnvValues: map[string]string{
					"DB_PASSWORD": "testpassword",
					"SALT":        "testsalt",
					"SIGNEDKEY":   "testsignedkey",
				},
			},
			expected: expectedData{
				errReadInConfig: nil,
				expectedConfig: models.Config{
					DB: models.DBConfig{
						Host:     "testhost",
						Port:     "testport",
						Username: "testusername",
						Password: "testpassword",
						DBName:   "testdbname",
						SSLMode:  "testsslmode",
					},
					Auth: models.AuthConfig{
						Salt:      "testsalt",
						SignedKey: "testsignedkey",
					},
				},
			},
		},
		{
			name: "Error check in config path param",
			input: inputData{
				configPath: "error_path",
				configName: "config_test",
			},
			expected: expectedData{
				errReadInConfig: customerrors.ErrReadInConfig,
			},
		},
		{
			name: "Error check in config name param",
			input: inputData{
				configPath: "testdata",
				configName: "errorpath",
			},
			expected: expectedData{
				errReadInConfig: customerrors.ErrReadInConfig,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			for key, inputEnvValue := range testCase.input.inputEnvValues {
				t.Setenv(key, inputEnvValue)

			}

			err := InitConfig(
				testCase.input.configPath,
				testCase.input.configName,
				&testCase.input.inputConfig,
			)

			if !errors.Is(err, testCase.expected.errReadInConfig) {
				t.Error("errInitConfit != errReadInConfig")
			}

			if errors.Is(err, customerrors.ErrReadInConfig) {
				return
			}

			if testCase.expected.expectedConfig != testCase.input.inputConfig {
				t.Error("input and expected configs are different")
			}
		})
	}
}
