package main

import (
	"os"

	"github.com/byteforge-run/tinyshop-tester/internal/stages"
	tester_utils "github.com/byteforge-run/tester-utils"
)

func main() {
	definition := stages.GetDefinition()
	os.Exit(tester_utils.Run(os.Args[1:], definition))
}
