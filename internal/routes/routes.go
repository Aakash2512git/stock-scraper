// routes/routes.go
package routes

import (
	company_handler "main/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetUpRoutes() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	api := app.Group("/api")
	api.Get("/companies/refresh", company_handler.RefreshCompanies)
	api.Get("/companies", company_handler.GetFullCompany)
	api.Get("/companies/:name", company_handler.GetOneCompany)

	return app
}
