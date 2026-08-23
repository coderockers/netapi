# `dns` — list of supported DNS providers

Returns every DNS provider for which NetAPI builds a domain list, with the alias you pass to [`download-dns`](download-dns.md) and the current number of domains known to use that provider.

- **Authentication:** none
- **Parameters:** none
- **Response:** plain CSV with a header row, about 160 lines

## Request

```
GET https://netapi.com/api2/?method=dns
```

## Response

```
alias,title,total
```

| Column | Type | Description |
|---|---|---|
| `alias` | string | Provider alias — the `dns_alias` value for `download-dns`. |
| `title` | string | Provider name, double-quoted. |
| `total` | integer | Number of domains currently using the provider's nameservers. |

Example:

```
alias,title,total
cloudflare,"Cloudflare",38002618
godaddy,"GoDaddy",47941850
a2hosting,"A2 Hosting",178948
4-cn,"4.CN",116617
```

## Examples

cURL:

```bash
curl "https://netapi.com/api2/?method=dns"
```

Python — the ten largest providers:

```python
import csv
import io
import requests

resp = requests.get("https://netapi.com/api2/", params={"method": "dns"}, timeout=30)
resp.raise_for_status()

providers = sorted(csv.DictReader(io.StringIO(resp.text)), key=lambda r: -int(r["total"]))
for p in providers[:10]:
    print(f'{p["alias"]:20} {p["title"]:30} {int(p["total"]):>12,}')
```

## Notes

- A domain is attributed to a provider by matching its authoritative nameservers against the provider's nameserver patterns. Each provider has a landing page with its nameserver pattern and daily statistics, e.g. [netapi.com/dns-cloudflare/](https://netapi.com/dns-cloudflare/).
- `total` is recalculated daily. The counts in the example above are illustrative; fetch the live list for current numbers.
- See also the [DNS providers overview](https://netapi.com/dns-providers/) and the [API reference](https://netapi.com/help/api/#dns-api-list).
