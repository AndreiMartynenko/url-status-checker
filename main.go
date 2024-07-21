package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func CheckURL(url string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		results <- fmt.Sprintf("ERROR: %s -%v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		results <- fmt.Sprintf("OK: %s", url)
	} else {
		results <- fmt.Sprintf("FAIL: %s is Status code: %d", url, resp.StatusCode)
	}
}

func ReadURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	urlMap := make(map[string]bool)
	var urls []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		if url == "" {
			continue // Skip empty lines
		}
		if !urlMap[url] {
			urls = append(urls, url)
			urlMap[url] = true
		} else {
			fmt.Printf("Duplicate URL found and deleted: %s\n", url)
		}
	}
	return urls, scanner.Err()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: url-status-checker <file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	urls, err := ReadURLsFromFile(filename)
	if err != nil {
		fmt.Printf("Error reading file %s\n", err)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	results := make(chan string, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go CheckURL(url, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	logFile, err := os.Create("log.txt")
	if err != nil {
		fmt.Printf("Error creating log file %s\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	for result := range results {
		fmt.Println(result)
		_, err := logFile.WriteString(result + "\n")
		if err != nil {
			fmt.Printf("Error writing to log file %s\n", err)
			os.Exit(1)
		}
	}
}
