package server

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/taldoflemis/gahoot/internal/rooms"
)

func (s *Server) registerRoomsRoutes(router fiber.Router) {
	router.Post("/rooms", s.createRoomHandler)
	router.Delete("/rooms/:id", s.removeRoomHandler)
}

func (s *Server) createRoomHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var req rooms.CreateRoomRequest
	err := c.BodyParser(&req)

	if err != nil {
		log.Println("Failed to parse request", "error", err)
		return fiber.ErrBadRequest
	}

	id, err := s.roomService.Get().CreateRoom(c.Context(), req)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(id)
}

func (s *Server) removeRoomHandler(c *fiber.Ctx) error {
	idReq := c.Params("id")
	id, err := uuid.Parse(idReq)
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = s.roomService.Get().RemoveRoom(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, rooms.ErrRoomNotFound):
			return fiber.ErrNotFound
		default:
			return fiber.ErrInternalServerError
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
