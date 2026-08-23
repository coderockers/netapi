# `download-whois` — domains by registrar

Downloads every domain registered through a given registrar across all zones, as a list or as a dataset that includes the registration and expiration dates taken from WHOIS/RDAP.

- **Authentication:** `token` + active paid plan. `dataset_type=dataset` requires the Plus plan or higher.
- **Response:** gzip-compressed CSV file (`.csv.gz`) with a header row; plain CSV with `format=plain`

## Request

```
GET https://netapi.com/api2/?method=download-whois&registrar_id={id}&dataset_type={list|dataset}&token={token}[&format=plain]
```

| Parameter | Required | Description |
|---|---|---|
| `registrar_id` | yes | Numeric registrar ID from [`registrars`](registrars.md), e.g. `146` for GoDaddy.com, LLC, `1068` for NameCheap, Inc. |
| `dataset_type` | yes | `list` — one domain per line · `dataset` — domain plus registration data and metadata |
| `token` | yes | Your API token |
| `format` | no | `plain` — uncompressed CSV (the GoDaddy dataset is several gigabytes uncompressed) |

## Response

Filename `{registrar_id}_list.csv.gz` or `{registrar_id}_dataset.csv.gz`.

`dataset_type=list`:

```
url
themis.abogado
abogada.abogado
```

`dataset_type=dataset`:

```
registrar,url,registered_at,expiring_at,majestic_rank,emails,phones,ip,ip_country
"GoDaddy.com, LLC",themis.abogado,2025-03-12,2027-03-12,,info@themis.abogado,,34.174.62.87,US
"GoDaddy.com, LLC",abogada.abogado,2024-12-21,2026-12-21,,law@law.law,,64.23.224.186,US
```

### Dataset columns

| Column | Description |
|---|---|
| `registrar` | Registrar legal entity name (the `name` column of `registrars`), double-quoted. |
| `url` | Domain name. |
| `registered_at` | Registration date, `YYYY-MM-DD`; empty if unknown. |
| `expiring_at` | Expiration date, `YYYY-MM-DD`; empty if unknown. |
| `majestic_rank` | Position in the Majestic Million; empty if not ranked. |
| `emails`, `phones` | Contact details found on the website, comma-separated inside one double-quoted field. |
| `ip` | IP address of the web server. |
| `ip_country` | Two-letter country code of the server IP. |

Unlike the zone datasets, registrar datasets do not include nameservers or hostname.

## Errors

| HTTP | Message |
|---|---|
| 401 | `401 Unauthorized: Missing user token.` |
| 403 | `403 Forbidden: Incorrect token.` |
| 403 | `403 Forbidden: No active/paid plan.` |
| 403 | `403 Forbidden: Your current plan does not allow downloading of detailed datasets.` |
| 405 | `405 Method Not Allowed: Missing Registrar ID.` / `Invalid Registrar ID.` |
| 405 | `405 Method Not Allowed: Missing dataset type.` / `Invalid dataset type.` |
| 404 | `404 Not Found: File not found. …` |

## Examples

cURL:

```bash
curl "https://netapi.com/api2/?method=download-whois&registrar_id=146&dataset_type=list&token=YOUR_API_TOKEN" \
  -o 146_list.csv.gz
```

Python — domains expiring within the next 30 days:

```python
import csv
import gzip
from datetime import date, timedelta
import requests

params = {
    "method": "download-whois",
    "registrar_id": 146,
    "dataset_type": "dataset",
    "token": "YOUR_API_TOKEN",
}
cutoff = (date.today() + timedelta(days=30)).isoformat()

with requests.get("https://netapi.com/api2/", params=params, stream=True, timeout=600) as resp:
    resp.raise_for_status()
    with gzip.open(resp.raw, mode="rt", encoding="utf-8", newline="") as f:
        for row in csv.DictReader(f):
            if row["expiring_at"] and row["expiring_at"] <= cutoff:
                print(row["url"], row["expiring_at"])
```

## Notes

- Registrar attribution and dates come from WHOIS/RDAP. Domains whose registrar could not be determined are not included in any registrar list.
- Web reference: [netapi.com/help/api/#whois-download](https://netapi.com/help/api/#whois-download); registrar pages: [netapi.com/whois-providers/](https://netapi.com/whois-providers/).
