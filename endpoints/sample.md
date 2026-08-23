# `sample` — free sample of a zone list or dataset

Returns the first rows of a zone's domain list or detailed dataset so you can check the format before subscribing. The sample has the same columns as the file served by [`download`](download.md) with `filter_type=active`.

- **Authentication:** none
- **Response:** plain CSV (not compressed) with a header row and 10 data rows, served as a file download

## Request

```
GET https://netapi.com/api2/?method=sample&zone_tld={tld}&dataset_type={list|dataset}
```

| Parameter | Required | Description |
|---|---|---|
| `zone_tld` | yes | TLD without the leading dot, e.g. `com`, `de`, `co.uk`; punycode for IDN zones. See [`zones`](zones.md). |
| `dataset_type` | yes | `list` — domains only · `dataset` — domains with metadata |

## Response

`dataset_type=list` (filename `{zone_id}_list.csv`):

```
url
abitofhappinessfarm.com
nuovabrianzatraslochi.com
```

`dataset_type=dataset` (filename `{zone_id}_dataset.csv`):

```
url,majestic_rank,dns1,dns2,hostname,emails,phones,ip,ip_country
abitofhappinessfarm.com,,dns1.register.com,dns2.register.com,cms5.weebly.com,"abitofhappinessfarm@gmail.com,hi@mystore.com",,199.34.228.164,US
nuovabrianzatraslochi.com,,ns1.register.it,ns2.register.it,ec2-3-67-141-185.eu-central-1.compute.amazonaws.com,ntraslochibrianza@alice.it,+39-039-2003656,3.67.141.185,DE
```

Column descriptions are in [download.md](download.md#dataset-columns).

## Errors

| HTTP | Message |
|---|---|
| 405 | `405 Method Not Allowed: Missing Zone TLD.` |
| 405 | `405 Method Not Allowed: Invalid zone tld.` |
| 405 | `405 Method Not Allowed: Missing dataset type.` |
| 405 | `405 Method Not Allowed: Invalid dataset type.` |
| 404 | `404 Not Found: File not found. …` — no sample has been generated for this zone yet |

## Examples

```bash
curl "https://netapi.com/api2/?method=sample&zone_tld=com&dataset_type=dataset"
```

```python
import requests

resp = requests.get(
    "https://netapi.com/api2/",
    params={"method": "sample", "zone_tld": "com", "dataset_type": "dataset"},
    timeout=30,
)
resp.raise_for_status()
print(resp.text)
```

## Notes

- Sample rows are taken from the live dataset and change when the zone is refreshed.
- Samples are linked from every zone page, e.g. [netapi.com/com-domains/](https://netapi.com/com-domains/).
