# goexpert-stress-test
Exercise of stress test CLI for Go Expert Postgraduate Course

## Description

A CLI application designed to perform load testing on web services. Users can specify the target URL, total number of requests, and the level of concurrency. The system generates a detailed report upon test completion.

## How to Install

### Prerequisites
- Go 1.21 or higher installed on your system

### Installation Steps

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the application:
   ```bash
   make docker-build
   ```

## How to Run

- Example:
   ```bash
   docker run --rm goexpert-stress-test --url=http://httpbin.org/get --requests=100 --concurrency=10
   ```
   This will run 100 requests with 10 concurrent connections to http://httpbin.org/get
