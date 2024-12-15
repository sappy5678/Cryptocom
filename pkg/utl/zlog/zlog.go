package zlog

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

// Log represents zerolog logger
type Log struct {
	logger *zerolog.Logger
}

// New instantiates new zero logger
func New() *Log {
	z := zerolog.New(os.Stdout)

	return &Log{
		logger: &z,
	}
}

// Log logs using zerolog
func (z *Log) Log(ctx context.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params["source"] = source

	if err != nil {
		params["error"] = err
		z.logger.Error().Fields(params).Msg(msg)

		return
	}

	z.logger.Info().Fields(params).Msg(msg)
}
