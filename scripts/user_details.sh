#!/usr/bin/env bash
set -euo pipefail

URL="http://localhost:8080/user/details"

# Grab token from login.sh
TOKEN=$(./login.sh | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Failed to retrieve access token via login.sh."
    exit 1
fi

if command -v jq >/dev/null 2>&1; then
    curl -sS -X GET "$URL" -H "Authorization: Bearer $TOKEN" | jq .
else
    curl -sS -X GET "$URL" -H "Authorization: Bearer $TOKEN"
fi
echo
