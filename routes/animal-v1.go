package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rafaeldajuda/entity"
)

func Routes(app *fiber.App, api *entity.API) {
	routeV1(app, api)
}

func routeV1(app *fiber.App, api *entity.API) {
	v1 := app.Group("/api/v1")

	v1.Get("/", api.HandlerV1.GetAnimals)
	v1.Get("/:id", api.HandlerV1.GetAnimal)
	v1.Post("/", api.HandlerV1.PostAnimal)
	v1.Put("/:id", api.HandlerV1.PutAnimal)
	v1.Delete("/:id", api.HandlerV1.DeleteAnimal)
}
