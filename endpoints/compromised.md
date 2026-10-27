# `compromised` — compromised domains (free)

Returns NetAPI's threat-intelligence feed: domain names currently reported as compromised or malicious by public threat feeds and notifications, or the complete history of every domain ever listed.

- **Authentication:** none
- **Response:** plain text, one entry per line, served as a file download
- **License:** [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) — free for any use with attribution to NetAPI

## Request

```
GET https://netapi.com/api2/?method=compromised&dataset_type={url|url-all}[&format=hosts]
```

| `dataset_type` | Contents |
|---|---|
| `url` | Domains seen on threat feeds in the last 24 hours — the **current** blocklist |
| `url-all` | Every domain ever listed since the feed started (more than a million) |

| `format` (optional) | Output |
|---|---|
| — | one domain per line (filenames `compromised_url.csv`, `compromised_url_history.csv`) |
| `hosts` | one `0.0.0.0 domain` line per domain, for hosts-based blockers such as dnsmasq `addn-hosts`, AdAway or Pi-hole (filenames `compromised_domains_hosts.txt`, `compromised_domains_history_hosts.txt`) |

## Response

The first line is a comment starting with `#` that names the list and its date; every following line is one domain. There is no CSV header.

`dataset_type=url`:

```
# NetAPI.com - currently compromised URLs (2026-10-27).
000359.xyz
00zyku.com
```

`dataset_type=url&format=hosts`:

```
# NetAPI.com - currently compromised domains (2026-10-27), hosts format.
0.0.0.0 000359.xyz
0.0.0.0 00zyku.com
```

## Errors

| HTTP | Message |
|---|---|
| 405 | `405 Method Not Allowed: Missing list type.` — `dataset_type` is missing |
| 405 | `405 Method Not Allowed: Invalid list type.` |

## Examples

cURL — save today's domain blocklist, and the same list as a hosts file:

```bash
curl -s "https://netapi.com/api2/?method=compromised&dataset_type=url" -o compromised_url.csv
curl -s "https://netapi.com/api2/?method=compromised&dataset_type=url&format=hosts" -o compromised_domains_hosts.txt
```

Python — load the current domains into a set and check a few names:

```python
import requests

resp = requests.get(
    "https://netapi.com/api2/", params={"method": "compromised", "dataset_type": "url"}, timeout=60
)
resp.raise_for_status()

blocked = {line.strip() for line in resp.text.splitlines() if line.strip() and not line.startswith("#")}
print(len(blocked), "domains currently listed")

for domain in ("example.com", "00zyku.com"):
    print(domain, "LISTED" if domain in blocked else "clean")
```

## Notes

- "Last 24 hours" means the entry was present on at least one of the upstream threat feeds within the last day. Entries drop off the current list as soon as they are no longer reported, which is why the current list is the one to use for blocking.
- The `url-all` list is historical. A domain appearing there was listed at some point; it is not necessarily compromised now — usually it has been cleaned up or the domain has expired.
- The lists contain host names only (registered domains and, where reported, subdomains) — no schemes or paths.
- A single domain can be checked without downloading the list: [`compromised-check`](json-api.md) in the JSON API.
- The feed is updated daily. Per-zone breakdowns are available through [`compromised-zone`](compromised-zone.md); statistics and charts on [netapi.com/compromised-urls/](https://netapi.com/compromised-urls/); abuse research by TLD and registrar at [netapi.com/research/](https://netapi.com/research/).
- A domain listed by mistake can be reported through the [removal request form](https://netapi.com/compromised-urls/removal/).
- Web reference: [netapi.com/help/api/#compromised-api](https://netapi.com/help/api/#compromised-api).
