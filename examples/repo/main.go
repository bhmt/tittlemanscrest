package main

import (
	"flag"
	"os"

	"github.com/bhmt/tittlemanscrest/examples/repo/click"
)

func main() {
	var repoType string

	flag.StringVar(&repoType, "r", "click", "set the repository example type")
	flag.Parse()

	switch repoType {
	case "click":
		click.Click()
	default:
		flag.Usage()
		os.Exit(1)
	}
}
