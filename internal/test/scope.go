package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matryer/is"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/hampgoodwin/errors"

	servicev1 "github.com/hampgoodwin/GoLuca/gen/proto/go/goluca/service/v1"
	"github.com/hampgoodwin/GoLuca/internal/config"
	"github.com/hampgoodwin/GoLuca/internal/database"
	"github.com/hampgoodwin/GoLuca/internal/environment"
)

type Scope struct {
	Config config.Config
	Env    environment.Environment

	Ctx context.Context
	Is  *is.I

	GRPCTestServer     *grpc.Server
	GRPCBufConn        *bufconn.Listener
	GRPCTestClient     servicev1.GoLucaServiceClient
	GRPCTestClientConn *grpc.ClientConn

	HTTPTestServer *httptest.Server
	HTTPClient     *http.Client

	DB         *pgxpool.Pool
	dbDatabase string
}

func GetScope(t *testing.T) Scope {
	t.Helper()
	s, err := NewScope(t)
	if err != nil {
		s.CleanupScope(t)
		t.Fatalf("creating new scope: %v", err)
	}
	return s
}

func NewScope(t *testing.T) (Scope, error) {
	s := Scope{}
	s.Ctx = context.Background()

	s.Env = environment.TestEnvironment
	s.Env.Log = zap.NewNop()
	s.Is = is.New(t)
	s.HTTPClient = &http.Client{Timeout: time.Second * 30} // TODO: Is this needed?

	if err := s.NewDatabase(t); err != nil {
		return s, errors.Wrap(err, "creating new test db")
	}

	t.Cleanup(func() { s.CleanupScope(t) })
	return s, nil
}

func (s *Scope) NewDatabase(t *testing.T) error {
	t.Helper()
	// Get a random name for a database
	gofakeit.Seed(0)
	s.dbDatabase = strings.ToLower(strings.Replace(gofakeit.Name(), " ", "", -1))
	var err error
	// Create a connection pool on the default database
	s.DB, err = database.NewDatabasePool(s.Ctx, s.Env.Config.Database.ConnectionString())
	if err != nil {
		return errors.Wrap(err, "opening new database connection")
	}

	// Create the new database with the existing database connection pool
	if err := database.CreateDatabase(s.DB, s.dbDatabase); err != nil {
		return errors.Wrap(err, "creating test database")
	}
	// Close the old connection
	s.DB.Close()

	// Open a connection to the newly created random database
	dbCFG := s.Env.Config.Database
	dbCFG.Database = s.dbDatabase
	s.DB, err = database.NewDatabasePool(s.Ctx, dbCFG.ConnectionString())
	if err != nil {
		return errors.Wrap(err, "opening new database connection")
	}

	// run migration on the new database
	if err := database.Migrate(s.DB); err != nil {
		return errors.Wrap(err, fmt.Sprintf("migrating test database %q", s.dbDatabase))
	}

	return nil
}

func (s *Scope) SetHTTP(t *testing.T, handler http.Handler) {
	t.Helper()
	s.HTTPTestServer = httptest.NewServer(handler)
	s.HTTPClient = &http.Client{Timeout: time.Second * 30}
}

func (s *Scope) SetGRPC(t *testing.T, controller servicev1.GoLucaServiceServer) {
	t.Helper()

	buffer := 101024 * 1024
	listener := bufconn.Listen(buffer)
	s.GRPCBufConn = listener

	grpcServer := grpc.NewServer()
	servicev1.RegisterGoLucaServiceServer(grpcServer, controller)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(s.Ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}
	s.GRPCTestClientConn = conn

	goLucaServiceClient := servicev1.NewGoLucaServiceClient(conn)
	s.GRPCTestClient = goLucaServiceClient
}

func (s *Scope) CleanupScope(t *testing.T) {
	t.Helper()

	if s.HTTPTestServer != nil {
		s.HTTPTestServer.Close()
	}

	if s.GRPCTestClientConn != nil {
		s.GRPCTestClientConn.Close()
	}
	if s.GRPCBufConn != nil {
		s.GRPCBufConn.Close()
	}
	if s.GRPCTestServer != nil {
		s.GRPCTestServer.GracefulStop()
	}

	s.DB.Close()

	db, _ := database.NewDatabasePool(s.Ctx, s.Env.Config.Database.ConnectionString())
	err := database.DropDatabase(db, s.dbDatabase)
	s.Is.NoErr(err)
}
