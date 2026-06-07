package main

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/bootstrap"

	_ "go.uber.org/automaxprocs"
)

func main() {
	bootstrap.Run()
}
