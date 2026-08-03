package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"dbmw/core/connection"
	"dbmw/core/query"
	"github.com/gofiber/fiber/v2"
)

type QueryHandler struct {
	connSvc  *connection.Service
	querySvc *query.Service
}

func NewQueryHandler(connSvc *connection.Service, querySvc *query.Service) *QueryHandler {
	return &QueryHandler{
		connSvc:  connSvc,
		querySvc: querySvc,
	}
}

func (h *QueryHandler) Execute(c *fiber.Ctx) error {
	var body struct {
		ConnectionID string `json:"connectionId"`
		Query        string `json:"query"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json: " + err.Error()})
	}

	connID := body.ConnectionID
	if connID == "" {
		connID = h.connSvc.GetActiveID()
	}
	if connID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "no database connection selected"})
	}

	conn, cfg, err := h.connSvc.GetConnector(c.Context(), connID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	name := "Unknown"
	if cfg != nil {
		name = cfg.Name
	}

	res, err := h.querySvc.Execute(c.Context(), conn, connID, name, body.Query)
	if err != nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"error":           err.Error(),
			"executionTimeMs": res.ExecutionTimeMs,
		})
	}

	return c.JSON(res)
}

func (h *QueryHandler) ListHistory(c *fiber.Ctx) error {
	connID := c.Query("connId")
	limitStr := c.Query("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	hist, err := h.querySvc.GetHistory(c.Context(), connID, limit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(hist)
}

func (h *QueryHandler) ClearHistory(c *fiber.Ctx) error {
	connID := c.Query("connId")
	if err := h.querySvc.ClearHistory(c.Context(), connID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "history cleared"})
}

func (h *QueryHandler) ExportCSV(c *fiber.Ctx) error {
	var body query.QueryResult
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	data, err := h.querySvc.ExportCSV(&body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="export.csv"`)
	return c.Send(data)
}

func (h *QueryHandler) ExportJSON(c *fiber.Ctx) error {
	var body query.QueryResult
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	data, err := h.querySvc.ExportJSON(&body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="export.json"`))
	return c.Send(data)
}
