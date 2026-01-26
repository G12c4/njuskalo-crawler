package main

import (
	"encoding/json"
	"os"
)

func loadProcessedURLs(filename string) (map[string]bool, error) {
	processed := make(map[string]bool)
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return processed, nil
		}
		return nil, err
	}

	var urls []string
	if err := json.Unmarshal(data, &urls); err != nil {
		return nil, err
	}

	for _, url := range urls {
		processed[url] = true
	}
	return processed, nil
}

func saveProcessedURLs(filename string, processed map[string]bool) error {
	var urls []string
	for url := range processed {
		urls = append(urls, url)
	}
	data, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func saveResults(filename string, results []CarDetail) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
