package server

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/taldoflemis/gahoot/cmd/web"
)

func (s *Server) RegisterFiberRoutes() {
	s.F.Get("/", s.HelloWorldHandler)

	s.F.Static("/js", "./cmd/web/js")

	s.F.Get("/web", adaptor.HTTPHandler(templ.Handler(web.HelloForm())))

	s.F.Post("/hello", func(c *fiber.Ctx) error {
		return web.HelloWebHandler(c)
	})
	s.registerRoomsRoutes()
}

func (s *Server) HelloWorldHandler(c *fiber.Ctx) error {
	resp := map[string]string{
		"message": "Hello World",
	}
	return c.JSON(resp)
}
