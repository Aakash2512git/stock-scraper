package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/models"
	"main/internal/redis"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func ScrapeCompanies() ([]models.CompanyInfo, error) {
	companies := []string{"RELIANCE", "TCS"} // add more companies here

	companyInfo := models.CompanyInfo{}
	companyInfos := make([]models.CompanyInfo, 0, len(companies))

	c := colly.NewCollector(
		colly.AllowedDomains("screener.in", "www.screener.in"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	// Scrape ratios

	c.OnHTML("ul#top-ratios li", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText("span.name"))
		value := strings.TrimSpace(e.ChildText("span.value"))

		switch name {
		case "Market Cap":
			companyInfo.MarketCap = cleanText(value)
		case "Current Price":
			companyInfo.CurrentPrice = cleanText(value)
		case "High / Low":
			companyInfo.HighLow = cleanText(value)
		case "Stock P/E":
			companyInfo.StockPE = cleanText(value)
		case "Book Value":
			companyInfo.BookValue = cleanText(value)
		case "Dividend Yield":
			companyInfo.DividendYield = cleanText(value)
		case "ROCE":
			companyInfo.ROCE = cleanText(value)
		case "ROE":
			companyInfo.ROE = cleanText(value)
		case "Face Value":
			companyInfo.FaceValue = cleanText(value)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		companyInfos = append(companyInfos, companyInfo)
		companyInfo = models.CompanyInfo{} //reset for new company
	})

	// Visit each company page
	for _, company := range companies {
		companyInfo.Company = company
		c.Visit(scrapeUrl(company))
	}

	// Print results as JSON
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(companyInfos)

	return companyInfos, nil
}

func scrapeUrl(company string) string {
	return fmt.Sprintf("https://www.screener.in/company/%s/consolidated/", company)
}
func cleanText(s string) string {
	// Remove \n, \t and extra spaces
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func Run(ctx context.Context, cronTime string) {
	// Simple version: interpret cronTime as duration ("24h")
	d, err := time.ParseDuration(cronTime)
	if err != nil {
		log.Printf("Invalid duration %s, defaulting to 24h", cronTime)
		d = 24 * time.Hour
	}

	// Run immediately
	runOnce(ctx)

	// Run periodically
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scraper stopped")
			return
		case <-ticker.C:
			runOnce(ctx)
		}
	}
}

func runOnce(ctx context.Context) {
	log.Println("Running scraper...")
	companies, err := ScrapeCompanies()
	if err != nil {
		log.Println("Error scraping:", err)
		return
	}
	redis.SaveCompanies(ctx, companies)
	log.Printf("Scraped %d companies\n", len(companies))
}
