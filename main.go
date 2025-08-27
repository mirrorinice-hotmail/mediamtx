// main executable.
package main

import (
	"log"
	"os"

	"github.com/bluenviron/mediamtx/internal/core"
)

func main() {
	log.Println("---------  Mediamtx Rino --------- v25.08.21.0101")

	s, ok := core.New(os.Args[1:])
	if !ok {
		os.Exit(1)
	}
	s.Wait()
}
