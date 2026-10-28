# `compromised` — compromised domains and IP addresses (free)

Returns NetAPI's threat-intelligence feed: domain names and IP addresses currently reported as compromised or malicious by public threat feeds, or the complete history of everything ever listed.

- **Authentication:** none
- **Response:** plain text, one entry per line, served as a file download
- **License:** [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) — free for any use with attribution to NetAPI

## Request

```
GET https://netapi.com/api2/?method=compromised&dataset_type={ip|url|ip-all|url-all}
```

| `dataset_type` | Contents | Size (Aug 2026) |
|---|---|---|
| `url` | Domains seen on threat feeds in the last 24 hours — the **current** blocklist | ~300,000 |
| `ip` | IP addresses seen in the last 24 hours — the **current** blocklist | ~45,000 |
| `url-all` | Every domain ever listed since the feed started | ~1.1 million |
| `ip-all` | Every IP address ever listed | ~1.9 million |

## Response

The first line is a comment with the list name and generation date; every following line is one domain or one IP address. There is no CSV header.

`dataset_type=url` (filename `compromised_url.csv`):

```
# Current compromised URLs (2026-08-23).
000359.xyz
00zyku.com
01.losbuhosweb.com.mx
```

`dataset_type=ip` (filename `compromised_ip.csv`):

```
# Current compromised IPs (2026-08-23).
182.127.178.111
61.52.45.140
```

Historical lists use the filenames `compromised_url_history.csv` and `compromised_ip_history.csv`.

## Errors

| HTTP | Message |
|---|---|
| 405 | `405 Method Not Allowed: Missing list type.` — `dataset_type` is missing |
| 405 | `405 Method Not Allowed: Invalid list type.` |

## Examples

cURL — save today's domain blocklist:

```bash
curl -s "https://netapi.com/api2/?method=compromised&dataset_type=url" -o compromised_url.csv
```

Python — load the current domains into a set and check a few names:

```python
import requests

resp = requests.get(
    "https://netapi.com/api2/", params={"method": "compromised", "dataset_type": "url"}, timeout=60
)
resp.raise_for_status()

blocked = {line.strip() for line in resp.text.splitlines() if line and not line.startswith("#")}
print(len(blocked), "domains currently listed")

for domain in ("example.com", "00zyku.com"):
    print(domain, "LISTED" if domain in blocked else "clean")
```

## Notes

- "Last 24 hours" means the entry was present on at least one of the upstream threat feeds within the last day. Entries drop off the current list as soon as they are no longer reported, which is why the current list is the one to use for blocking.
- The `*-all` lists are historical. A domain or IP appearing there was listed at some point; it is not necessarily compromised now — usually it has been cleaned up or the domain has expired.
- The `url` lists contain host names only (registered domains and, where reported, subdomains) — no schemes or paths.
- Feeds are rebuilt daily. Per-zone breakdowns are available through [`compromised-zone`](compromised-zone.md); statistics and charts on [netapi.com/compromised-urls/](https://netapi.com/compromised-urls/) and [netapi.com/compromised-ips/](https://netapi.com/compromised-ips/); abuse research by TLD and registrar at [netapi.com/research/](https://netapi.com/research/).
- Web reference: [netapi.com/help/api/#compromised-api](https://netapi.com/help/api/#compromised-api).
