package services

import (
	"api-load-tester/config"
	"api-load-tester/models"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LoadTesterService struct {
	Config *config.ApiConfig
}

func NewLoadTesterService(config *config.ApiConfig) *LoadTesterService {
	return &LoadTesterService{
		Config: config,
	}
}

func (lt *LoadTesterService) SendRequest(request models.ApiRequest) (int, error) {

	// Convert body string to io.Reader
	bodyReader := strings.NewReader(request.Body)

	// Create a new HTTP request
	req, err := http.NewRequest(request.Method, request.URL, bodyReader)
	if err != nil {
		return -1, err
	}

	// Set custom headers if provided
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	// Create an HTTP client with a timeout

	client := &http.Client{
		Timeout: time.Duration(lt.Config.Timeout) * time.Second, // Make timeout configurable
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close() // Ensure the response body is closed

	// Return the HTTP status code
	return resp.StatusCode, nil
}

func (lt *LoadTesterService) runRequest(wg *sync.WaitGroup) {
	defer wg.Done()

	request := models.ApiRequest{
		URL:     lt.Config.URL,
		Method:  lt.Config.Method,
		Headers: lt.Config.Headers,
		Body:    lt.Config.Body,
	}

	// Send the request
	statusCode, err := lt.SendRequest(request)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Printf("Response Status: %d\n", statusCode)
	}
}

func (lt *LoadTesterService) RunLoadTest() {
	var wg sync.WaitGroup
	startTime := time.Now()

	// Calculate the number of batches needed
	batches := lt.Config.TotalRequests / lt.Config.Concurrency
	remainingRequests := lt.Config.TotalRequests % lt.Config.Concurrency

	// Send requests in batches
	for batch := 0; batch < batches; batch++ {
		log.Printf("Starting Batch %d\n", batch+1)

		wg.Add(lt.Config.Concurrency)
		for i := 0; i < lt.Config.Concurrency; i++ {
			go lt.runRequest(&wg)
		}
	}

	// Send remaining requests (if any)
	if remainingRequests > 0 {
		wg.Add(remainingRequests)
		for i := 0; i < remainingRequests; i++ {
			go lt.runRequest(&wg)
		}
	}

	wg.Wait() // Wait for all requests to complete

	// Print the total time taken for the load test
	elapsedTime := time.Since(startTime)
	log.Printf("Load test completed in %v\n", elapsedTime)
}
