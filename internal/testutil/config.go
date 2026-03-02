package testutil

import (
	"sync"
	"testing"

	"github.com/Nikkoz/task-service/internal/config"
)

var (
	cfgOnce sync.Once
	cfg     config.Config
	cfgErr  error
)

func GetConfig(t *testing.T) config.Config {
	t.Helper()

	cfgOnce.Do(func() {
		cfg, cfgErr = config.Load()
	})

	if cfgErr != nil {
		t.Fatalf("load config: %v", cfgErr)
	}

	return cfg
}
