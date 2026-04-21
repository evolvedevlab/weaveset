#!/bin/sh
set -e

PUID=${PUID:-1000}
PGID=${PGID:-1000}

# Create group and user
addgroup -g "$PGID" appuser 2>/dev/null || true
adduser -D -u "$PUID" -G appuser -s /bin/sh appuser 2>/dev/null || true

# Fix all permissions in /app
chown -R "$PUID:$PGID" /app

# Switch user and run
exec su-exec appuser "$@"
