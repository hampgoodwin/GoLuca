package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/abelgoodwin1988/GoLuca/api"
	"github.com/abelgoodwin1988/GoLuca/internal/config"
	"github.com/abelgoodwin1988/GoLuca/internal/configloader"
	"github.com/abelgoodwin1988/GoLuca/internal/data"
	"github.com/abelgoodwin1988/GoLuca/internal/lucalog"
	"github.com/abelgoodwin1988/GoLuca/internal/setup"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	if err := lucalog.Set(); err != nil {
		fmt.Println("failed to create new logger")
		os.Exit(1)
	}

	// receive interrupt and terminate signals
	go signalsReceiver()

	// Create setup coordinator
	setup.Set()

	// no dependencies
	go configloader.Load()
	go api.Register()

	// have dependencies
	go data.CreateDB()
	// defer data.DBPool.Close()
	go data.Migrate()

	<-setup.C.ReadyForServer()
	server := &http.Server{
		Handler: setup.C.Router.Val,
		Addr:    fmt.Sprintf("%s:%s", config.Env.APIHost, config.Env.APIPort),
	}
	setup.C.Mu.Lock()
	setup.C.Server.Ready = true
	setup.C.Server.Val = server
	setup.C.Mu.Unlock()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			setup.C.Err <- errors.Wrap(err, "api failure")
		}
	}()
	lucalog.Logger.Info("server listening", zap.String("address", setup.C.Server.Val.Addr))
	<-setup.C.Done
}

func signalsReceiver() {
	// listen for exit signals
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
	lucalog.Logger.Info("exiting")
	// thought should be moved out, exitting passes to done channel on setup, and then handles graceful exit
	os.Exit(1)
}
