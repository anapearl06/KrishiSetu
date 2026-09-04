# Government Mandi Data Acquisition — DATA-01

Automated ingestion of the **Current Daily Price of Various Commodities from Various Markets (Mandi)** dataset from the official [data.gov.in](https://data.gov.in/) API.

## Deliverable

```
dataingestion/
├── ingest.py          # automated ingestion script (curl + pagination)
├── report.py          # generates the data-quality report
└── data/
    └── raw/
        ├── mandi_prices.json   # 16,869 raw records (unmodified)
        └── report.json         # data-quality summary
```

## Source

- **Resource ID:** `9ef84268-d588-465a-a308-a864a43d0070`
- **API:** `https://api.data.gov.in/resource/9ef84268-d588-465a-a308-a864a43d0070`
- **Fields preserved (raw):** `State`, `District`, `Market`, `Commodity`, `Variety`, `Grade`, `Arrival_Date`, `Min_Price`, `Max_Price`, `Modal_Price`

## How it works

The underlying store is Elasticsearch and enforces `offset + limit <= 10000`. The dataset has 16,869 rows, so a straight page-through cannot reach all records. The ingestion splits the fetch **by State** (each State is far below the 10k window) and pulls each State in a single request:

```
API
 ↓
request 1  → filters[state.keyword]=Andaman and Nicobar  (44)
request 2  → filters[state.keyword]=Andhra Pradesh       (336)
request 3  → filters[state.keyword]=Assam                (96)
 ...
request 28 → filters[state.keyword]=West Bengal          (375)
 ↓
all records (16,869)
 ↓
data/raw/mandi_prices.json
```

> **Note:** the plain `filters[state]=` filter is unreliable — it ignores the filter and returns global slices. The correct filter field id is `filters[state.keyword]`. The script uses it with `offset=0&limit=<state_count>` so every State's records arrive in one clean request with no pagination gaps or overlaps.

### Re-run

```bash
python3 ingest.py   # writes data/raw/mandi_prices.json
python3 report.py   # writes data/raw/report.json
```

## Data quality summary (as fetched)

| Metric | Value |
|---|---|
| Total records | **16,869** |
| Date range | 2026-09-02 → 2026-09-02 (single day) |
| States | 28 |
| Districts | 476 |
| Markets | 1,805 |
| Commodities | 238 |
| Varieties | 496 |
| Grades | 14 |
| Missing values | 0 (across all 10 fields) |
| Duplicate rows | 0 |
