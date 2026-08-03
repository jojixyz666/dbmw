package handlers

import (
	"net/http"

	"dbmw/core/project"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	projectSvc *project.Service
}

func NewProjectHandler(projectSvc *project.Service) *ProjectHandler {
	return &ProjectHandler{projectSvc: projectSvc}
}

func (h *ProjectHandler) Detect(c *fiber.Ctx) error {
	path := c.Query("path", ".")
	info, err := h.projectSvc.DetectFramework(path)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(info)
}

func (h *ProjectHandler) GenerateConfig(c *fiber.Ctx) error {
	var body struct {
		Path   string             `json:"path"`
		Config project.DBMWConfig `json:"config"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid format: " + err.Error()})
	}

	if body.Path == "" {
		body.Path = "."
	}

	createdPath, err := h.projectSvc.GenerateConfig(body.Path, body.Config)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "dbmw.yml generated successfully",
		"path":    createdPath,
	})
}
