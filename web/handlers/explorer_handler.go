package handlers

import (
	"net/http"

	"dbmw/core/connection"
	"dbmw/core/explorer"
	"github.com/gofiber/fiber/v2"
)

type ExplorerHandler struct {
	connSvc     *connection.Service
	explorerSvc *explorer.Service
}

func NewExplorerHandler(connSvc *connection.Service, explorerSvc *explorer.Service) *ExplorerHandler {
	return &ExplorerHandler{
		connSvc:     connSvc,
		explorerSvc: explorerSvc,
	}
}

func (h *ExplorerHandler) getConnector(c *fiber.Ctx) (connection.Connector, error) {
	connID := c.Query("connId")
	if connID == "" {
		connID = h.connSvc.GetActiveID()
	}
	if connID == "" {
		return nil, fiber.NewError(http.StatusBadRequest, "no active database connection selected")
	}
	conn, _, err := h.connSvc.GetConnector(c.Context(), connID)
	return conn, err
}

func (h *ExplorerHandler) ListDatabases(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.explorerSvc.GetDatabases(c.Context(), conn)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ExplorerHandler) ListSchemas(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	database := c.Query("database")
	res, err := h.explorerSvc.GetSchemas(c.Context(), conn, database)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ExplorerHandler) ListTables(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	schema := c.Query("schema")
	res, err := h.explorerSvc.GetTables(c.Context(), conn, schema)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ExplorerHandler) GetTableDetails(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	table := c.Params("table")
	if table == "" {
		table = c.Query("table")
	}
	if table == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "table parameter is required"})
	}
	schema := c.Query("schema")
	res, err := h.explorerSvc.GetTableDetails(c.Context(), conn, schema, table)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ExplorerHandler) ListColumns(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	table := c.Params("table")
	if table == "" {
		table = c.Query("table")
	}
	if table == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "table parameter is required"})
	}
	schema := c.Query("schema")
	res, err := h.explorerSvc.GetColumns(c.Context(), conn, schema, table)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ExplorerHandler) ListIndexes(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	table := c.Params("table")
	if table == "" {
		table = c.Query("table")
	}
	if table == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "table parameter is required"})
	}
	schema := c.Query("schema")
	res, err := h.explorerSvc.GetIndexes(c.Context(), conn, schema, table)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ExplorerHandler) ListForeignKeys(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	table := c.Params("table")
	if table == "" {
		table = c.Query("table")
	}
	if table == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "table parameter is required"})
	}
	schema := c.Query("schema")
	res, err := h.explorerSvc.GetForeignKeys(c.Context(), conn, schema, table)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *ExplorerHandler) ListViews(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	schema := c.Query("schema")
	res, err := h.explorerSvc.GetViews(c.Context(), conn, schema)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
