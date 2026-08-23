# `registrars` — list of known registrars

Returns every registrar NetAPI has identified from WHOIS/RDAP data, with the numeric ID you pass to [`download-whois`](download-whois.md) and [`sample-whois`](sample-whois.md), and the current number of domains registered through that registrar.

- **Authentication:** none
- **Parameters:** none
- **Response:** plain CSV with a header row, about 4,200 lines

## Request

```
GET https://netapi.com/api2/?method=registrars
```

## Response

```
id,name,brand_name,total
```

| Column | Type | Description |
|---|---|---|
| `id` | integer | Registrar ID — the `registrar_id` value for `download-whois`. IDs are stable. |
| `name` | string | Legal entity name as it appears in WHOIS, double-quoted, e.g. `"GoDaddy.com, LLC"`. |
| `brand_name` | string | Consumer brand the entity operates under, double-quoted; empty when the brand is the same as the name or unknown. |
| `total` | integer | Number of domains currently registered through this registrar. |

Example:

```
id,name,brand_name,total
146,"GoDaddy.com, LLC","GoDaddy",65992964
1068,"NameCheap, Inc.","Namecheap",24351298
69,"Tucows Domains Inc.","Tucows, Hover",9960806
2,"Network Solutions, LLC","Network Solutions",4720656
1,"Reserved","",321448
```

## Examples

cURL:

```bash
curl "https://netapi.com/api2/?method=registrars"
```

Python — find the ID of a registrar by brand:

```python
import csv
import io
import requests

resp = requests.get("https://netapi.com/api2/", params={"method": "registrars"}, timeout=30)
resp.raise_for_status()

for row in csv.DictReader(io.StringIO(resp.text)):
    if "namecheap" in (row["brand_name"] + row["name"]).lower():
        print(row["id"], row["name"], row["total"])
```

## Notes

- The same legal entity may appear more than once under slightly different WHOIS spellings; each spelling has its own ID and its own domain list.
- `total` is recalculated daily. The counts in the example above are illustrative.
- Every registrar has a landing page with daily statistics and abuse figures, e.g. [netapi.com/whois-godaddy-com-llc/](https://netapi.com/whois-godaddy-com-llc/); see the [registrars overview](https://netapi.com/whois-providers/) and the [API reference](https://netapi.com/help/api/#whois-registrars).
