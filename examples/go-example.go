//go:build ignore

// Download a zone domain list from NetAPI and stream it line by line.
//
// The response is a gzip-compressed CSV. It is decompressed on the fly, so
// even the .com list can be processed without buffering it in memory.
//
// Standard library only. Usage: go run go-example.go
package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const apiURL = "https://netapi.com/api2/"
const apiToken = "YOUR_API_TOKEN" // https://netapi.com/dashboard/

func main() {
	params := url.Values{}
	params.Set("method", "download")
	params.Set("zone_tld", "net")       // any TLD from ?method=zones, or "all-zones"
	params.Set("dataset_type", "list")  // "list" or "dataset"
	params.Set("filter_type", "active") // "active", "new" or "deleted"
	params.Set("token", apiToken)

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Errors are short plain-text messages, e.g. "403 Forbidden: No active/paid plan."
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "HTTP %d: %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "not a gzip response:", err)
		os.Exit(1)
	}
	defer gz.Close()

	scanner := bufio.NewScanner(gz)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // dataset rows with many emails can be long
	n := 0
	for scanner.Scan() {
		if n == 0 {
			fmt.Println("Header:", scanner.Text())
		} else {
			fmt.Println(scanner.Text())
		}
		n++
		if n > 10 {
			break // remove to process the whole file
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
