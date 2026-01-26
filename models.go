package main

type CarDetail struct {
	URL          string `json:"url"`
	Title        string `json:"title"`
	Price        string `json:"price"`
	Location     string `json:"location"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Type         string `json:"type"`
	Year         string `json:"year"`
	ModelYear    string `json:"model_year"`
	Mileage      string `json:"mileage"`
	Engine       string `json:"engine"`
	Power        string `json:"power"`
	Displacement string `json:"displacement"`
	Gearbox      string `json:"gearbox"`
	Gears        string `json:"gears"`
	Condition    string `json:"condition"`
	ServiceBook  string `json:"service_book"`
}

type AdInfo struct {
	URL   string
	Price string
}
