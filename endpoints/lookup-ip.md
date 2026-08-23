# `lookup-ip` — domains hosted on one IP address

Reverse-IP lookup: returns up to three domains whose websites resolved to the given IP address at the last crawl, with their hostname and nameservers.

- **Authentication:** `token` + active paid plan (any plan)
- **Response:** up to 3 CSV rows, no header; plain text

## Request

```
GET https://netapi.com/api2/?method=lookup-ip&ip={ip}&token={token}
```

| Parameter | Required | Description |
|---|---|---|
| `ip` | yes | IPv4 or IPv6 address, e.g. `172.67.69.160`. Maximum length 99 characters. |
| `token` | yes | Your API token |

## Response

One line per domain, at most three lines:

```
url,hostname,dns1,dns2
```

| Field | Description |
|---|---|
| `url` | Domain hosted on the IP. |
| `hostname` | Web-server hostname recorded for that domain. |
| `dns1`, `dns2` | Authoritative nameservers of the domain; `dns2` is empty if only one is known. |

Example (illustrative values):

```
example.com,web1.example.com,ns1.example.com,ns2.example.com
example-shop.com,web1.example.com,ns1.example.com,ns2.example.com
```

If no domain is known for the IP the response body is the text `NOT FOUND` with HTTP status 200.

## Errors

| HTTP | Message |
|---|---|
| 401 | `401 Unauthorized: Missing user token.` |
| 403 | `403 Forbidden: Incorrect token.` |
| 403 | `403 Forbidden: No active/paid plan.` |
| 405 | `405 Method Not Allowed: Missing IP.` |
| 405 | `405 Method Not Allowed: IP address is too long.` |

## Examples

cURL:

```bash
curl "https://netapi.com/api2/?method=lookup-ip&ip=172.67.69.160&token=YOUR_API_TOKEN"
```

Python:

```python
import csv
import io
import requests

resp = requests.get(
    "https://netapi.com/api2/",
    params={"method": "lookup-ip", "ip": "172.67.69.160", "token": "YOUR_API_TOKEN"},
    timeout=30,
)
resp.raise_for_status()

if resp.text.strip() == "NOT FOUND":
    print("no domains known for this IP")
else:
    for url, hostname, dns1, dns2 in csv.reader(io.StringIO(resp.text)):
        print(url, hostname, dns1, dns2)
```

## Notes

- The response is capped at three domains; shared hosting and CDN addresses can carry far more. Use the domain datasets ([`download`](download.md), `ip` column) if you need every domain on an address.
- Results reflect the crawler index, not a live PTR query.
- The same lookup is available in the browser at [netapi.com/lookup-ip/](https://netapi.com/lookup-ip/).
- Web reference: [netapi.com/help/api/#lookup-ip-api](https://netapi.com/help/api/#lookup-ip-api).
