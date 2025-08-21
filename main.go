package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type result struct {
	duration time.Duration
	status   int
	err      error
}

type report struct {
	totalRequests      int32
	successfulRequests int32
	statusCodes        map[int]int32
	totalDuration      time.Duration
	responseTimes      []time.Duration
	errors             []error
}

func main() {
	url := flag.String("url", "", "URL to test")
	totalRequests := flag.Int("requests", 0, "Total number of requests")
	concurrency := flag.Int("concurrency", 1, "Number of concurrent requests")

	flag.Parse()

	if *url == "" || *totalRequests <= 0 || *concurrency <= 0 {
		log.Fatal("Missing required parameters. Usage: --url=<url> --requests=<number> --concurrency=<number>")
	}

	runStress(*url, *totalRequests, *concurrency)
}

func runStress(url string, totalRequests, concurrency int) {
	startTime := time.Now()
	results := make(chan result, totalRequests)
	var wg sync.WaitGroup
	report := &report{
		statusCodes:   make(map[int]int32),
		responseTimes: make([]time.Duration, 0, totalRequests),
		errors:        make([]error, 0),
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			for j := 0; j < totalRequests/concurrency; j++ {
				makeRequest(client, url, results)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		report.totalRequests++

		if res.status == http.StatusOK {
			report.successfulRequests++
		}

		if res.status > 0 {
			report.statusCodes[res.status]++
		}

		if res.err != nil {
			report.errors = append(report.errors, res.err)
		} else {
			report.responseTimes = append(report.responseTimes, res.duration)
		}
	}

	report.totalDuration = time.Since(startTime)
	printReport(report)
}

func makeRequest(client *http.Client, url string, results chan<- result) {
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	status := 0
	if resp != nil {
		status = resp.StatusCode
		resp.Body.Close()
	}

	results <- result{
		duration: duration,
		status:   status,
		err:      err,
	}
}

func getPercentile(times []time.Duration, percentile int) time.Duration {
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })

	if len(times) == 0 {
		return 0
	}

	index := (percentile * len(times)) / 100

	if index >= len(times) {
		index = len(times) - 1
	}

	return times[index]
}

func printReport(r *report) {
	fmt.Println("\n=== STRESS TEST REPORT ===")
	fmt.Printf("Total time: %v\n", r.totalDuration.Round(time.Millisecond))
	fmt.Printf("Total requests: %d\n", r.totalRequests)
	fmt.Printf("Successful requests (HTTP 200): %d\n", r.successfulRequests)

	if len(r.errors) > 0 {
		fmt.Printf("\nErrors: %d\n", len(r.errors))
	}

	fmt.Println("\nStatus code distribution:")
	for status, count := range r.statusCodes {
		fmt.Printf("  HTTP %d: %d requests\n", status, count)
	}

	if len(r.responseTimes) > 0 {
		sort.Slice(r.responseTimes, func(i, j int) bool {
			return r.responseTimes[i] < r.responseTimes[j]
		})

		fmt.Println("\nLatency percentiles:")
		fmt.Printf("  p50: %v\n", getPercentile(r.responseTimes, 50).Round(time.Millisecond))
		fmt.Printf("  p75: %v\n", getPercentile(r.responseTimes, 75).Round(time.Millisecond))
		fmt.Printf("  p90: %v\n", getPercentile(r.responseTimes, 90).Round(time.Millisecond))
		fmt.Printf("  p99: %v\n", getPercentile(r.responseTimes, 99).Round(time.Millisecond))

		rps := float64(len(r.responseTimes)) / r.totalDuration.Seconds()
		fmt.Printf("\nRequests per second: %.2f\n", rps)
	}
}
