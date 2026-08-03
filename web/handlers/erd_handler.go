package handlers

import (
	"net/http"

	"dbmw/core/connection"
	"dbmw/core/erd"
	"github.com/gofiber/fiber/v2"
)

type ERDHandler struct {
	connSvc *connection.Service
	erdSvc  *erd.Service
}

func NewERDHandler(connSvc *connection.Service, erdSvc *erd.Service) *ERDHandler {
	return &ERDHandler{
		connSvc: connSvc,
		erdSvc:  erdSvc,
	}
}

func (h *ERDHandler) Generate(c *fiber.Ctx) error {
	connID := c.Query("connId")
	if connID == "" {
		connID = h.connSvc.GetActiveID()
	}
	if connID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "no active database connection selected"})
	}

	conn, _, err := h.connSvc.GetConnector(c.Context(), connID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	schema := c.Query("schema")
	graph, err := h.erdSvc.GenerateGraph(c.Context(), conn, schema)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(graph)
}
