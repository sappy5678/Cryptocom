package zlog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sappy5678/cryptocom/pkg/utl/zlog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	log := zlog.New()
	assert.NotNil(t, log)
}

func TestLog(t *testing.T) {
	log := zlog.New()

	assert.NotPanics(t, func() {
		log.Log(context.Background(), "test", "test", nil, nil)
	})
	assert.NotPanics(t, func() {
		log.Log(context.Background(), "test", "test", errors.New("test"), nil)
	})
	assert.NotPanics(t, func() {
		log.Log(context.Background(), "test", "test", errors.New("test"), map[string]interface{}{"test": "test"})
	})
}
