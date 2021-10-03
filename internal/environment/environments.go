package environment

import (
	"github.com/hampgoodwin/GoLuca/internal/config"
)

var TestEnvironment = Environment{
	Config: config.Local,
}
