package main

import (
	"flag"

	"github.com/chengkz2023/My-GVA/server/internal/app/bootstrap"
)

func main() {
	flag.String("c", "", "choose config file.")
	flag.Parse()
	bootstrap.Run()
}
