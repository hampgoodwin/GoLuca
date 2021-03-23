package setup

import (
	"sync"
	"time"

	"github.com/abelgoodwin1988/GoLuca/internal/lucalog"
	"go.uber.org/zap"
)

var C = &Ch{}

func Set() {
	C = &Ch{
		Mu:           sync.RWMutex{},
		Err:          make(chan error, 1),
		Done:         make(chan bool, 1),
		Log:          LogMsg{true},
		ConfigLoader: ConfigLoaderMsg{false},
		Router:       RouterMsg{false, nil},
		DB:           DBMsg{false, nil},
		Migration:    MigrationMsg{false},
		Server:       ServerMsg{false, nil},
	}
}

func (c *Ch) ReadyForDBCreation() chan bool {
	rdyChan := make(chan bool, 1)
	go func() {
		for {
			if c.Log.Ready && c.ConfigLoader.Ready {
				rdyChan <- true
				break
			}
		}
	}()
	return rdyChan
}

func (c *Ch) ReadyForMigration() chan bool {
	rdyChan := make(chan bool, 1)
	go func() {
		for {
			if c.DB.Ready {
				rdyChan <- true
				break
			}
		}
	}()
	return rdyChan
}
func (c *Ch) ReadyForServer() chan bool {
	rdyChan := make(chan bool, 1)
	go func() {
		for {
			time.Sleep(time.Second * 1)
			lucalog.Logger.Info("readyForServer",
				zap.Bool("DB", c.DB.Ready),
				zap.Bool("Migration", c.Migration.Ready),
				zap.Bool("Router", c.Router.Ready),
				zap.Bool("ConfigLoader", c.ConfigLoader.Ready),
				zap.Bool("LOG", c.Log.Ready),
			)
			if c.DB.Ready && c.Migration.Ready && c.Router.Ready && c.ConfigLoader.Ready && c.Log.Ready {
				rdyChan <- true
				break
			}
		}
	}()
	return rdyChan
}

func (c *Ch) Cleanup() chan bool {
	doneCh := make(chan bool, 1)
	lucalog.Logger.Info("closing server")
	c.Server.Val.Close()
	lucalog.Logger.Info("server closed")
	lucalog.Logger.Info("closing db")
	c.DB.Val.Close()
	lucalog.Logger.Info("db closed")
	return doneCh
}
