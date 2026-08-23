# `download` — domain lists and datasets by zone

Downloads the full list of registered domains in a zone (or in all zones at once), optionally with metadata, or only the domains that were added or dropped in the last 24 hours.

- **Authentication:** `token` + active paid plan. `dataset_type=dataset` requires the Plus plan or higher.
- **Response:** gzip-compressed CSV file (`.csv.gz`) with a header row; plain CSV with `format=plain`

## Request

```
GET https://netapi.com/api2/?method=download&zone_tld={tld}&dataset_type={list|dataset}&filter_type={active|new|deleted}&token={token}[&format=plain]
```

| Parameter | Required | Description |
|---|---|---|
| `zone_tld` | yes | First-level TLD without the leading dot (`com`, `net`, `uk`), punycode for IDN zones (`xn--p1ai`), or `all-zones` for one combined file of every zone. Valid values come from [`zones`](zones.md). |
| `dataset_type` | yes | `list` — one domain per line · `dataset` — domain plus metadata (see columns below) |
| `filter_type` | yes | `active` — all currently registered domains · `new` — domains first seen in the last 24 hours · `deleted` — domains that dropped out of the zone in the last 24 hours |
| `token` | yes | Your API token from the [dashboard](https://netapi.com/dashboard/) |
| `format` | no | `plain` — uncompressed CSV. Plain files for large zones run to several gigabytes; use only when you cannot handle gzip. |

`new` and `deleted` exist only for zones with `isUpdatedDaily=1` in `zones`. For the other zones request `active`.

## Response

The file is served as `application/octet-stream` with a `Content-Disposition` filename built from the request, for example `net_active_list.csv.gz`, `all-zones_new_dataset.csv.gz`, `com_deleted_list.csv`.

`dataset_type=list`:

```
url
example.net
sample-site.net
```

`dataset_type=dataset`:

```
url,majestic_rank,dns1,dns2,hostname,emails,phones,ip,ip_country
example.net,8421,ns1.example.net,ns2.example.net,web1.example.net,"info@example.net,sales@example.net",+1-202-555-0143,203.0.113.10,US
sample-site.net,,ns1.hosting.net,ns2.hosting.net,srv-12.hosting.net,,,198.51.100.7,DE
```

### Dataset columns

| Column | Description |
|---|---|
| `url` | Domain name. |
| `majestic_rank` | Majestic rank, integer; empty if the domain is not ranked. |
| `dns1`, `dns2` | Authoritative nameservers. `dns2` is empty when only one nameserver is known. |
| `hostname` | Hostname of the web server. |
| `emails` | Email addresses found on the website, comma-separated inside one double-quoted field; empty if none. |
| `phones` | Phone numbers found on the website, same format as `emails`. |
| `ip` | IP address of the web server. |
| `ip_country` | Two-letter country code of the server IP (geolocation). |

Metadata is collected when the domain's website is crawled; fields are empty for domains without a reachable website.

## Errors

| HTTP | Message |
|---|---|
| 401 | `401 Unauthorized: Missing user token.` |
| 403 | `403 Forbidden: Incorrect token.` |
| 403 | `403 Forbidden: No active/paid plan.` |
| 403 | `403 Forbidden: Your current plan does not allow downloading of detailed datasets.` |
| 405 | `405 Method Not Allowed: Missing zone type.` — `zone_tld` is missing |
| 405 | `405 Method Not Allowed: Invalid zone tld.` |
| 405 | `405 Method Not Allowed: Missing dataset type.` / `Invalid dataset type.` |
| 405 | `405 Method Not Allowed: Missing filter type.` / `Invalid filter type.` |
| 404 | `404 Not Found: File not found. …` — the file does not exist, typically `new` or `deleted` requested for a zone with `isUpdatedDaily=0` |

## Examples

cURL — save the compressed file:

```bash
curl "https://netapi.com/api2/?method=download&zone_tld=net&dataset_type=list&filter_type=active&token=YOUR_API_TOKEN" \
  -o net_active_list.csv.gz
```

cURL — unpack on the fly and count the domains:

```bash
curl -s "https://netapi.com/api2/?method=download&zone_tld=net&dataset_type=list&filter_type=new&token=YOUR_API_TOKEN" \
  | gunzip | tail -n +2 | wc -l
```

Python — stream the dataset without writing the archive to disk:

```python
import csv
import gzip
import requests

params = {
    "method": "download",
    "zone_tld": "net",
    "dataset_type": "dataset",
    "filter_type": "active",
    "token": "YOUR_API_TOKEN",
}

with requests.get("https://netapi.com/api2/", params=params, stream=True, timeout=600) as resp:
    resp.raise_for_status()
    with gzip.open(resp.raw, mode="rt", encoding="utf-8", newline="") as f:
        for row in csv.DictReader(f):
            print(row["url"], row["ip"], row["ip_country"])
            break  # remove to process the whole file
```

More examples: [examples/](../examples/).

## Notes

- A complete download of a large zone (`.com`, `all-zones`) takes a while even when compressed; use a client that streams to disk rather than buffering in memory.
- Files are regenerated once a day; requesting the same file again on the same day returns the same content.
- Every download is logged against your token.
- Web reference: [netapi.com/help/api/#list-api](https://netapi.com/help/api/#list-api).
