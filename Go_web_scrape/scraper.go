package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"regexp"
	// "github.com/PuerkitoBio/goquery"
)

func main() {
	urlPtr := flag.String("url", "https://example.com", "The URL to scrape")
	flag.Parse()

	// Target URL
	url := *urlPtr
	_, err := neturl.Parse(*urlPtr)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		return
	}
	fmt.Printf("Scraping URL: %s\n", url)

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	bodyString := string(bodyBytes)

	// extract the Title (using Regex)
	// regex is not ideal for HTML
	titleRegex := regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	titleMatch := titleRegex.FindStringSubmatch(bodyString)

	if len(titleMatch) > 1 {
		fmt.Printf("\nPage Title: %s\n", titleMatch[1])
	} else {
		fmt.Println("\nNo title found")
	}

	// extract Links (using Regex)
	// Matches href="url" or href='url'
	linkRegex := regexp.MustCompile(`(?i)href=["'](.*?)["']`)
	links := linkRegex.FindAllStringSubmatch(bodyString, -1)

	fmt.Printf("\nFound %d links:\n", len(links))
	for i, link := range links {
		if i >= 10 {
			fmt.Println("... (more than 10 links found)")
			break
		}
		fmt.Printf("- %s\n", link[1])
	}
}
