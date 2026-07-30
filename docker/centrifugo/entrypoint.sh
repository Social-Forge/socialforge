#!/bin/sh
set -e

# Use /tmp for config (writable directory)
CENTRIFUGO_CONFIG_PATH="/tmp/config.json"

# Check if template exists
if [ ! -f "/centrifugo/config.json" ]; then
    echo "Error: config.json template not found"
    exit 1
fi

# Replace environment variables in config using envsubst
envsubst < /centrifugo/config.json > "$CENTRIFUGO_CONFIG_PATH"

# Start Centrifugo with the generated config
exec centrifugo --config="$CENTRIFUGO_CONFIG_PATH"
