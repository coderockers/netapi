# `download-dns` — domains by DNS provider

Downloads every domain whose authoritative nameservers belong to a given DNS provider — for example all domains on Cloudflare, GoDaddy or Google nameservers — across all zones, as a list or as a detailed dataset.

- **Authentication:** `token` + active paid plan. `dataset_type=dataset` requires the Plus plan or higher.
- **Response:** gzip-compressed CSV file (`.csv.gz`) with a header row; plain CSV with `format=plain`

## Request

```
GET https://netapi.com/api2/?method=download-dns&dns_alias={alias}&dataset_type={list|dataset}&token={token}[&format=plain]
```

| Parameter | Required | Description |
|---|---|---|
| `dns_alias` | yes | Provider alias from [`dns`](dns.md), e.g. `cloudflare`, `godaddy`, `google`. |
| `dataset_type` | yes | `list` — one domain per line · `dataset` — domain plus metadata |
| `token` | yes | Your API token |
| `format` | no | `plain` — uncompressed CSV (the Cloudflare dataset alone is several gigabytes uncompressed) |

## Response

Filename `{alias}_list.csv.gz` or `{alias}_dataset.csv.gz`.

`dataset_type=list`:

```
url
example.com
example.org
```

`dataset_type=dataset` — same columns as the zone datasets:

```
url,majestic_rank,dns1,dns2,hostname,emails,phones,ip,ip_country
example.com,1203,ada.ns.cloudflare.com,rob.ns.cloudflare.com,,"hello@example.com",,104.21.5.77,US
```

Column descriptions: [download.md](download.md#dataset-columns). `dns1`/`dns2` point to the requested provider's nameservers.

## Errors

| HTTP | Message |
|---|---|
| 401 | `401 Unauthorized: Missing user token.` |
| 403 | `403 Forbidden: Incorrect token.` |
| 403 | `403 Forbidden: No active/paid plan.` |
| 403 | `403 Forbidden: Your current plan does not allow downloading of detailed datasets.` |
| 405 | `405 Method Not Allowed: Missing DNS alias.` / `Invalid DNS alias.` |
| 405 | `405 Method Not Allowed: Missing dataset type.` / `Invalid dataset type.` |
| 404 | `404 Not Found: File not found. …` |

## Examples

cURL:

```bash
curl "https://netapi.com/api2/?method=download-dns&dns_alias=cloudflare&dataset_type=list&token=YOUR_API_TOKEN" \
  -o cloudflare_list.csv.gz
```

Python:

```python
import csv
import gzip
import requests

params = {
    "method": "download-dns",
    "dns_alias": "cloudflare",
    "dataset_type": "dataset",
    "token": "YOUR_API_TOKEN",
}

with requests.get("https://netapi.com/api2/", params=params, stream=True, timeout=600) as resp:
    resp.raise_for_status()
    with gzip.open(resp.raw, mode="rt", encoding="utf-8", newline="") as f:
        reader = csv.DictReader(f)
        for i, row in enumerate(reader):
            print(row["url"], row["dns1"])
            if i == 9:
                break
```

## Notes

- A domain is attributed to a provider by its nameserver hostnames (e.g. `*.ns.cloudflare.com`), not by the IP it resolves to. Domains whose nameservers match no known provider are not included in any DNS list.
- Lists are rebuilt daily from all zones.
- Web reference: [netapi.com/help/api/#dns-api](https://netapi.com/help/api/#dns-api); provider pages: [netapi.com/dns-providers/](https://netapi.com/dns-providers/).
