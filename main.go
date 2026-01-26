package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

const (
	processedFile = "processed_urls.json"
	resultsFile   = "results.json"
)

func main() {
	// 1. Load already processed URLs
	processedURLs, err := loadProcessedURLs(processedFile)
	if err != nil {
		log.Fatalf("could not load processed URLs: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		JavaScriptEnabled: playwright.Bool(false),
		UserAgent:         playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		ExtraHttpHeaders: map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"Accept-Language":           "hr-HR,hr;q=0.9,en-US;q=0.8,en;q=0.7",
			"Cache-Control":             "max-age=0",
			"Sec-Ch-Ua":                 `"Chromium";v="122", "Not(A:Brand";v="24", "Google Chrome";v="122"`,
			"Sec-Ch-Ua-Mobile":          "?0",
			"Sec-Ch-Ua-Platform":        `"Windows"`,
			"Sec-Fetch-Dest":            "document",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-Site":            "none",
			"Sec-Fetch-User":            "?1",
			"Upgrade-Insecure-Requests": "1",
		},
	})
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}
	defer context.Close()

	currentSearchURL := "https://www.njuskalo.hr/search/?keywords=vw+tiguan&showAllCategories=1&price[min]=10000&price[max]=15000&condition[used]=1&adsWithImages=1"

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	var newAds []AdInfo
	seenURLs := make(map[string]bool)
	pageNum := 1

	// 2. Collect all ads from all pages, skipping already processed ones
	for currentSearchURL != "" {
		fmt.Printf("Navigating to search page %d: %s\n", pageNum, currentSearchURL)
		ads, nextURL, err := scrapeSearchPage(page, currentSearchURL)
		if err != nil {
			fmt.Printf("Error scraping search page %d: %v\n", pageNum, err)
			break
		}

		for _, ad := range ads {
			if !seenURLs[ad.URL] {
				seenURLs[ad.URL] = true
				if !processedURLs[ad.URL] {
					newAds = append(newAds, ad)
				}
			}
		}
		fmt.Printf("Found %d new ads so far.\n", len(newAds))

		if nextURL == "" {
			break
		}

		currentSearchURL = nextURL
		pageNum++
		time.Sleep(2 * time.Second)
	}

	if len(newAds) == 0 {
		fmt.Println("No new ads found. Exiting.")
		return
	}

	fmt.Printf("\nProcessing %d new ads...\n", len(newAds))

	var newResults []CarDetail
	for i, ad := range newAds {
		fmt.Printf("\r\033[K[%d/%d] Scraping Detail: %s", i+1, len(newAds), ad.URL)
		detail, err := scrapeDetail(page, ad.URL)

		if err != nil {
			fmt.Printf("Error scraping %s: %v\n", ad.URL, err)
			continue
		}
		if detail.Price == "" {
			detail.Price = ad.Price
		}
		newResults = append(newResults, detail)
		processedURLs[ad.URL] = true

		// Save progress periodically
		if (i+1)%5 == 0 {
			saveProcessedURLs(processedFile, processedURLs)
		}

		time.Sleep(time.Duration(1500+time.Now().UnixNano()%1500) * time.Millisecond)
	}
	fmt.Println() // Newline after progress loop

	// 3. Filter Results
	filteredResults := FilterCars(newResults)

	// 4. Save Results
	if err := saveResults(resultsFile, filteredResults); err != nil {
		log.Printf("Error saving results: %v", err)
	}

	// 5. Update Processed URLs
	if err := saveProcessedURLs(processedFile, processedURLs); err != nil {
		log.Printf("Error updating processed URLs: %v", err)
	}

	// 6. Print Results as JSON
	if len(filteredResults) > 0 {
		fmt.Println("\nFiltered Results:")
		finalJSON, _ := json.MarshalIndent(filteredResults, "", "  ")
		fmt.Println(string(finalJSON))
	} else {
		fmt.Println("\nNo cars matched your filters in this session.")
	}

	fmt.Println("\nScraping session complete!")
}

func scrapeSearchPage(page playwright.Page, url string) ([]AdInfo, string, error) {
	if _, err := page.Goto(url); err != nil {
		return nil, "", err
	}

	var ads []AdInfo
	items, err := page.QuerySelectorAll(".EntityList-item--Regular")
	if err != nil {
		return nil, "", err
	}

	for _, item := range items {
		linkEl, _ := item.QuerySelector(".entity-title a")
		priceEl, _ := item.QuerySelector(".price--eur, .price--hrk, .entity-price .price")

		if linkEl != nil {
			href, _ := linkEl.GetAttribute("href")
			if href != "" {
				if !strings.HasPrefix(href, "http") {
					href = "https://www.njuskalo.hr" + href
				}

				price := ""
				if priceEl != nil {
					price, _ = priceEl.TextContent()
					price = strings.TrimSpace(price)
				}

				ads = append(ads, AdInfo{URL: href, Price: price})
			}
		}
	}

	nextPageEl, _ := page.QuerySelector(".Pagination-item--next a")
	nextURL := ""
	if nextPageEl != nil {
		href, _ := nextPageEl.GetAttribute("href")
		if href != "" {
			if !strings.HasPrefix(href, "http") {
				nextURL = "https://www.njuskalo.hr" + href
			} else {
				nextURL = href
			}
		}
	}

	return ads, nextURL, nil
}

func scrapeDetail(page playwright.Page, url string) (CarDetail, error) {
	if _, err := page.Goto(url); err != nil {
		return CarDetail{}, err
	}

	content, _ := page.Content()
	if strings.Contains(content, "ShieldSquare Captcha") || strings.Contains(content, "Robot ne smije naškoditi čovjeku") {
		return CarDetail{}, fmt.Errorf("blocked by captcha")
	}

	detail := CarDetail{URL: url}

	titleEl, _ := page.QuerySelector(".ClassifiedDetailSummary-title")
	if titleEl != nil {
		detail.Title, _ = titleEl.TextContent()
		detail.Title = strings.TrimSpace(detail.Title)
	}

	priceEl, _ := page.QuerySelector(".price--hrk, .price--eur, .ClassifiedDetailSummary-price--eur")
	if priceEl != nil {
		detail.Price, _ = priceEl.TextContent()
		detail.Price = strings.TrimSpace(detail.Price)
	}

	locEl, _ := page.QuerySelector(".ClassifiedDetailSummary-address, .entity-description-item--location, .entity-description-main .entity-description-item")
	if locEl != nil {
		detail.Location, _ = locEl.TextContent()
		detail.Location = strings.TrimSpace(detail.Location)
		if strings.Contains(detail.Location, "Lokacija vozila:") {
			detail.Location = strings.TrimSpace(strings.Replace(detail.Location, "Lokacija vozila:", "", 1))
		}
	}

	items, _ := page.QuerySelectorAll(".ClassifiedDetailBasicDetails-list dt, .ClassifiedDetailBasicDetails-list dd")
	for i := 0; i < len(items)-1; i += 2 {
		label, _ := items[i].TextContent()
		value, _ := items[i+1].TextContent()
		label, value = strings.TrimSpace(label), strings.TrimSpace(value)

		switch label {
		case "Lokacija vozila":
			detail.Location = value
		case "Marka automobila":
			detail.Brand = value
		case "Model automobila":
			detail.Model = value
		case "Tip automobila":
			detail.Type = value
		case "Godina proizvodnje":
			detail.Year = value
		case "Godina modela":
			detail.ModelYear = value
		case "Prijeđeni kilometri":
			detail.Mileage = value
		case "Motor":
			detail.Engine = value
		case "Snaga motora":
			detail.Power = value
		case "Radni obujam":
			detail.Displacement = value
		case "Mjenjač":
			detail.Gearbox = value
		case "Broj stupnjeva":
			detail.Gears = value
		case "Stanje":
			detail.Condition = value
		case "Servisna knjiga":
			detail.ServiceBook = value
		}
	}

	return detail, nil
}
