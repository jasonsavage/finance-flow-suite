Agent, you are a backend developer in charge of implementing the transactions feature for the Finance Flow application. 

After login, the user will be able to upload *.csv files from their bank which will be parsed and stored in the database. 

## 1. Database Implementation

The database will need to store the following information:

- transaction_id: VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY
- account_id: VARCHAR(255) NOT NULL
- date: DATE NOT NULL
- description: VARCHAR(255) NOT NULL
- category: VARCHAR(255)
- deposit: DECIMAL(10, 2) NOT NULL
- withdrawal: DECIMAL(10, 2) NOT NULL
- bank_account_name: VARCHAR(255) NOT NULL
- created_at: TIMESTAMP NOT NULL DEFAULT NOW()
- updated_at: TIMESTAMP NOT NULL DEFAULT NOW()


## 2. API Routing & Handlers

The transaction endpoints are exposed under the protected middleware route `/transactions` which requires a valid JWT token to access.

- **`POST /transactions/upload?bankAccount=Chase Checking`**: 
  - **Payload**: `bankAccount` query parameter and *.csv file content in the request body which should use the standard multipart/form-data uploads.
  - **Logic**: Parses the *.csv file content and saves the transactions to the database. Returns a 200 OK upon success along with the number of transactions uploaded.

- **`GET /transactions/list?from=YYYY-MM-DD&to=YYYY-MM-DD`**: 
  - **Payload**: Optional `from` and `to` date query parameters. 
  - **Logic**: Retrieves all transactions for the authenticated user's account from the database and returns them in a JSON format. If `from` or `to` are provided, return transactions within that date range.
  - **Authorization**: Only allow access to the user who owns the account.

  
  ## 3. CSV Parsing Logic

  The CSV parsing logic should be implemented in a way that it can handle different CSV formats from different banks. This utility or group of utilities should be part of the "helpers" package.

  **Step 1**: Generate a unique ID for each transaction by hashing the row data and use it to populate the `transaction_id` field. This will allow future deduping of transactions. 

  **Step 2**: Use the bank account name provided in the query parameter to populate the `bank_account_name` field for each transaction. 

  **Step 3**: Find the Date column in the CSV file. Search the header row for a column name that contains the word "date". Map the column to the `date` field. If more than one column is found, default to the first. 

  **Step 4**: Find the Description column in the CSV file. Search the header row for a column name that contains the word "description". Map the column to the `description` field. If more than one column is found, default to the first. 

  **Step 5**: Detect the bank transaction amount format. Based on the bullet points below, determine the correct way to populate the `deposit` and `withdrawal` columns.

  - Credit/Debit format
    - each transaction has either a `credit` or `debit` column. The `credit` column maps to the `deposit` column and the `debit` column maps to the `withdrawal` column.
    
  - Type/Amount format
    - each transaction has a `type` column which can be either "credit" or "debit" (convert to lowercase before comparison). If the type is "credit", the `deposit` column is populated with the amount. If the type is "debit", the `withdrawal` column is populated with the amount. Both columns `type` and `amount` may contain other words so you will have to search for each column instead of doing s direct match.

  - Positive/Negative format
    - each transaction only has an amount column and the value is prefixed with either a "+" or "-". If the amount is prefixed with a "-", the value should be mapped to the `withdrawal` column otherwise the value should be mapped to the `deposit` column. The prefix should be removed before mapping the value to the column.

**Step 6**: Ignore the `category` column in the CSV file along with any other columns not described above. 

Once the CSV file is parsed successfully, the transactions should be saved to the database. 

## 4. Testing & Validation

Create a shell script in the `/scripts` directory to test the implementation on a local development server without needing a frontend client. 

- **`transaction_upload.sh`**: Should accept 2 arguments: the path to the CSV file and the bank account name. It should then upload the CSV file to the database and return the number of transactions uploaded.

- **`transactions.sh`**: Should fetch all transactions from the authenticated user's account and write them to a file named `transactions.json` in the `/scripts` directory.
  