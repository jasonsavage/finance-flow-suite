#!/usr/bin/env bash
set -euo pipefail

# Usage: ./transaction_upload.sh [csv_file] [bank_account_name]
# Requires a running API server, and relies on ./login.sh to fetch a valid token

if [ $# -lt 2 ]; then
  echo "Usage: ./transaction_upload.sh <path_to_csv> <bank_account_name>"
  exit 1
fi

CSV_FILE="$1"
BANK_ACCOUNT="$2"
# Replace spaces with %20 for URL encoding
BANK_ACCOUNT_ENCODED="${BANK_ACCOUNT// /%20}"
URL="http://localhost:8080/transactions/upload?bankAccount=$BANK_ACCOUNT_ENCODED"

# Grab token from login.sh
# Login.sh outputs an extra newline, and the payload is JSON. Let's extract exactly the token.
TOKEN=$(./login.sh | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Failed to retrieve access token via login.sh."
    exit 1
fi

echo "Uploading $CSV_FILE mapped to $BANK_ACCOUNT..."

curl -sS -X POST "$URL" \
    -H "Authorization: Bearer $TOKEN" \
    -F "file=@${CSV_FILE};type=text/csv"

echo
