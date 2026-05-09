#!/usr/bin/env bash

set -euo pipefail

# May 1, 2026 is Day 1.
START_DATE="${START_DATE:-2026-05-01}"
WEEKEND_SKIP_START="${WEEKEND_SKIP_START:-2026-05-09}"
TODAY="${TODAY:-$(date +%F)}"

start_seconds="$(date -d "$START_DATE" +%s)"
today_seconds="$(date -d "$TODAY" +%s)"
calendar_day_number="$(( (today_seconds - start_seconds) / 86400 + 1 ))"

if (( calendar_day_number < 1 )); then
  echo "No daily post yet"
  exit 0
fi

day_of_week="$(date -d "$TODAY" +%u)"
if (( day_of_week >= 6 )); then
  echo "Weekend"
  exit 0
fi

skip_start_seconds="$(date -d "$WEEKEND_SKIP_START" +%s)"
skipped_weekend_days=0

if (( today_seconds >= skip_start_seconds )); then
  current_seconds="$skip_start_seconds"

  while (( current_seconds <= today_seconds )); do
    current_day_of_week="$(date -d "@$current_seconds" +%u)"

    if (( current_day_of_week >= 6 )); then
      skipped_weekend_days="$((skipped_weekend_days + 1))"
    fi

    current_seconds="$((current_seconds + 86400))"
  done
fi

day_number="$((calendar_day_number - skipped_weekend_days))"

echo "day ${day_number}"
