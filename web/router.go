package web

import (
	"dbmw/core/connection"
	"dbmw/core/data"
	"dbmw/core/erd"
	"dbmw/core/explorer"
	"dbmw/core/project"
	"dbmw/core/query"
	"dbmw/storage"
	"dbmw/web/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupRouter configures Fiber middlewares, API routes, and embedded SPA static filesystem.
func SetupRouter(
	app *fiber.App,
	connSvc *connection.Service,
	expSvc *explorer.Service,
	qSvc *query.Service,
	dataSvc *data.Service,
	erdSvc *erd.Service,
	projSvc *project.Service,
	cfgStore *storage.ConfigStore,
) error {
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))

	connHandler := handlers.NewConnectionHandler(connSvc)
	expHandler := handlers.NewExplorerHandler(connSvc, expSvc)
	qHandler := handlers.NewQueryHandler(connSvc, qSvc)
	dataHandler := handlers.NewDataHandler(connSvc, dataSvc)
	erdHandler := handlers.NewERDHandler(connSvc, erdSvc)
	projHandler := handlers.NewProjectHandler(projSvc)
	cfgHandler := handlers.NewConfigHandler(cfgStore)

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "version": "0.0.1"})
	})

	// Config
	api.Get("/config", cfgHandler.Get)
	api.Post("/config", cfgHandler.Update)

	// Connections
	api.Get("/connections", connHandler.List)
	api.Post("/connections", connHandler.Save)
	api.Post("/connections/test", connHandler.Test)
	api.Post("/connections/active", connHandler.SetActive)
	api.Get("/connections/:id", connHandler.Get)
	api.Delete("/connections/:id", connHandler.Delete)

	// Explorer
	api.Get("/explorer/databases", expHandler.ListDatabases)
	api.Get("/explorer/schemas", expHandler.ListSchemas)
	api.Get("/explorer/tables", expHandler.ListTables)
	api.Get("/explorer/tables/:table/details", expHandler.GetTableDetails)
	api.Get("/explorer/columns/:table", expHandler.ListColumns)
	api.Get("/explorer/indexes/:table", expHandler.ListIndexes)
	api.Get("/explorer/foreign-keys/:table", expHandler.ListForeignKeys)
	api.Get("/explorer/views", expHandler.ListViews)

	// Query
	api.Post("/query/execute", qHandler.Execute)
	api.Get("/query/history", qHandler.ListHistory)
	api.Delete("/query/history", qHandler.ClearHistory)
	api.Post("/query/export/csv", qHandler.ExportCSV)
	api.Post("/query/export/json", qHandler.ExportJSON)

	// Data browse & CRUD
	api.Get("/data/browse/:table", dataHandler.Browse)
	api.Post("/data/browse/:table", dataHandler.Browse)
	api.Post("/data/insert/:table", dataHandler.Insert)
	api.Post("/data/update/:table", dataHandler.Update)
	api.Post("/data/delete/:table", dataHandler.Delete)

	// ERD
	api.Get("/erd/generate", erdHandler.Generate)

	// Project
	api.Get("/project/detect", projHandler.Detect)
	api.Post("/project/generate", projHandler.GenerateConfig)

	// Serve embedded Frontend SPA
	distFS, err := GetFileSystem()
	if err == nil {
		app.Use("/", filesystem.New(filesystem.Config{
			Root:         distFS,
			Index:        "index.html",
			NotFoundFile: "index.html",
		}))
	}

	return nil
}
