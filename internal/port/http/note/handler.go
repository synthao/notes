package note

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/synthao/notes/internal/domain"
	"net/http"
	"strconv"
	"time"
)

type CreateRequest struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type GetListResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetOneResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Handler struct {
	app     *fiber.App
	service domain.Service
}

func NewHandler(app *fiber.App, service domain.Service) *Handler {
	return &Handler{app: app, service: service}
}

func (h *Handler) InitRoutes() {
	h.app.Post("/api/notes", func(ctx *fiber.Ctx) error {

		var req CreateRequest

		err := ctx.BodyParser(&req)
		if err != nil {
			return ctx.JSON(fiber.Map{"error": "Failed to create a record. Payload parsing error"})
		}

		id, err := h.service.Create(&domain.Note{Name: req.Name, Text: req.Text})
		if err != nil {
			return ctx.JSON(fiber.Map{"error": "Failed to create a record"})
		}

		return ctx.Status(http.StatusCreated).JSON(fiber.Map{"id": id})
	})

	h.app.Get("/api/notes", func(ctx *fiber.Ctx) error {
		data, err := h.service.GetList(10, 0)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(map[string]any{
				"error": "Something went wrong while fetching data",
			})
		}

		return ctx.JSON(fromDomainToGetListResponse(data))
	})

	h.app.Get("/api/notes/:id<int>", func(ctx *fiber.Ctx) error {
		id, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			return ctx.JSON(fiber.Map{"error": err.Error()})
		}

		data, err := h.service.GetOne(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ctx.SendStatus(http.StatusNotFound)
			}

			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Something went wrong while fetching record",
			})
		}

		return ctx.JSON(fromDomainToGetOneResponse(data))
	})

	h.app.Delete("/api/notes/:id<int>", func(ctx *fiber.Ctx) error {
		id, err := strconv.Atoi(ctx.Params("id"))
		if err != nil {
			return err
		}

		if err := h.service.Delete(id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ctx.SendStatus(http.StatusNotFound)
			}
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return ctx.SendStatus(http.StatusNoContent)
	})
}

func fromDomainToGetOneResponse(data *domain.Note) *GetOneResponse {
	return &GetOneResponse{
		ID:        data.ID,
		Name:      data.Name,
		Text:      data.Text,
		CreatedAt: data.CreatedAt,
	}
}

func fromDomainToGetListResponse(data []domain.Note) []GetListResponse {
	res := make([]GetListResponse, len(data))

	for i, item := range data {
		res[i] = GetListResponse{
			ID:   item.ID,
			Name: item.Name,
		}
	}

	return res
}
