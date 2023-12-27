package server

import (
	"github.com/taldoflemis/gahoot/cmd/web"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.Get("/", s.HelloWorldHandler)
	s.Get("/health", s.healthHandler)

	s.Static("/js", "./cmd/web/js")

	s.Get("/web", adaptor.HTTPHandler(templ.Handler(web.HelloForm())))

	s.Post("/hello", func(c *fiber.Ctx) error {
		return web.HelloWebHandler(c)
	})
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := map[string]string{
		"message": "Hello World",
	}
	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
