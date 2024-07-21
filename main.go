package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func CheckURL(url string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	client := http.Client{
		Timeout: 5, *time.Second,
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
