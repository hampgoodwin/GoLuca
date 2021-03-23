package setup

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
)

type LogMsg struct {
	Ready bool
}

type RouterMsg struct {
	Ready bool
	Val   *chi.Mux
}

type DBMsg struct {
	Ready bool
	Val   *pgxpool.Pool
}
type MigrationMsg struct {
	Ready bool
}

type ConfigLoaderMsg struct {
	Ready bool
}

type ServerMsg struct {
	Ready bool
	Val   *http.Server
}

type Ch struct {
	Mu           sync.RWMutex
	Err          chan error
	Done         chan bool
	Log          LogMsg
	ConfigLoader ConfigLoaderMsg
	DB           DBMsg
	Migration    MigrationMsg
	Router       RouterMsg
	Server       ServerMsg
}
