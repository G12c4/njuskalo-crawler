package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// FilterCars applies business logic to filter results before saving or notifying
func FilterCars(cars []CarDetail) []CarDetail {
	fmt.Printf("Applying filters to %d cars...\n", len(cars))

	reasons := make(map[string]int)
	filtered := []CarDetail{}
	for _, car := range cars {
		// 1. Brand: VW
		if !isVW(car.Brand) {
			reasons["Brand: "+car.Brand]++
			continue
		}

		// 2. Automatic transmission
		if !isAutomatic(car.Gearbox) {
			reasons["Gearbox: "+car.Gearbox]++
			continue
		}

		// 3. Year: 2014 - 2016
		year := parseYear(car.Year)
		if year < 2014 || year > 2016 {
			reasons[fmt.Sprintf("Year: %d", year)]++
			continue
		}

		// 4. Price: 10,000 - 15,000 EUR
		price := parsePrice(car.Price)
		if price < 10000 || price > 15000 {
			reasons[fmt.Sprintf("Price: %d", price)]++
			continue
		}

		// Transform Power kW to HP
		car.Power = convertKWtoHP(car.Power)

		filtered = append(filtered, car)
	}

	rejectedCount := len(cars) - len(filtered)
	if rejectedCount > 0 {
		reasonStrings := []string{}
		for r, count := range reasons {
			reasonStrings = append(reasonStrings, fmt.Sprintf("%s (%d)", r, count))
		}
		fmt.Printf("Rejected %d/%d reasons [%s]\n", rejectedCount, len(cars), strings.Join(reasonStrings, ", "))
	}

	fmt.Printf("Filters applied: %d cars remaining.\n", len(filtered))
	return filtered
}

func convertKWtoHP(powerStr string) string {
	// Extracts digits from "110 kW"
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(powerStr)
	if match == "" {
		return powerStr
	}

	kw, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return powerStr
	}

	// 1 kW = 1.35962 HP
	hp := math.Round(kw * 1.35962)
	return fmt.Sprintf("%d HP", int(hp))
}

func isVW(brand string) bool {
	b := strings.ToLower(brand)
	return strings.Contains(b, "vw") || strings.Contains(b, "volkswagen")
}

func isAutomatic(gearbox string) bool {
	g := strings.ToLower(gearbox)
	return strings.Contains(g, "automatski") || strings.Contains(g, "dsg") || strings.Contains(g, "triptronic") || strings.Contains(g, "sekvencijski")
}

func parseYear(yearStr string) int {
	re := regexp.MustCompile(`\d{4}`)
	match := re.FindString(yearStr)
	if match == "" {
		return 0
	}
	year, _ := strconv.Atoi(match)
	return year
}

func parsePrice(priceStr string) int {
	re := regexp.MustCompile(`[^\d]`)
	clean := re.ReplaceAllString(priceStr, "")
	if clean == "" {
		return 0
	}
	price, _ := strconv.Atoi(clean)
	return price
}
