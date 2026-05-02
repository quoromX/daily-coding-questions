#!/usr/bin/env bash

set -euo pipefail

# May 1, 2026 is Day 1, so May 2, 2026 is Day 2.
START_DATE="${START_DATE:-2026-05-01}"
TODAY="${TODAY:-$(date +%F)}"

start_seconds="$(date -d "$START_DATE" +%s)"
today_seconds="$(date -d "$TODAY" +%s)"
day_number="$(( (today_seconds - start_seconds) / 86400 + 1 ))"

if (( day_number < 1 )); then
  echo "No daily post yet"
  exit 0
fi

echo "day ${day_number}"
