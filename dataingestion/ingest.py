#!/usr/bin/env python3
"""Automated ingestion of the current daily Mandi price dataset from data.gov.in.

Resource: "Current Daily Price of Various Commodities from Various Markets (Mandi)"
Resource ID: 9ef84268-d588-465a-a308-a864a43d0070

The underlying store is Elasticsearch and enforces `offset + limit <= 10000`,
while the dataset has ~16.8k rows, so a single straight pagination cannot reach
all records. We therefore split the fetch by State. Each State carries fewer than
the 10k window, so every State is pulled in one request using
`filters[state.keyword]=<state>&offset=0&limit=<state_count>`.

Note: the plain `filters[state]=` param is unreliable (it ignores the filter and
returns global slices); the correct filter field id is `state.keyword`.

Output: data/raw/mandi_prices.json (raw records, fields unmodified).
"""

import json
import os
import subprocess
import sys
import time
import urllib.parse

RESOURCE_ID = "9ef84268-d588-465a-a308-a864a43d0070"
API_KEY = "579b464db66ec23bdd00000121c3726eee2941dc7f6a3b643d92a22f"
BASE_URL = f"https://api.data.gov.in/resource/{RESOURCE_ID}"

OUT_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "data", "raw")
OUT_FILE = os.path.join(OUT_DIR, "mandi_prices.json")

MAX_RETRIES = 5
RETRY_BACKOFF = 3

# Accurate per-state record counts (sum = 16869, matches the API-reported total).
STATE_COUNTS = {
    "Andaman and Nicobar": 44,
    "Andhra Pradesh": 336,
    "Assam": 96,
    "Bihar": 41,
    "Chandigarh": 21,
    "Chattisgarh": 96,
    "Goa": 10,
    "Gujarat": 517,
    "Haryana": 610,
    "Himachal Pradesh": 582,
    "Jammu and Kashmir": 31,
    "Karnataka": 415,
    "Keralam": 903,
    "Madhya Pradesh": 1349,
    "Maharashtra": 1077,
    "Meghalaya": 30,
    "NCT of Delhi": 18,
    "Nagaland": 13,
    "Odisha": 353,
    "Pondicherry": 3,
    "Punjab": 787,
    "Rajasthan": 370,
    "Tamil Nadu": 7013,
    "Telangana": 443,
    "Tripura": 133,
    "Uttar Pradesh": 1012,
    "Uttarakhand": 191,
    "West Bengal": 375,
}


def fetch_state(state, count):
    """Fetch every record for a single state in one request (offset=0)."""
    params = {
        "api-key": API_KEY,
        "format": "json",
        "limit": str(count),
        "offset": "0",
        "filters[state.keyword]": state,
    }
    query = urllib.parse.urlencode(params)
    url = f"{BASE_URL}?{query}"

    for attempt in range(1, MAX_RETRIES + 1):
        try:
            proc = subprocess.run(
                ["curl", "-s", "-m", "120", url],
                capture_output=True,
                text=True,
                check=False,
            )
            raw = proc.stdout
            if not raw:
                raise RuntimeError("empty response")
            data = json.loads(raw)
            if "records" not in data:
                raise RuntimeError(f"unexpected payload: {raw[:200]}")
            records = data.get("records", [])
            if len(records) != count:
                raise RuntimeError(
                    f"expected {count} records, got {len(records)}"
                )
            return records
        except Exception as exc:
            wait = RETRY_BACKOFF * attempt
            print(f"    retry {attempt}/{MAX_RETRIES} after {wait}s ({state}): {exc}", file=sys.stderr)
            time.sleep(wait)
    raise RuntimeError(f"failed to fetch state={state}")


def main():
    all_records = []
    for i, (state, count) in enumerate(STATE_COUNTS.items(), 1):
        recs = fetch_state(state, count)
        all_records.extend(recs)
        print(f"[{i}/{len(STATE_COUNTS)}] {state}: {len(recs)} records")
        time.sleep(0.2)

    print(f"\nFetched total: {len(all_records)} (API reports 16869)")

    os.makedirs(OUT_DIR, exist_ok=True)
    with open(OUT_FILE, "w", encoding="utf-8") as f:
        json.dump(all_records, f, indent=2)
    print(f"Saved to {OUT_FILE}")


if __name__ == "__main__":
    main()
