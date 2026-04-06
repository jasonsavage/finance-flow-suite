Agent, you are a backend developer in charge of implementing the user details feature for the Finance Flow application.

Once a user has registered for an account, they will be able to edit there first name, last name, and email address. We need to create routes that allow the user to update their account details.

## 1. API Routing & Handlers

- **`GET /user/details`**: 
  - **Payload**: None
  - **Logic**: Returns the authenticated user's details. Do not include the password hash, created_at, or updated_at in the response.
  
- **`PUT /user/details`**: 
  - **Payload**: 
    - `first_name`: String
    - `last_name`: String
    - `email`: String
  - **Logic**: Updates the user's details in the database and returns a 200 OK upon success. The only fields that can be updated are first_name, last_name, and email. If the user's email is not unique or invalid, return a validation error.

## . Testing & Validation

Create a shell script in the `/scripts` directory to test the implementation on a local development server without needing a frontend client. 

- **`user_details.sh`**: Should make a GET request to `user/details` and echo the JSON response.
- **`update_user_details.sh`**: Should accept first_name, last_name, and email as arguments and make a PUT request to `user/details` endpoint.


