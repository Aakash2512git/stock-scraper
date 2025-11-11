package models

type CompanyInfo struct {
	Company       string `json:"company"`
	MarketCap     string `json:"market_cap"`
	CurrentPrice  string `json:"current_price"`
	HighLow       string `json:"high_low"`
	StockPE       string `json:"stock_pe"`
	BookValue     string `json:"book_value"`
	DividendYield string `json:"dividend_yield"`
	ROCE          string `json:"roce"`
	ROE           string `json:"roe"`
	FaceValue     string `json:"face_value"`
}

type CompanyRequest struct {
	Company string `json:"company" example:"RELIANCE"`
}
