#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

MD_MAX=450
GO_MAX=700

GRANDFATHER=(
  "docs/commands/agents.md"
)

is_grandfathered() {
  local f=$1
  for g in "${GRANDFATHER[@]}"; do [[ "$f" == "$g" ]] && return 0; done
  return 1
}

violations=0

while IFS= read -r f; do
  [ -f "$f" ] || continue
  n=$(wc -l <"$f")
  if (( n > MD_MAX )); then
    if is_grandfathered "$f"; then
      echo "note: $f is $n lines (grandfathered; auto-generated — please split)"
    else
      echo "DOC TOO LONG  $f: $n lines (max $MD_MAX) — split or trim"
      violations=$((violations + 1))
    fi
  fi
done < <(
  { find docs -name '*.md' 2>/dev/null; echo ROADMAP.md; echo README.md; } | sort -u
)

while IFS= read -r f; do
  n=$(wc -l <"$f")
  if (( n > GO_MAX )); then
    if is_grandfathered "$f"; then
      echo "note: $f is $n lines (grandfathered; please split)"
    else
      echo "GO FILE TOO LONG  $f: $n lines (max $GO_MAX) — split by topic"
      violations=$((violations + 1))
    fi
  fi
done < <(find . -name '*.go' -not -path './.git/*' -not -name '*_test.go' | sed 's#^\./##' | sort)

for g in "${GRANDFATHER[@]}"; do
  [ -f "$g" ] || continue
  cnt=$(wc -l <"$g")
  if [[ "$g" == *.md ]]; then
    (( cnt <= MD_MAX )) && echo "note: $g is now under $MD_MAX lines — remove from GRANDFATHER"
  elif [[ "$g" == *.go ]]; then
    (( cnt <= GO_MAX )) && echo "note: $g is now under $GO_MAX lines — remove from GRANDFATHER"
  fi
done

if (( violations > 0 )); then
  echo "FAIL: $violations file(s) over the length budget."
  exit 1
fi
echo "OK: all docs <= $MD_MAX lines, all Go source <= $GO_MAX lines."
