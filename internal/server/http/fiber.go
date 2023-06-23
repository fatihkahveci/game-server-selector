package http

import (
	"game-server-selector/internal/domain"
	"game-server-selector/internal/models"
	"game-server-selector/internal/services"
	"game-server-selector/internal/validator"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type FiberServer struct {
	app           *fiber.App
	serverDomain  domain.ServerDomain
	configService services.ConfigService
	metricService services.MetricService
}

func NewFiberServer(
	s domain.ServerDomain,
	c services.ConfigService,
	m services.MetricService,
) *FiberServer {
	app := fiber.New()

	return &FiberServer{
		app:           app,
		serverDomain:  s,
		configService: c,
		metricService: m,
	}
}

func (s *FiberServer) InitMiddleware() {
	s.app.Use(recover.New())
	s.app.Use(logger.New())
	s.app.Use(s.metricMiddleware)
}

func (s *FiberServer) Router() error {
	v1 := s.app.Group("/v1")
	v1.Get("/server/list", s.serverList)
	v1.Post("/server/search", s.serverSearch)

	v1.Post("/server/create", s.serverCreate)
	v1.Post("/server/update/:id", s.serverUpdate)

	return nil
}

func (s *FiberServer) Stop() error {
	return s.app.Shutdown()
}

func (s *FiberServer) Start() error {
	s.InitMiddleware()
	s.Router()

	port := s.configService.GetHttpPort()
	certPath := s.configService.GetSslCertPath()
	keyPath := s.configService.GetSslKeyPath()

	if certPath != "" && keyPath != "" {
		return s.app.ListenTLS(port, certPath, keyPath)
	}
	return s.app.Listen(port)
}

func (s *FiberServer) serverList(ctx *fiber.Ctx) error {
	list, err := s.serverDomain.ListServers()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"list":    list,
	})
}

func (s *FiberServer) serverCreate(ctx *fiber.Ctx) error {
	serverRequest := models.CreateServerRequest{}
	if err := ctx.BodyParser(&serverRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	errors := validator.ValidateServer(serverRequest)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)

	}
	serverData, err := s.serverDomain.CreateServer(serverRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.JSON(fiber.Map{
		"server":  serverData,
		"success": true,
	})
}

func (s *FiberServer) serverSearch(ctx *fiber.Ctx) error {
	searchRequest := []models.SearchRequest{}
	if err := ctx.BodyParser(&searchRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	list, err := s.serverDomain.SearchServers(searchRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"list":    list,
	})
}

func (s *FiberServer) serverUpdate(ctx *fiber.Ctx) error {
	serverRequest := models.UpdateServerRequest{}
	if err := ctx.BodyParser(&serverRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	errors := validator.ValidateServerUpdate(serverRequest)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)

	}
	server, err := s.serverDomain.UpdateServer(ctx.Params("id"), serverRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"server":  server,
	})
}

func (s *FiberServer) metricMiddleware(ctx *fiber.Ctx) error {
	start := time.Now()
	method := ctx.Route().Method
	err := ctx.Next()
	status := fiber.StatusInternalServerError

	if ctx.Route().Path == "metrics" {
		return ctx.Next()
	}

	if err != nil {
		if e, ok := err.(*fiber.Error); ok {
			status = e.Code
		}
	} else {
		status = ctx.Response().StatusCode()
	}

	handler := ctx.Route().Path
	statusCode := strconv.Itoa(status)
	s.metricService.IncRequestTotal(handler, statusCode, method)
	elapsed := float64(time.Since(start).Nanoseconds()) / 1e9
	s.metricService.ObserveRequestDuration(handler, statusCode, method, elapsed)

	return err
}
