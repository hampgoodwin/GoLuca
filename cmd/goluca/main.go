package goluca

import (
	"fmt"

	"github.com/abelgoodwin1988/GoLuca/internal/configloader"
	"github.com/abelgoodwin1988/GoLuca/internal/data"
)

func main() {
	configloader.Load()
	data.CreateDB()
	run()
}

func run() {
	fmt.Println("run")
}
