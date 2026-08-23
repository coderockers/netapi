# `compromised-zone` — currently compromised domains in one TLD (free)

Returns the subset of the current compromised-domain feed that belongs to a single domain zone — for example every `.com` or `.uk` domain reported in the last 24 hours. IP addresses and historical entries are not included.

- **Authentication:** none
- **Response:** plain text, one domain per line, served as a file download
- **License:** [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) — free for any use with attribution to NetAPI

## Request

```
GET https://netapi.com/api2/?method=compromised-zone&zone_tld={tld}
```

| Parameter | Required | Description |
|---|---|---|
| `zone_tld` | yes | First-level TLD without the leading dot: `com`, `uk`, `xn--p1ai`. Case-insensitive. Must be one of the zones in [`zones`](zones.md); second-level zones such as `co.uk` are not accepted — request `uk` and filter the result. |

## Response

Filename `compromised_{tld}.csv`. The first line is a comment with the zone and date; each following line is one domain.

```
# Current compromised domains in .com zone (2026-08-23).
00zyku.com
024fpv.com
```

## Errors

| HTTP | Message |
|---|---|
| 405 | `405 Method Not Allowed: Missing zone type.` — `zone_tld` is missing |
| 405 | `405 Method Not Allowed: Invalid zone tld.` |

## Examples

cURL:

```bash
curl -s "https://netapi.com/api2/?method=compromised-zone&zone_tld=com" -o compromised_com.csv
```

Python — count compromised domains in a few zones:

```python
import requests

for tld in ("com", "xyz", "top", "uk"):
    resp = requests.get(
        "https://netapi.com/api2/",
        params={"method": "compromised-zone", "zone_tld": tld},
        timeout=60,
    )
    resp.raise_for_status()
    count = sum(1 for line in resp.text.splitlines() if line and not line.startswith("#"))
    print(f".{tld}: {count:,}")
```

## Notes

- The same entries appear in `compromised&dataset_type=url`; this endpoint just saves you the filtering.
- For each zone, the share of compromised domains and its 30-day trend are shown on the zone page (e.g. [netapi.com/com-domains/](https://netapi.com/com-domains/)) and in the [TLD abuse report](https://netapi.com/research/tld-abuse/).
- Web reference: [netapi.com/help/api/#compromised-zone](https://netapi.com/help/api/#compromised-zone).
