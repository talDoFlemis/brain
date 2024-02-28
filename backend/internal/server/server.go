package server

import (
	"context"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/taldoflemis/gahoot/internal/rooms"
)

type Server struct {
	weaver.Implements[weaver.Main]

	F *fiber.App

	roomService weaver.Ref[rooms.Roomer]

	api weaver.Listener `weaver:"api"`
}

func Serve(_ context.Context, app *Server) error {
	app.F = fiber.New()
	app.RegisterFiberRoutes()
	srv := otelhttp.NewHandler(adaptor.FiberApp(app.F), "gahoot")
	return http.Serve(app.api, srv)
}
