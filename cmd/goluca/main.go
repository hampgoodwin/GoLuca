package goluca

import (
	"fmt"

	"github.com/abelgoodwin1988/GoLuca/internal/configloader"
	"github.com/abelgoodwin1988/GoLuca/internal/db"
)

func main() {
	configloader.Load()
	db.Create()
	run()
}

func run() {
	fmt.Println("run")
}
