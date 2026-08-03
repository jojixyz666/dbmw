package handlers

import (
	"net/http"

	"dbmw/storage"
	"github.com/gofiber/fiber/v2"
)

type ConfigHandler struct {
	configStore *storage.ConfigStore
}

func NewConfigHandler(store *storage.ConfigStore) *ConfigHandler {
	return &ConfigHandler{configStore: store}
}

func (h *ConfigHandler) Get(c *fiber.Ctx) error {
	cfg, err := h.configStore.Get()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cfg)
}

func (h *ConfigHandler) Update(c *fiber.Ctx) error {
	var body storage.AppConfig
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "parse error: " + err.Error()})
	}

	if err := h.configStore.Save(body); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "config saved", "config": body})
}
