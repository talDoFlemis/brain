package server

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/gofiber/fiber/v2"

	"github.com/taldoflemis/gahoot/internal/rooms"
)

type Server struct {
	weaver.Implements[weaver.Main]

	f *fiber.App

	roomService weaver.Ref[rooms.Roomer]

	api weaver.Listener `weaver:"api"`
}

func Serve(_ context.Context, app *Server) error {
	app.f = fiber.New()
	app.RegisterFiberRoutes()
	return app.f.Listener(app.api)
}
