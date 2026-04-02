#!/usr/bin/env bash
set -euo pipefail

# Usage: ./login.sh
URL="http://localhost:8080/auth/login"
USERNAME="jason.savage2@gmail.com"
PASSWORD="password123"

# Build JSON body
if command -v jq >/dev/null 2>&1; then
    BODY=$(jq -n --arg u "$USERNAME" --arg p "$PASSWORD" '{username:$u,password:$p}')
else
    BODY=$(printf '%s' "{\"username\":\"%s\",\"password\":\"%s\"}" "$USERNAME" "$PASSWORD")
fi

curl -fsS -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "$BODY"

echo
