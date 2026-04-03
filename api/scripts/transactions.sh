#!/usr/bin/env bash
set -euo pipefail

# Usage: ./transactions.sh [from_date] [to_date]
# Dates should be in YYYY-MM-DD format.

URL="http://localhost:8080/transactions/list"

# Append optional query params if passed
if [ $# -ge 2 ]; then
    URL="${URL}?from=$1&to=$2"
elif [ $# -eq 1 ]; then
    URL="${URL}?from=$1"
fi

# Grab token from login.sh
TOKEN=$(./login.sh | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Failed to retrieve access token via login.sh."
    exit 1
fi

echo "Fetching transactions from $URL"

if command -v jq >/dev/null 2>&1; then
    curl -sS -X GET "$URL" -H "Authorization: Bearer $TOKEN" | jq .
else
    curl -sS -X GET "$URL" -H "Authorization: Bearer $TOKEN"
fi

echo
