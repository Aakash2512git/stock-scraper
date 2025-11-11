# 🐿️ Stock Scraper (Go + Redis + Docker)

A simple **Stock Scraper** built in **Go (Golang)** that scrapes stock data for various companies and stores it in **Redis**.  
A background job runs every 24 hours to update the data automatically.  
This project was built purely for **hands-on learning and experimentation** with Go, Redis, and Docker.

---

## 🚀 Features

- 🔄 Scrapes stock data for multiple companies  
- 🕒 Scheduled job runs every 24 hours to refresh data  
- 💾 Stores scraped data in Redis  
- 🐳 Docker Compose setup for Redis container  
- 🧠 Focused on learning Go fundamentals: concurrency, jobs, and data persistence

---

## 🧰 Tech Stack

- **Language:** Go (Golang)
- **Database/Cache:** Redis
- **Containerization:** Docker & Docker Compose
- **Scheduler:** Custom Go job (runs every 24h)
- **Libraries:** 
  - `net/http` for scraping
  - `encoding/json` for data handling
  - `go-redis/redis/v8` for Redis integration
  - `time` and `cron` (optional) for scheduling

---

## 🛠️ Setup & Run

### 1️⃣ Clone the repository
```bash
git clone https://github.com/<your-username>/stock-scraper.git
cd stock-scraper

2️⃣ Run via Docker Compose
 docker-compose up -d

 Run the Go scraper
 go run main.go


