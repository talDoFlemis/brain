package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"

	"github.com/taldoflemis/gahoot/internal/server"
)

func main() {
	if err := weaver.Run(context.Background(), server.Serve); err != nil {
		panic(err)
	}
}
