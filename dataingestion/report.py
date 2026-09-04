#!/usr/bin/env python3
"""Generate a data-quality report for the ingested mandi raw dataset.

Outputs data/raw/report.json (machine-readable) and prints a human summary.
"""

import json
import os
from collections import Counter
from datetime import datetime

HERE = os.path.dirname(os.path.abspath(__file__))
RAW = os.path.join(HERE, "data", "raw", "mandi_prices.json")
REPORT = os.path.join(HERE, "data", "raw", "report.json")

FIELDS = [
    "state", "district", "market", "commodity", "variety", "grade",
    "arrival_date", "min_price", "max_price", "modal_price",
]


def parse_date(d):
    try:
        return datetime.strptime(d, "%d/%m/%Y").date()
    except (ValueError, TypeError):
        return None


def full_sig(r):
    return tuple(r.get(k) for k in FIELDS)


def main():
    with open(RAW, encoding="utf-8") as f:
        recs = json.load(f)

    total = len(recs)

    dates = [parse_date(r.get("arrival_date")) for r in recs]
    dates = [d for d in dates if d]
    date_min = str(min(dates)) if dates else None
    date_max = str(max(dates)) if dates else None

    missing = {}
    for fld in FIELDS:
        missing[fld] = sum(1 for r in recs if r.get(fld) in (None, ""))

    cnt = Counter(full_sig(r) for r in recs)
    dup_sigs = [k for k, v in cnt.items() if v > 1]
    duplicate_rows = sum(v for v in cnt.values() if v > 1) - len(dup_sigs)

    report = {
        "resource_id": "9ef84268-d588-465a-a308-a864a43d0070",
        "source": "https://api.data.gov.in/resource/9ef84268-d588-465a-a308-a864a43d0070",
        "total_records": total,
        "date_range": {"min": date_min, "max": date_max},
        "num_states": len({r["state"] for r in recs}),
        "num_districts": len({(r["state"], r["district"]) for r in recs}),
        "num_markets": len({(r["state"], r["market"]) for r in recs}),
        "num_commodities": len({r["commodity"] for r in recs}),
        "num_varieties": len({r["variety"] for r in recs}),
        "num_grades": len({r["grade"] for r in recs}),
        "missing_values": missing,
        "duplicate_signatures": len(dup_sigs),
        "duplicate_rows": duplicate_rows,
    }

    with open(REPORT, "w", encoding="utf-8") as f:
        json.dump(report, f, indent=2)

    print(json.dumps(report, indent=2))
    print(f"\nReport written to {REPORT}")


if __name__ == "__main__":
    main()
