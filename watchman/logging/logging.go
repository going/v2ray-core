package logging

import (
	"sync"

	"go.uber.org/zap"
)

// Logging contains zap logger
type Logging struct {
	Logger *zap.Logger
}

var once sync.Once // for the singleton
var log *Logging

// GetLogger returns a singleton instace for logging
func GetInstance() *Logging {
	dlog, err := zap.NewDevelopment()
	if err != nil {
		return nil
	}
	once.Do(func() {
		if log == nil {
			log = &Logging{
				Logger: dlog,
			}
		}
	})
	return log
}
