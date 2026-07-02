#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
docs="$root/docs"
sidebar="$docs/_sidebar.md"
base="https://s1.danny.vn"
today="$(date +%Y-%m-%d)"

[ -f "$sidebar" ] || { echo "error: $sidebar not found" >&2; exit 1; }

# Extract first paragraph from a markdown file (for llms.txt descriptions).
first_para() {
  awk '
    /^#/ { found=1; next }
    found && /^[^#>|[]/ && !/^$/ && !/^---/ && !/^```/ {
      gsub(/\*\*/, ""); gsub(/`/, "")
      printf "%s ", $0; count++
      if (count >= 2) exit
    }
    found && /^$/ && count > 0 { exit }
  ' "$1" | head -c 140
}

# --- sitemap.xml -----------------------------------------------------------

sitemap="$docs/sitemap.xml"
{
  echo '<?xml version="1.0" encoding="UTF-8"?>'
  echo '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">'
  echo "  <url>"
  echo "    <loc>${base}/#/</loc>"
  echo "    <lastmod>${today}</lastmod>"
  echo "    <changefreq>weekly</changefreq>"
  echo "  </url>"

  grep -oP '\(([^)]+\.md)\)|\(([^)]+/)\)' "$sidebar" | tr -d '()' | while read -r href; do
    case "$href" in http*) continue ;; esac
    path="${href%.md}"
    path="${path%/}"
    [ -z "$path" ] && continue
    freq="monthly"
    case "$path" in guides/*) freq="weekly" ;; esac
    echo "  <url>"
    echo "    <loc>${base}/#/${path}</loc>"
    echo "    <lastmod>${today}</lastmod>"
    echo "    <changefreq>${freq}</changefreq>"
    echo "  </url>"
  done

  echo '</urlset>'
} > "$sitemap"

# --- llms.txt ---------------------------------------------------------------

llms="$docs/llms.txt"
{
  echo "# s1ctl"
  echo ""
  echo "> CLI and Go SDK for operating SentinelOne Singularity Platform as code."
  echo "> Pull live state, review in git diff, push back. 370+ commands across REST, SDL, and GraphQL surfaces."
  echo ""

  while IFS= read -r line; do
    # Section header: - **Title**
    title="$(echo "$line" | sed -n 's/^- \*\*\(.*\)\*\*$/\1/p')"
    if [ -n "$title" ]; then
      echo ""
      echo "## ${title}"
      echo ""
      continue
    fi

    # Link: - [Title](href) or  - [Title](href)
    if echo "$line" | grep -qP '^\s*- \['; then
      link_title="$(echo "$line" | sed -n 's/.*\[\([^]]*\)\].*/\1/p')"
      href="$(echo "$line" | sed -n 's/.*(\([^)]*\)).*/\1/p')"
      [ -z "$link_title" ] || [ -z "$href" ] && continue

      case "$href" in
        http*) echo "- [${link_title}](${href})"; continue ;;
      esac

      path="${href%.md}"
      path="${path%/}"
      url="${base}/#/${path}"
      [ -z "$path" ] && url="${base}/#/"

      file="$docs/$href"
      [ "$href" = "/" ] && file="$docs/README.md"
      [ -d "$file" ] && file="${file%/}/README.md"
      desc=""
      if [ -f "$file" ]; then
        desc="$(first_para "$file")"
        [ ${#desc} -ge 137 ] && desc="${desc%% *}..."
      fi

      if [ -n "$desc" ]; then
        echo "- [${link_title}](${url}): ${desc}"
      else
        echo "- [${link_title}](${url})"
      fi
    fi
  done < "$sidebar"
  echo ""
} > "$llms"

# --- llms-full.txt ----------------------------------------------------------

llmsfull="$docs/llms-full.txt"
{
  echo "# s1ctl — full documentation"
  echo ""
  echo "> CLI and Go SDK for operating SentinelOne Singularity Platform as code."
  echo "> Pull live state, review in git diff, push back. 370+ commands across REST, SDL, and GraphQL surfaces."
  echo ""
  echo "---"

  seen=""
  grep -oP '\(([^)]+)\)' "$sidebar" | tr -d '()' | while read -r href; do
    case "$href" in http*) continue ;; esac
    file="$docs/$href"
    [ "$href" = "/" ] && file="$docs/README.md"
    [[ "$href" == */ ]] && file="${file}README.md"
    [ -d "$file" ] && file="${file}/README.md"
    [ -f "$file" ] || continue
    real="$(realpath "$file")"
    echo "$seen" | grep -qF "$real" && continue
    seen="${seen}${real}"$'\n'
    echo ""
    echo "---"
    echo ""
    cat "$file"
  done
  echo ""
} > "$llmsfull"

echo "generated: sitemap.xml ($(grep -c '<url>' "$sitemap") URLs), llms.txt ($(wc -l < "$llms") lines), llms-full.txt ($(wc -c < "$llmsfull") bytes)"
