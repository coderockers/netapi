# NetAPI Top 1 Million — free daily ranking of the most popular domains

A free, daily-rebuilt CSV list of the one million most popular domains on the Internet, in the same shape as the discontinued Alexa Top 1M file. It is intended for research, security checks, allow-listing and large-scale crawls.

- **URL:** `https://netapi.com/netapi_top1mln.csv` (permanent; not under `/api2/`)
- **Authentication:** none
- **Response:** plain CSV, ~1,000,001 lines (header + 1,000,000 rows), about 22 MB
- **License:** [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) — free for any use, including commercial, with attribution to NetAPI
- **Landing page and methodology:** [netapi.com/top-1-mln-websites/](https://netapi.com/top-1-mln-websites/)

## Request

```
GET https://netapi.com/netapi_top1mln.csv
```

No parameters are required. An optional `token` parameter (your API token) ties the download to your account in the request log; it does not change the response.

## Response

```
NetAPI_Rank,Domain
1,google.com
2,microsoft.com
3,apple.com
```

| Column | Description |
|---|---|
| `NetAPI_Rank` | Integer rank from 1 (most popular) to 1,000,000. |
| `Domain` | Registered domain without subdomain — `google.com`, not `www.google.com`. |

Rows are sorted by ascending rank.

## How the ranking is built

The list merges several public popularity signals — link-graph / SEO authority and traffic-oriented signals such as DNS query popularity and real-user Chrome experience data — with reciprocal rank fusion (RRF), so no single noisy source dominates the result. Hostnames are reduced to registered domains using NetAPI's TLD map. The outcome is a synthetic popularity approximation, not a measurement of page views. The full description is on the [landing page](https://netapi.com/top-1-mln-websites/).

## Examples

cURL:

```bash
curl -s https://netapi.com/netapi_top1mln.csv -o netapi_top1mln.csv
```

Python — rank lookup for a handful of domains:

```python
import csv
import io
import requests

resp = requests.get("https://netapi.com/netapi_top1mln.csv", timeout=120)
resp.raise_for_status()

rank = {row["Domain"]: int(row["NetAPI_Rank"]) for row in csv.DictReader(io.StringIO(resp.text))}

for domain in ("google.com", "wikipedia.org", "example.com"):
    print(domain, rank.get(domain, "not in Top 1M"))
```

See also [examples/python-example-top1m.py](../examples/python-example-top1m.py).

## Citing the list

If you use the NetAPI Top 1M in research, products, security tooling or other datasets, please cite NetAPI and link to `https://netapi.com/top-1-mln-websites/`. A suggested citation:

> NetAPI. *NetAPI Top 1 Million Websites* (daily ranking). https://netapi.com/top-1-mln-websites/ — retrieved YYYY-MM-DD.

## Notes

- The file is regenerated once a day; the "Last updated" timestamp is shown on the landing page. Fetch it no more than once a day.
- Because the ranking is rebuilt daily, keep a dated copy if your work needs reproducibility.
- Web reference: [netapi.com/help/api/#top1m-api](https://netapi.com/help/api/#top1m-api).
