package test

import (
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/config"
)

func TestNewTestEnvironment(t *testing.T) {
	s := GetScope(t)
	s.Is.True(s != (Scope{}))
	s.Is.Equal(config.Local, s.Env.Config)
}
