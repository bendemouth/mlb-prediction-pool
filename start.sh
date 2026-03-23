#!/bin/sh
set -e

# Start the Go backend in the background.
# It listens on localhost:8080 and nginx proxies /api/ to it.
/app/main &

# Start nginx in the foreground (PID 1) so the container lifecycle is tied to it.
exec nginx -g "daemon off;"
