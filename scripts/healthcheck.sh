#!/usr/bin/env bash
set -euo pipefail

URL="http://localhost:8080/healthcheck"

echo "Checking health of API server at $URL"
curl -isS -X GET "$URL"

echo
echo
