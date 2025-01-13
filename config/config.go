package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	URL           string            // API endpoint URL
	Method        string            // HTTP method (GET, POST, etc.)
	Headers       map[string]string // Request headers
	Body          string            // Request body (for POST/PUT)
	Concurrency   int               // Number of concurrent requests
	TotalRequests int               // Total number of requests to send
	Timeout       int               //timeout of the request
}

func LoadConfig() *ApiConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("API_URL")
	if url == "" {
		log.Fatal("API_URL is mandatory")
	}

	method := os.Getenv("API_METHOD")
	if method == "" {
		log.Fatal("API_METHOD is mandatory")
	}

	var headers map[string]string
	headerJson := os.Getenv("API_HEADERS")
	err = json.Unmarshal([]byte(headerJson), &headers)
	if err != nil {
		log.Fatal("Error parsing API_HEADERS:", err)
	}

	concurrencyStr := os.Getenv("API_CONCURRENCY")
	concurrency := 50 // Default value for concurrency
	if concurrencyStr != "" {
		concurrency, err = strconv.Atoi(concurrencyStr)
		if err != nil {
			log.Printf("Error parsing API_CONCURRENCY: %v. Using default value %d.", err, concurrency)
		}
	} else {
		log.Printf("API_CONCURRENCY not set. Using default value %d.", concurrency)
	}

	totalRequestsStr := os.Getenv("API_TOTAL_REQUESTS")
	totalRequests := 100 // Default value for total requests
	if totalRequestsStr != "" {
		totalRequests, err = strconv.Atoi(totalRequestsStr)
		if err != nil {
			log.Printf("Error parsing API_TOTAL_REQUESTS: %v. Using default value %d.", err, totalRequests)
		}
	} else {
		log.Printf("API_TOTAL_REQUESTS not set. Using default value %d.", totalRequests)
	}

	timeoutStr := os.Getenv("TIME_OUT")
	timeout := 10
	if timeoutStr != "" {
		timeout, err := strconv.Atoi(timeoutStr)
		if err != nil {
			log.Printf("Error parsing TIME_OUT: %v. Using default value %d.", err, timeout)
		}
	} else {
		log.Printf("TIME_OUT not set. Using default value %d.", timeout)

	}

	return &ApiConfig{
		URL:           url,
		Method:        method,
		Headers:       headers,
		Body:          os.Getenv("API_BODY"),
		Concurrency:   concurrency,
		TotalRequests: totalRequests,
		Timeout:       timeout,
	}

}
