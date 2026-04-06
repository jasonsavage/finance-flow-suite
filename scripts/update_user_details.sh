#!/usr/bin/env bash
set -euo pipefail

if [ $# -lt 3 ]; then
  echo "Usage: ./update_user_details.sh <first_name> <last_name> <email>"
  exit 1
fi

FIRST_NAME="$1"
LAST_NAME="$2"
EMAIL="$3"
URL="http://localhost:8080/user/details"

TOKEN=$(./login.sh | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Failed to retrieve access token via login.sh."
    exit 1
fi

if command -v jq >/dev/null 2>&1; then
    BODY=$(jq -n --arg f "$FIRST_NAME" --arg l "$LAST_NAME" --arg e "$EMAIL" '{first_name:$f,last_name:$l,email:$e}')
else
    BODY=$(printf '%s' "{\"first_name\":\"%s\",\"last_name\":\"%s\",\"email\":\"%s\"}" "$FIRST_NAME" "$LAST_NAME" "$EMAIL")
fi

echo "Updating user to: $FIRST_NAME $LAST_NAME, $EMAIL"

curl -sS -X PUT "$URL" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$BODY"

echo
