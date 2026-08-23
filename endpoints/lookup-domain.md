# `lookup-domain` — details of one domain

Returns what NetAPI knows about a single registered domain: nameservers, web-server hostname, IP and country, registration and expiration dates, and the registrar.

- **Authentication:** `token` + active paid plan (any plan)
- **Response:** one CSV row, no header; plain text

## Request

```
GET https://netapi.com/api2/?method=lookup-domain&domain={domain}&token={token}
```

| Parameter | Required | Description |
|---|---|---|
| `domain` | yes | Registered domain name, e.g. `example.com`. Pass the bare domain — no scheme, path or subdomain (anything else is looked up literally and returns `NOT FOUND`). The last label must be a zone listed in [`zones`](zones.md). Maximum length 199 characters. |
| `token` | yes | Your API token |

## Response

A single line with nine comma-separated fields:

```
url,dns1,dns2,hostname,ip,ip_country,registered_at,expiring_at,registrar_id
```

| Field | Description |
|---|---|
| `url` | The requested domain. |
| `dns1`, `dns2` | Authoritative nameservers; `dns2` is empty if only one is known. |
| `hostname` | Hostname of the web server. |
| `ip` | IP address of the web server. |
| `ip_country` | Two-letter country code of the server IP. |
| `registered_at` | Registration date, `YYYY-MM-DD`; empty if unknown. |
| `expiring_at` | Expiration date, `YYYY-MM-DD`; empty if unknown. |
| `registrar_id` | Numeric registrar ID; resolve it with [`registrars`](registrars.md). Empty if unknown. |

Example (illustrative values):

```
example.com,ns1.example.com,ns2.example.com,web1.example.com,203.0.113.10,US,2015-08-14,2027-08-13,146
```

If the domain is not in the database the response body is the text `NOT FOUND` with HTTP status 200.

## Errors

| HTTP | Message |
|---|---|
| 401 | `401 Unauthorized: Missing user token.` |
| 403 | `403 Forbidden: Incorrect token.` |
| 403 | `403 Forbidden: No active/paid plan.` |
| 405 | `405 Method Not Allowed: Missing domain name.` |
| 405 | `405 Method Not Allowed: Invalid domain zone.` — the TLD is not a supported zone |
| 405 | `405 Method Not Allowed: Domain name is too long.` |

## Examples

cURL:

```bash
curl "https://netapi.com/api2/?method=lookup-domain&domain=example.com&token=YOUR_API_TOKEN"
```

Python:

```python
import csv
import requests

FIELDS = ["url", "dns1", "dns2", "hostname", "ip", "ip_country",
          "registered_at", "expiring_at", "registrar_id"]

resp = requests.get(
    "https://netapi.com/api2/",
    params={"method": "lookup-domain", "domain": "example.com", "token": "YOUR_API_TOKEN"},
    timeout=30,
)
resp.raise_for_status()

if resp.text.strip() == "NOT FOUND":
    print("not in database")
else:
    row = next(csv.reader([resp.text.strip()]))
    print(dict(zip(FIELDS, row)))
```

## Notes

- One domain per request; there is no bulk variant. For many domains use the [`download`](download.md) datasets.
- Values reflect the last crawl of the domain, not a live DNS query.
- The same lookup is available in the browser at [netapi.com/lookup-domain/](https://netapi.com/lookup-domain/).
- Web reference: [netapi.com/help/api/#lookup-domain-api](https://netapi.com/help/api/#lookup-domain-api).
