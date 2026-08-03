package handlers

import (
	"dbmw/core/connection"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ConnectionHandler struct {
	connSvc *connection.Service
}

func NewConnectionHandler(connSvc *connection.Service) *ConnectionHandler {
	return &ConnectionHandler{connSvc: connSvc}
}

func (h *ConnectionHandler) List(c *fiber.Ctx) error {
	list, err := h.connSvc.ListConnections()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"connections": list,
		"activeId":    h.connSvc.GetActiveID(),
	})
}

func (h *ConnectionHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := h.connSvc.GetConnection(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if item == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "connection not found"})
	}
	return c.JSON(item)
}

func (h *ConnectionHandler) Save(c *fiber.Ctx) error {
	var body connection.ConnectionConfig
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json body: " + err.Error()})
	}

	saved, err := h.connSvc.SaveConnection(body)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(saved)
}

func (h *ConnectionHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.connSvc.DeleteConnection(id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted", "id": id})
}

func (h *ConnectionHandler) Test(c *fiber.Ctx) error {
	var body connection.ConnectionConfig
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json body: " + err.Error()})
	}

	if err := h.connSvc.TestConnection(c.Context(), body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Connection test successful!",
	})
}

func (h *ConnectionHandler) SetActive(c *fiber.Ctx) error {
	var body struct {
		ID string `json:"id"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "connection id is required"})
	}

	// Verify it can connect
	_, _, err := h.connSvc.GetConnector(c.Context(), body.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.connSvc.SetActiveID(body.ID)
	return c.JSON(fiber.Map{"activeId": body.ID, "message": "Connection activated"})
}
