package config_test

import (
	"testing"

	"github.com/sappy5678/cryptocom/pkg/utl/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestLoad(t *testing.T) {
	defer goleak.VerifyNone(t)
	cases := []struct {
		name     string
		path     string
		wantData *config.Configuration
		wantErr  bool
	}{
		{
			name:    "Fail on non-existing file",
			path:    "notExists",
			wantErr: true,
		},
		{
			name:    "Fail on wrong file format",
			path:    "testdata/config.invalid.yaml",
			wantErr: true,
		},
		{
			name: "Success",
			path: "testdata/config.testdata.yaml",
			wantData: &config.Configuration{

				Server: &config.Server{
					Port:         ":8080",
					Debug:        true,
					ReadTimeout:  15,
					WriteTimeout: 20,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.path)
			assert.Equal(t, tt.wantData, cfg)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
