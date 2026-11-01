//go:build ignore

// Read the NetAPI compromised IP feed with threat data (JSON Lines).
//
// The addresses were caught attacking NetAPI Nightwatch, NetAPI's own honeypot
// network. The feed is free, needs no token and is rebuilt every hour. Each line
// is one JSON object: ip, first_seen, last_seen, category (A infection attempt,
// B password guessing, C exploit), top_threat, threats, country, asn, as_name,
// sensors.
//
// The response is streamed line by line, nothing is buffered in memory.
//
// Standard library only. Usage: go run go-example-compromised-ips.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
)

const apiURL = "https://netapi.com/api2/"

type threat struct {
	Type   string   `json:"type"`
	Name   string   `json:"name"`
	Family string   `json:"family"`
	CVE    string   `json:"cve"`
	Mitre  []string `json:"mitre"`
	First  string   `json:"first"`
	Last   string   `json:"last"`
	Days   int      `json:"days"`
}

type record struct {
	IP        string   `json:"ip"`
	Category  string   `json:"category"`
	TopThreat string   `json:"top_threat"`
	FirstSeen string   `json:"first_seen"`
	LastSeen  string   `json:"last_seen"`
	Country   string   `json:"country"`
	ASN       int      `json:"asn"`
	ASName    string   `json:"as_name"`
	Sensors   int      `json:"sensors"`
	Threats   []threat `json:"threats"`
}

type pair struct {
	key   string
	count int
}

func top(counter map[string]int, n int) []pair {
	pairs := make([]pair, 0, len(counter))
	for k, v := range counter {
		pairs = append(pairs, pair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].count > pairs[j].count })
	if len(pairs) > n {
		pairs = pairs[:n]
	}
	return pairs
}

func main() {
	params := url.Values{}
	params.Set("method", "compromised")
	params.Set("dataset_type", "ip") // "ip" = last 24 hours, "ip-all" = every address ever listed
	params.Set("format", "jsonl")
	// Optional filters, the same as the list builder on https://netapi.com/compromised-ips/
	// params.Set("category", "B,C")
	// params.Set("country_not", "us")
	// params.Set("min_sensors", "2")

	resp, err := http.Get(apiURL + "?" + params.Encode())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// API errors are plain text, e.g. "405 Method Not Allowed: Invalid category: ..."
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "HTTP %d: %s\n", resp.StatusCode, body)
		os.Exit(1)
	}

	categories := map[string]int{}
	threats := map[string]int{}
	countries := map[string]int{}
	listed := map[string]struct{}{}
	first := ""

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var r record
		if err := json.Unmarshal(line, &r); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if first == "" {
			first = r.IP
		}
		listed[r.IP] = struct{}{}
		categories[r.Category]++
		country := r.Country
		if country == "" {
			country = "??"
		}
		countries[country]++
		for _, t := range r.Threats {
			threats[t.Type]++
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%d addresses in the list\n", len(listed))
	fmt.Println("By category:", categories)
	fmt.Println("Top attack types:", top(threats, 5))
	fmt.Println("Top countries:", top(countries, 5))

	for _, ip := range []string{"8.8.8.8", first} {
		if ip == "" {
			continue
		}
		if _, ok := listed[ip]; ok {
			fmt.Println(ip, "-> LISTED")
		} else {
			fmt.Println(ip, "-> not listed")
		}
	}
}
