# `download-ip` — reverse-IP dataset (IP → domains)

Downloads the complete reverse-IP dataset: every web-server IP address NetAPI has seen, with the list of domains hosted on it. One file covers all zones.

- **Authentication:** `token` + active paid plan (any plan)
- **Response:** gzip-compressed CSV file (`.csv.gz`) with a header row; plain CSV with `format=plain`

## Request

```
GET https://netapi.com/api2/?method=download-ip&token={token}[&format=plain]
```

| Parameter | Required | Description |
|---|---|---|
| `token` | yes | Your API token |
| `format` | no | `plain` — uncompressed CSV (several gigabytes) |

## Response

Filename `ip2domain.csv.gz` (or `ip2domain.csv`). Two columns:

```
IP,domains
```

| Column | Description |
|---|---|
| `IP` | IPv4 or IPv6 address of a web server. |
| `domains` | All domains whose website resolved to that IP at the last crawl, comma-separated inside one double-quoted field. |

Example:

```
IP,domains
203.0.113.10,"example.de,example-shop.de"
198.51.100.7,beispiel.de
2001:db8::1,"ipv6-site.org,another-site.net"
```

## Errors

| HTTP | Message |
|---|---|
| 401 | `401 Unauthorized: Missing user token.` |
| 403 | `403 Forbidden: Incorrect token.` |
| 403 | `403 Forbidden: No active/paid plan.` |
| 404 | `404 Not Found: File not found. …` |

## Examples

cURL:

```bash
curl "https://netapi.com/api2/?method=download-ip&token=YOUR_API_TOKEN" -o ip2domain.csv.gz
```

Python — build an IP → domains index for one /24 network:

```python
import csv
import gzip
import ipaddress
import requests

network = ipaddress.ip_network("203.0.113.0/24")
params = {"method": "download-ip", "token": "YOUR_API_TOKEN"}

with requests.get("https://netapi.com/api2/", params=params, stream=True, timeout=1800) as resp:
    resp.raise_for_status()
    with gzip.open(resp.raw, mode="rt", encoding="utf-8", newline="") as f:
        for row in csv.DictReader(f):
            try:
                ip = ipaddress.ip_address(row["IP"])
            except ValueError:
                continue
            if ip in network:
                print(ip, row["domains"].split(","))
```

## Notes

- For a single IP, use [`lookup-ip`](lookup-ip.md) instead of downloading the whole file.
- The dataset is rebuilt daily from the crawler index; an IP is listed only when at least one domain's website currently resolves to it.
- Web reference: [netapi.com/help/api/#ip2domains-api](https://netapi.com/help/api/#ip2domains-api).
