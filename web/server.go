package web

import (
	"fmt"

	"dbmw/core/connection"
	"dbmw/core/data"
	"dbmw/core/erd"
	"dbmw/core/explorer"
	"dbmw/core/project"
	"dbmw/core/query"
	"dbmw/storage"

	"github.com/gofiber/fiber/v2"
)

// Server wraps the Fiber application instance.
type Server struct {
	app      *fiber.App
	port     int
	connSvc  *connection.Service
	expSvc   *explorer.Service
	qSvc     *query.Service
	dataSvc  *data.Service
	erdSvc   *erd.Service
	projSvc  *project.Service
	cfgStore *storage.ConfigStore
}

// NewServer initializes the Fiber web server instance.
func NewServer(
	port int,
	connSvc *connection.Service,
	expSvc *explorer.Service,
	qSvc *query.Service,
	dataSvc *data.Service,
	erdSvc *erd.Service,
	projSvc *project.Service,
	cfgStore *storage.ConfigStore,
) (*Server, error) {
	app := fiber.New(fiber.Config{
		AppName:               "DBMW v0.0.1",
		DisableStartupMessage: true,
	})

	if err := SetupRouter(app, connSvc, expSvc, qSvc, dataSvc, erdSvc, projSvc, cfgStore); err != nil {
		return nil, err
	}

	return &Server{
		app:      app,
		port:     port,
		connSvc:  connSvc,
		expSvc:   expSvc,
		qSvc:     qSvc,
		dataSvc:  dataSvc,
		erdSvc:   erdSvc,
		projSvc:  projSvc,
		cfgStore: cfgStore,
	}, nil
}

// Start listens on the configured host & port.
func (s *Server) Start() error {
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)
	return s.app.Listen(addr)
}

// Shutdown gracefully stops the Fiber server and all active connections.
func (s *Server) Shutdown() error {
	s.connSvc.CloseAll()
	return s.app.Shutdown()
}

// App returns raw fiber.App (useful for tests).
func (s *Server) App() *fiber.App {
	return s.app
}
