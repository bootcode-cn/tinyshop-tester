package main

import (
	"os"

	"github.com/bootcode-cn/tinyshop-tester/internal/stages"
	tester_utils "github.com/bootcode-cn/tester-utils"
)

func main() {
	definition := stages.GetDefinition()
	os.Exit(tester_utils.Run(os.Args[1:], definition))
}
