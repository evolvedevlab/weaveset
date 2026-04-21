#!/usr/bin/env bash

FILEPATH="${1:-site/content/list/.changed}"
HUGO_DIR="site"
SLEEP=30

# If running inside container, prefix /app
if [[ -d "/app" ]]; then
    echo "container detected"
    FILEPATH="/app/$FILEPATH"
    HUGO_DIR="/app/$HUGO_DIR"
fi

LAST_MOD=0

while true; do
    if [[ -f "$FILEPATH" ]]; then
        MOD=$(stat -c %Y "$FILEPATH" 2>/dev/null)

        if [[ "$MOD" -gt "$LAST_MOD" ]]; then
            echo "changes detected → rebuilding"

            if ! hugo -s "$HUGO_DIR" --minify --quiet; then
                echo "rebuild error" >&2
                sleep $SLEEP
                continue
            fi

            LAST_MOD=$MOD
        fi
    fi

    sleep $SLEEP
done
