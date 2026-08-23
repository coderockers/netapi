//go:build ignore

// Read the NetAPI compromised-domain / compromised-IP feed.
//
// The feed is free and needs no token. The first line is a "#" comment with
// the list name and date; every other line is one domain or IP address.
//
// Standard library only. Usage: go run go-example-compromised.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const apiURL = "https://netapi.com/api2/"

func main() {
	params := url.Values{}
	params.Set("method", "compromised")
	params.Set("dataset_type", "url") // "url", "ip", "url-all" or "ip-all"

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "HTTP %d: %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	blocked := make(map[string]struct{})
	var first []string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		blocked[line] = struct{}{}
		if len(first) < 10 {
			first = append(first, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%d entries in the current list\n", len(blocked))
	fmt.Println("First 10:", first)
	if len(first) == 0 {
		return
	}

	for _, domain := range []string{"example.com", first[0]} {
		if _, ok := blocked[domain]; ok {
			fmt.Println(domain, "-> LISTED")
		} else {
			fmt.Println(domain, "-> not listed")
		}
	}
}
