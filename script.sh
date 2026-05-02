#!/usr/bin/env bash

set -euo pipefail

README_FILE="README.md"
TMP_FILE="$(mktemp)"

cleanup() {
  rm -f "$TMP_FILE"
}
trap cleanup EXIT

extract_range_start() {
  local dir_name="$1"
  dir_name="${dir_name#day}"
  echo "${dir_name%%-*}"
}

extract_range_label() {
  local dir_name="$1"
  dir_name="${dir_name#day}"
  echo "Days ${dir_name}"
}

extract_day_number() {
  local file_name="$1"
  file_name="$(basename "$file_name")"
  echo "${file_name%.md}"
}

clean_title() {
  local raw_title="$1"

  raw_title="${raw_title#\# }"
  raw_title="${raw_title#\#\# }"
  raw_title="${raw_title#\*\*}"
  raw_title="${raw_title%\*\*}"
  raw_title="${raw_title#Daily Coding Challenge - }"
  raw_title="${raw_title#Daily Coding Question - }"

  if [[ "$raw_title" =~ ^Day[[:space:]]+[0-9]+:[[:space:]]+(.+)$ ]]; then
    raw_title="${BASH_REMATCH[1]}"
  fi

  echo "$raw_title"
}

escape_table_text() {
  local text="$1"

  text="${text//|/\\|}"
  echo "$text"
}

get_problem_title() {
  local file_path="$1"
  local day_number="$2"
  local first_line

  first_line="$(sed -n '/[^[:space:]]/ { p; q; }' "$file_path")"

  if [[ -z "$first_line" ]]; then
    echo "Day $day_number"
    return
  fi

  clean_title "$first_line"
}

write_header() {
  {
    echo "# Daily Coding Questions"
    echo
    echo "A growing collection of daily coding challenges for practicing problem solving, algorithms, and clean implementation. Use this index to jump straight to any available day's problem statement."
    echo
    echo "## Problems Index"
    echo
  } >> "$TMP_FILE"
}

write_directory_section() {
  local dir_path="$1"
  local dir_name
  local label
  local open_attribute=""
  local files

  dir_name="$(basename "$dir_path")"
  label="$(extract_range_label "$dir_name")"

  if [[ "$dir_name" == "day1-100" ]]; then
    open_attribute=" open"
  fi

  {
    echo "<details${open_attribute}>"
    echo "<summary><strong>${label}</strong></summary>"
    echo
  } >> "$TMP_FILE"

  mapfile -t files < <(find "$dir_path" -maxdepth 1 -type f -name '*.md' | sort -V)

  if [[ "${#files[@]}" -eq 0 ]]; then
    echo "No problems added yet." >> "$TMP_FILE"
  else
    {
      echo "| Day | Problem |"
      echo "| --- | --- |"
    } >> "$TMP_FILE"

    for file_path in "${files[@]}"; do
      local day_number
      local title

      day_number="$(extract_day_number "$file_path")"
      title="$(get_problem_title "$file_path" "$day_number")"
      title="$(escape_table_text "$title")"

      echo "| ${day_number} | [${title}](${file_path}) |" >> "$TMP_FILE"
    done
  fi

  {
    echo
    echo "</details>"
    echo
  } >> "$TMP_FILE"
}

main() {
  local day_dirs

  write_header

  mapfile -t day_dirs < <(
    find . -maxdepth 1 -type d -name 'day*-*' -printf '%f\n' \
      | while read -r dir_name; do
          printf '%s\t%s\n' "$(extract_range_start "$dir_name")" "$dir_name"
        done \
      | sort -n -k1,1 \
      | cut -f2-
  )

  if [[ "${#day_dirs[@]}" -eq 0 ]]; then
    echo "No problem directories found yet." >> "$TMP_FILE"
  else
    for dir_name in "${day_dirs[@]}"; do
      write_directory_section "$dir_name"
    done
  fi

  mv "$TMP_FILE" "$README_FILE"
  trap - EXIT
}

main "$@"
