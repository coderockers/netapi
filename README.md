# NetAPI — Domain Data API (v2)

[NetAPI](https://netapi.com/) provides daily-updated domain data as downloadable CSV files and a small set of lookup endpoints: full domain lists for more than 1,500 TLDs, detailed datasets (nameservers, hosting IP, country, emails, phone numbers, Majestic rank), newly registered domains, domains grouped by DNS provider or registrar, and free threat-intelligence feeds.

This repository is the official companion to the API. It documents every endpoint, ships an OpenAPI specification and contains ready-to-run examples in cURL, Python, Node.js, Go and PHP.

- **Base URL:** `https://netapi.com/api2/`
- **Transport:** HTTPS, `GET` only, parameters in the query string
- **Response format:** CSV. Downloads are gzip-compressed (`.csv.gz`) unless you pass `format=plain`; list and lookup endpoints return plain text.
- **Web documentation:** [netapi.com/help/api/](https://netapi.com/help/api/) (always the most up-to-date reference) · [FAQ](https://netapi.com/help/faq/)

---

## Quick start

```bash
# 1. Which zones are available? (no token needed)
curl "https://netapi.com/api2/?method=zones"

# 2. Download the full .net domain list (requires a token and an active plan)
curl "https://netapi.com/api2/?method=download&zone_tld=net&dataset_type=list&filter_type=active&token=YOUR_API_TOKEN" \
  -o net_active_list.csv.gz
```

Get your API token in the [NetAPI dashboard](https://netapi.com/dashboard/) — sign up with an email address or a Google account, no password required.

---

## Endpoints

All requests go to `https://netapi.com/api2/` and select the operation with the `method` query parameter.

### Reference lists (free, no token)

| Method | Description | Docs |
|---|---|---|
| `zones` | All supported TLDs with daily-update and ccTLD flags | [zones.md](endpoints/zones.md) |
| `dns` | Supported DNS providers (aliases for `download-dns`) | [dns.md](endpoints/dns.md) |
| `registrars` | Known registrars (IDs for `download-whois`) | [registrars.md](endpoints/registrars.md) |

### Downloads (token + active plan)

| Method | Description | Docs |
|---|---|---|
| `download` | Domain list or dataset for one zone or all zones; active, new or deleted domains | [download.md](endpoints/download.md) |
| `download-dns` | Domains that use a given DNS provider | [download-dns.md](endpoints/download-dns.md) |
| `download-whois` | Domains registered with a given registrar, with registration and expiration dates | [download-whois.md](endpoints/download-whois.md) |

### Lookups (token + active plan)

| Method | Description | Docs |
|---|---|---|
| `lookup-domain` | Nameservers, hosting IP, country, registration dates and registrar of one domain | [lookup-domain.md](endpoints/lookup-domain.md) |
| `lookup-ip` | Domains hosted on one IP address (up to 3) | [lookup-ip.md](endpoints/lookup-ip.md) |

### Threat intelligence (free, no token)

| Method | Description | Docs |
|---|---|---|
| `compromised` | Compromised domains or IP addresses — current (last 24 hours) or all-time | [compromised.md](endpoints/compromised.md) |
| `compromised-zone` | Currently compromised domains in one TLD | [compromised-zone.md](endpoints/compromised-zone.md) |

### Free datasets outside `/api2/`

| URL | Description | Docs |
|---|---|---|
| `https://netapi.com/netapi_top1mln.csv` | NetAPI Top 1 Million — a free daily ranking of the most popular domains | [top-1m.md](endpoints/top-1m.md) · [dedicated repo](https://github.com/coderockers/top-1mln-websites) |

---

## Authentication

Pass your token as the `token` query parameter:

```
https://netapi.com/api2/?method=lookup-domain&domain=example.com&token=YOUR_API_TOKEN
```

- `zones`, `dns`, `registrars`, `compromised`, `compromised-zone` and the Top 1M CSV need no token.
- `lookup-domain`, `lookup-ip`, `download`, `download-dns` and `download-whois` need a valid token **and an active paid plan**. Detailed datasets (`dataset_type=dataset`) require the Plus plan or higher; the Basic plan covers domain lists only. See [Plans](https://netapi.com/plans/).

Keep your token private: it identifies your account, and every request made with it is logged against that account.

---

## Parameters shared by several endpoints

| Parameter | Values | Used by |
|---|---|---|
| `dataset_type` | `list` — one domain per line · `dataset` — domain plus metadata | `download`, `download-dns`, `download-whois` |
| `filter_type` | `active` — all currently registered domains · `new` — domains first seen in the last 24 hours · `deleted` — domains that disappeared in the last 24 hours | `download` |
| `format` | `plain` — uncompressed CSV instead of `.csv.gz` | `download`, `download-dns`, `download-whois` |
| `token` | your API token | see Authentication |

`new` and `deleted` are available only for zones that are fully refreshed daily (`isUpdatedDaily=1` in `zones`). Other zones support `active` only.

---

## Response format

- Reference lists and downloads are CSV with a header row. Fields that may contain commas (provider names, registrar names, email and phone lists) are double-quoted.
- `lookup-domain` and `lookup-ip` return CSV rows **without** a header. When nothing is found they return the text `NOT FOUND` with HTTP 200.
- `compromised` and `compromised-zone` return one entry per line, preceded by a single `#` comment line with the list name and date.
- Downloads are served as `application/octet-stream` with a `Content-Disposition` filename, e.g. `net_active_list.csv.gz`.
- Plain-text downloads can be very large (the `.com` dataset is several gigabytes uncompressed). Use `format=plain` only when you really need it.

### Dataset columns

Zone datasets (`download`, `download-dns`):

```
url,majestic_rank,dns1,dns2,hostname,emails,phones,ip,ip_country
```

Registrar datasets (`download-whois`):

```
registrar,url,registered_at,expiring_at,majestic_rank,emails,phones,ip,ip_country
```

| Column | Meaning |
|---|---|
| `url` | domain name |
| `majestic_rank` | Majestic rank, integer; empty if not ranked |
| `dns1`, `dns2` | authoritative nameservers |
| `hostname` | hostname of the web server |
| `emails`, `phones` | contact details found on the website; comma-separated inside one quoted field |
| `ip` | web server IP address |
| `ip_country` | two-letter country code of the server IP (geolocation) |
| `registrar` | sponsoring registrar (legal entity name) |
| `registered_at`, `expiring_at` | registration and expiration dates, `YYYY-MM-DD` |

---

## Errors

Errors are returned as a short plain-text message with a matching HTTP status code:

| HTTP | Message | Meaning |
|---|---|---|
| 401 | `401 Unauthorized: Missing user token.` | `token` is missing on an endpoint that requires it |
| 403 | `403 Forbidden: Incorrect token.` | the token does not belong to any account |
| 403 | `403 Forbidden: No active/paid plan.` | download or lookup requested without an active subscription |
| 403 | `403 Forbidden: Your current plan does not allow downloading of detailed datasets.` | `dataset_type=dataset` on the Basic plan |
| 404 | `404 Not Found: File not found. …` | the requested file is not on the server (for example, a zone that has no `new` list yet) |
| 405 | `405 Method Not Allowed: Missing API method.` | `method` is missing |
| 405 | `405 Method Not Allowed: Unsupported method (…).` | unknown `method` |
| 405 | `405 Method Not Allowed: Missing …` / `Invalid …` | a required parameter is missing or has an unknown value |

Each endpoint page lists the exact messages it can return.

---

## Examples

| File | What it does |
|---|---|
| [examples/curl-example.sh](examples/curl-example.sh) | download a zone list with cURL |
| [examples/curl-example-compromised.sh](examples/curl-example-compromised.sh) | fetch the current compromised-domain feed |
| [examples/python-example.py](examples/python-example.py) | download a zone dataset and stream it row by row without unpacking to disk |
| [examples/python-example-compromised.py](examples/python-example-compromised.py) | read the compromised feed |
| [examples/python-example-top1m.py](examples/python-example-top1m.py) | load the Top 1M ranking and look up a domain's rank |
| [examples/node-example.js](examples/node-example.js) | download a zone list in Node.js 18+ and stream it through zlib (no dependencies) |
| [examples/node-example-compromised.js](examples/node-example-compromised.js) | read the compromised feed in Node.js |
| [examples/go-example.go](examples/go-example.go) | download a zone list in Go and stream it through `compress/gzip` (standard library only; `go run examples/go-example.go`) |
| [examples/go-example-compromised.go](examples/go-example-compromised.go) | read the compromised feed in Go |
| [examples/php-example.php](examples/php-example.php) | download and decompress a zone list in PHP |
| [examples/php-example-compromised.php](examples/php-example-compromised.php) | read the compromised feed in PHP |

---

## OpenAPI

[openapi.yaml](openapi.yaml) describes the API in OpenAPI 3.0 format. Because all operations share one path and are selected with the `method` parameter, the specification exposes a single `GET /api2/` operation whose `method` enum and parameter descriptions cover every endpoint, plus the Top 1M download.

---

## Data updates and coverage

- Zones with `isUpdatedDaily=1` (the gTLDs and a few ccTLDs) are fully refreshed every day, and their `new` and `deleted` lists are available.
- Zones with `isUpdatedDaily=0` (most ccTLDs, including large ones such as `.de` or `.uk`) get partial updates every day and a full refresh every one to two months. Only the `active` list is available for these zones.
- DNS-provider and registrar datasets, the compromised feeds and the Top 1M ranking are rebuilt daily.
- `zones`, `dns` and `registrars` return live counts, so the numbers change from day to day.

## Crawler

Website metadata (hostname, emails, phone numbers) is collected by our crawler, which identifies itself as `NetAPI/1.1 (+https://netapi.com/bot.html)` and honours `robots.txt`. See [netapi.com/bot.html](https://netapi.com/bot.html) for details and opt-out instructions.

## Licensing

- Paid downloads and lookups are licensed under the [NetAPI Terms of Service](https://netapi.com/tos/).
- The free datasets — the compromised feeds and the Top 1M ranking — are released under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). Please credit NetAPI with a link to `https://netapi.com/` when you use them.
- The contents of this repository (documentation, specification, example code) are released under the [MIT License](LICENSE).

## Support

- Questions about the API or your account: [netapi.com/contact-us/](https://netapi.com/contact-us/)
- Mistakes in this documentation: open an issue or a pull request in this repository.
