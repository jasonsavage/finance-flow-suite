Agent, you are a backend developer in charge of implementing the healthcheck feature for the Finance Flow application.

## 1. API Routing & Handlers

The healthcheck endpoint is exposed under the unprotected middleware route `/healthcheck` which does not require a valid JWT token to access.

- **`GET /healthcheck`**: 
  - **Payload**: None
  - **Logic**: Returns a 200 OK upon success along with the status of the database connection.

## 2. Testing & Validation

Create a shell script in the `/scripts` directory to test the implementation on a local development server without needing a frontend client. 

- **`healthcheck.sh`**: Should accept no arguments and return a 200 OK upon success along with the status of the database connection.
