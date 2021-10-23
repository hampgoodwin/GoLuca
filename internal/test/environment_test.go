package test

import (
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/matryer/is"
)

func TestNewTestEnvironment(t *testing.T) {
	i := is.New(t)
	s, err := NewScope(t)
	i.NoErr(err)
	i.True(s != (Scope{}))
	i.Equal(config.Local, s.Env.Config)
	i.NoErr(err)
}
