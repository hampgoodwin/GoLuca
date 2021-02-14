package configloader

import (
	"os"

	"github.com/abelgoodwin1988/GoLuca/internal/config"
)

// Load loads at rest configuration into memory
func Load() {
	config.Env = config.Environment{
		DBConnectionString: os.Getenv("DBCONNSTRING"),
		DBDriverName:       os.Getenv("DBDRIVERNAME"),
	}
}
