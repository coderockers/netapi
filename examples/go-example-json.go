//go:build ignore

// NetAPI JSON API (/api-json/): JSON answers for integrations and AI tools.
// Free methods need no token; the paid ones (and higher limits) take the token as a Bearer header.
// Standard library only: go run go-example-json.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	token = "YOUR_API_TOKEN"
	base  = "https://netapi.com/api-json/"
)

type apiError struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		RetryAfter int    `json:"retry_after"`
	} `json:"error"`
}

// call runs one method and decodes the JSON answer into out; withToken adds the Bearer header.
func call(method string, params map[string]string, withToken bool, out any) error {
	query := url.Values{"method": {method}}
	for key, value := range params {
		query.Set(key, value)
	}
	req, err := http.NewRequest(http.MethodGet, base+"?"+query.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if withToken {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var e apiError // errors are JSON too
		_ = json.NewDecoder(resp.Body).Decode(&e)
		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Fprintf(os.Stderr, "rate limited, retry after %d s\n", e.Error.RetryAfter)
		}
		return fmt.Errorf("%d %s: %s", resp.StatusCode, e.Error.Code, e.Error.Message)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func main() {
	// 1. compromised check (free)
	var check struct {
		ListedNow  bool `json:"listed_now"`
		ListedEver bool `json:"listed_ever"`
	}
	if err := call("compromised-check", map[string]string{"value": "example.com"}, false, &check); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("example.com listed now:", check.ListedNow, "| ever:", check.ListedEver)

	// 2. newly registered domains containing "crypto" (free)
	var fresh struct {
		Results []struct {
			Domain  string `json:"domain"`
			AddedOn string `json:"added_on"`
		} `json:"results"`
		Truncated bool `json:"truncated"`
	}
	if err := call("search-new", map[string]string{"q": "crypto", "days": "3", "limit": "5"}, false, &fresh); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, row := range fresh.Results {
		fmt.Println(row.AddedOn, row.Domain)
	}

	// 3. domain lookup (active plan required)
	var lookup struct {
		DNS       []string `json:"dns"`
		ExpiresOn string   `json:"expires_on"`
	}
	if err := call("lookup-domain", map[string]string{"domain": "example.com"}, true, &lookup); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("nameservers:", lookup.DNS, "| expires:", lookup.ExpiresOn)
}
