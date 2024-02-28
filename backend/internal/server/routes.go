package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) RegisterFiberRoutes() {
	api := s.F.Group("/api")

	api.Get("/", s.HelloWorldHandler)

	s.registerRoomsRoutes(api)
}

func (s *Server) HelloWorldHandler(c *fiber.Ctx) error {
	resp := map[string]string{
		"message": "Hello World",
	}
	return c.JSON(resp)
}
