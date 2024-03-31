package web

import (
	"fmt"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	_ "github.com/taldoflemis/brain.test/docs"
)

type Router struct {
	config    Config
	zapLogger *zap.Logger
	handlers  []Handler
}

type Handler interface {
	RegisterRoutes(router fiber.Router)
}

func NewRouter(
	config Config,
	zapLogger *zap.Logger,
	handlers []Handler,
) *Router {
	return &Router{
		config:    config,
		zapLogger: zapLogger,
		handlers:  handlers,
	}
}

func (r *Router) Serve() error {
	app := fiber.New(
		fiber.Config{
			ServerHeader:          "brain.test",
			AppName:               "brain.test v0.69420",
			DisableStartupMessage: true,
			ErrorHandler:          ErrorHandlerMiddleware,
		},
	)
	app.Use(recover.New())

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: r.zapLogger,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: r.config.CORSAllowOrigins,
		AllowMethods: r.config.CORSAllowMethods,
		AllowHeaders: r.config.CORSAllowHeaders,
	}))

	api := app.Group(r.config.Prefix)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	api.Get("/swagger/*", swagger.HandlerDefault)

	// Register routes
	for _, handler := range r.handlers {
		handler.RegisterRoutes(api)
	}

	r.zapLogger.Info(
		"Starting server",
		zap.String("listen_ip", r.config.ListenIP),
		zap.Int("port", r.config.Port),
	)

	err := app.Listen(fmt.Sprintf("%s:%d", r.config.ListenIP, r.config.Port))

	if err != nil {
		r.zapLogger.Error("Failed to start server", zap.Error(err))
		return err
	}

	return nil
}
