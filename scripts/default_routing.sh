#!/bin/bash

REDIS_HOST="localhost"
REDIS_PORT="6379"

CONFIG_FILE="default_routing.json"
HASH_KEY="routing-distributed-cache"
PREFIX_ROUTING_CFG_KEY="routing:cfg:"

redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" DEL "$HASH_KEY"

jq -c '.[]' "$CONFIG_FILE" | while read -r data; do
    SOURCE_URL=$(jq -r '.sourceURL' <<< "$data")
    METHOD_HTTP=$(jq -r '.methodHTTP' <<< "$data")

    KEY="${PREFIX_ROUTING_CFG_KEY}${METHOD_HTTP}|${SOURCE_URL}"

    redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" HSET "$HASH_KEY" "$KEY" "$data"

    if [ $? -eq 0 ]; then
        echo "new routing cfg added to Redis: $KEY"
    else
        echo "error adding routing cfg to Redis: $KEY"
    fi
done