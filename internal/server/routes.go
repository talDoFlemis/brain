package server

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/taldoflemis/gahoot/cmd/web"
)

func (s *Server) RegisterFiberRoutes() {
	s.f.Get("/", s.HelloWorldHandler)

	s.f.Static("/js", "./cmd/web/js")

	s.f.Get("/web", adaptor.HTTPHandler(templ.Handler(web.HelloForm())))

	s.f.Post("/hello", func(c *fiber.Ctx) error {
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
