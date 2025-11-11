package company_handler

import (
	"context"
	"fmt"
	"main/internal/models"
	"main/internal/redis"
	"main/internal/scraper"

	"github.com/gofiber/fiber/v2"
)

// InjestData saves scraped companies into Redis
func InjestData(ctx context.Context, companies []models.CompanyInfo) error {
	if err := redis.SaveCompanies(ctx, companies); err != nil {
		return err
	}
	fmt.Println("✅ Companies saved to Redis")
	return nil
}

// RefreshCompanies godoc
// @Summary Referesh all companies
// @Description Returns a list of companies scraped from screener.in
// @Tags companies
// @Produce json
// @Success 200 {array} models.CompanyInfo
// @Router /api/companies/refresh [get]
func RefreshCompanies(c *fiber.Ctx) error {
	ctx := context.Background()

	// Scrape companies
	companies, err := scraper.ScrapeCompanies()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to scrape companies",
		})
	}

	// Save in Redis
	if err := InjestData(ctx, companies); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to save companies in Redis",
		})
	}

	return c.JSON(fiber.Map{
		"message":   "✅ Companies refreshed successfully",
		"companies": companies,
	})
}

// GetCompanies godoc
// @Summary Get a company details
// @Description Returns info about a specific company scraped from screener.in
// @Tags companies
// @Produce json
// @Param name path string true "Company name"
// @Success 200 {object} models.CompanyInfo
// @Router /api/companies/{name} [get]
func GetOneCompany(c *fiber.Ctx) error {
	ctx := context.Background()
	name := c.Params("name") // from URL like /api/companies/RELIANCE

	company, err := redis.GetCompanyInfo(ctx, name)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "company not found",
		})
	}

	return c.JSON(company)
}

// GetCompanies godoc
// @Summary Get all companies
// @Description Returns a list of companies scraped from screener.in
// @Tags companies
// @Produce json
// @Success 200 {array} models.CompanyInfo
// @Router /api/companies [get]
func GetFullCompany(c *fiber.Ctx) error {
	ctx := context.Background()

	companies, err := redis.GetAllCompanies(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch companies",
		})
	}

	return c.JSON(companies)
}
