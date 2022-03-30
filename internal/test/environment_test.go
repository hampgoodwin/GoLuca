package test

import (
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/config"
)

func TestNewTestEnvironment(t *testing.T) {
	s, err := NewScope(t)
	s.Is.NoErr(err)
	s.Is.True(s != (Scope{}))
	s.Is.Equal(config.Local, s.Env.Config)
	s.Is.NoErr(err)
}
