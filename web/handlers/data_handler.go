package handlers

import (
	"net/http"

	"dbmw/core/connection"
	"dbmw/core/data"
	"github.com/gofiber/fiber/v2"
)

type DataHandler struct {
	connSvc *connection.Service
	dataSvc *data.Service
}

func NewDataHandler(connSvc *connection.Service, dataSvc *data.Service) *DataHandler {
	return &DataHandler{
		connSvc: connSvc,
		dataSvc: dataSvc,
	}
}

func (h *DataHandler) getConnector(c *fiber.Ctx) (connection.Connector, error) {
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

func (h *DataHandler) Browse(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	table := c.Params("table")
	schema := c.Query("schema")

	var opts data.BrowseOptions
	if len(c.Body()) > 0 {
		_ = c.BodyParser(&opts)
	}
	if opts.Page <= 0 {
		opts.Page = c.QueryInt("page", 1)
	}
	if opts.PageSize <= 0 {
		opts.PageSize = c.QueryInt("pageSize", 25)
	}

	page, err := h.dataSvc.Browse(c.Context(), conn, schema, table, opts)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(page)
}

func (h *DataHandler) Insert(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	table := c.Params("table")
	schema := c.Query("schema")

	var values map[string]any
	if err := c.BodyParser(&values); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid values json: " + err.Error()})
	}

	if err := h.dataSvc.Insert(c.Context(), conn, schema, table, values); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "row inserted successfully"})
}

func (h *DataHandler) Update(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	table := c.Params("table")
	schema := c.Query("schema")

	var body struct {
		PrimaryKey map[string]any `json:"primaryKey"`
		Values     map[string]any `json:"values"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body: " + err.Error()})
	}

	if err := h.dataSvc.Update(c.Context(), conn, schema, table, body.PrimaryKey, body.Values); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "row updated successfully"})
}

func (h *DataHandler) Delete(c *fiber.Ctx) error {
	conn, err := h.getConnector(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	table := c.Params("table")
	schema := c.Query("schema")

	var pk map[string]any
	if err := c.BodyParser(&pk); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid primaryKey body: " + err.Error()})
	}

	if err := h.dataSvc.Delete(c.Context(), conn, schema, table, pk); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "row deleted successfully"})
}
