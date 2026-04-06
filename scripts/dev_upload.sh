#!/usr/bin/env bash
set -euo pipefail

./transaction_upload.sh "../local/2026-04_Capital_One_Billing.csv" "CapitalOne Billing"
./transaction_upload.sh "../local/2026-04_Capital_One_Venture.csv" "CapitalOne Venture"
./transaction_upload.sh "../local/2026-04_PNC_Growth_x7718.csv" "PNC Growth"
./transaction_upload.sh "../local/2026-04_PNC_Reserve_x7697.csv" "PNC Reserve"
./transaction_upload.sh "../local/2026-04_PNC_Spend_x8276.csv" "PNC Spend"

./transactions.sh