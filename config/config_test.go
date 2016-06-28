package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockLoadFile(data []byte, err error) func() ([]byte, error) {
	return func() ([]byte, error) {
		return data, err
	}
}

func TestConfigHappyPath(t *testing.T) {
	orgFunc := loadFile
	defer func() {
		loadFile = orgFunc
	}()

	testData := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "Happy Path",
			data: []byte("{}"),
			err:  nil,
		},
		{
			name: "Correct Value format",
			data: []byte("{\"ListenPort\": 8080}"),
			err:  nil,
		},
	}

	for _, test := range testData {
		loadFile = mockLoadFile(test.data, test.err)
		err := loadConfig()
		assert.Nil(t, err, test.name)

		if test.name == "Correct Value format" {
			var expected int
			expected = 8080
			assert.Equal(t, expected, ConfigParams.ListenPort, test.name)
		}
	}
}

func TestConfigError(t *testing.T) {
	orgFunc := loadFile
	defer func() {
		loadFile = orgFunc
	}()

	testData := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "Load File Error",
			data: nil,
			err:  errors.New("Failed to load file"),
		},
		{
			name: "Malformed JSON",
			data: []byte("{"),
			err:  nil,
		},
		{
			name: "InCorrect Value format",
			data: []byte("{\"ListenPort\": \"8080\"}"),
			err:  nil,
		},
	}

	for _, test := range testData {
		loadFile = mockLoadFile(test.data, test.err)
		err := loadConfig()
		assert.NotNil(t, err, test.name)
	}
}
