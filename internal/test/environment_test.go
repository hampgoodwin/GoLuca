package test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/hampgoodwin/GoLuca/internal/config"
)

func TestGetScope(t *testing.T) {
	s := GetScope(t)
	s.Is.True(s != (Scope{}))
	s.Is.Equal(config.Local, s.Env.Config)
}

func TestNewScope(t *testing.T) {
	s, err := NewScope(t)
	s.Is.NoErr(err)
	s.Is.True(s != (Scope{}))
}
func TestNewDatabase_ContextCanceled(t *testing.T) {
	s, err := NewScope(t)
	s.Is.NoErr(err)
	s.Is.True(s != (Scope{}))
	ctx, done := context.WithCancel(s.Ctx)
	s.Ctx = ctx
	done()

	err = s.NewDatabase(t)
	s.Is.True(strings.Contains(err.Error(), "context canceled"))
}

func TestNewScope_SetHTTP(t *testing.T) {
	s, err := NewScope(t)
	s.Is.NoErr(err)
	s.Is.True(s != (Scope{}))

	mux := http.NewServeMux()

	s.SetHTTP(t, mux)

	r, err := http.Get(s.HTTPTestServer.URL)
	s.Is.NoErr(err)
	s.Is.True(r != nil)
}
