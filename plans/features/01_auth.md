# Feature: 01 Basic Authentication Flow

This document details the complete implementation of the basic authorization flow for the Go backend application. It outlines the database, routing, and architectural decisions made to ensure a secure and stateless authentication system.

## 1. Database Implementation

A `users` table has been created via the `001_create_tables.sql` migration, containing the following fields:
- `id`: `VARCHAR(255) PRIMARY KEY`. A unique identifier generated automatically during user registration prefixed with `USR-` (e.g., `USR-123...`).
- `username`: `VARCHAR(255) NOT NULL UNIQUE`. Serves as the primary login identifier (often an email).
- `password_hash`: `VARCHAR(255) NOT NULL`. Stores the securely hashed equivalent of the user's password using `bcrypt`.
- `first_name`: `VARCHAR(255) NOT NULL`. Initialized as an empty string during registration.
- `last_name`: `VARCHAR(255) NOT NULL`. Initialized as an empty string during registration.
- `email`: `VARCHAR(255) NOT NULL UNIQUE`. Automatically defaults to the provided `username` during registration to satisfy the `UNIQUE` and `NOT NULL` constraints, allowing users to update it later.
- `created_at`: `TIMESTAMP NOT NULL DEFAULT NOW()`.
- `updated_at`: `TIMESTAMP NOT NULL DEFAULT NOW()`.

## 2. API Routing & Handlers

The authentication endpoints are exposed under public `/auth` routes, while all other business logic is placed behind protected middleware routes.

### Public Routes
- **`POST /auth/register`**: 
  - **Payload**: Requires `username` and `password`.
  - **Logic**: Securely hashes the password with `bcrypt`, generates a `id`, copies `username` to the `email` field, and saves the new user to the `users` database table using the repository pattern. Returns a session token upon success.
- **`POST /auth/login`**: 
  - **Payload**: Requires `username` and `password`.
  - **Logic**: Queries the `users` table for the matching `username`. Verifies the hashed password utilizing `bcrypt`. Returns a session token and the user's data (excluding the password hash).

### Protected Routes & Middleware
- **Middleware Check**: We utilize a `RequireAuth` middleware that intercepts requests made to protected routes.
- **JWT Verification**: The middleware expects an `Authorization: Bearer <token>` header. It validates the cryptographic signature of the token using the `JWT_SECRET` environment variable.
- **Context Injection**: Successfully authenticated requests have their numerical `userID` safely injected into the request `Context`, allowing handlers to identify the active user.

## 3. Session Management (Stateless JWT)

This application employs stateless JSON Web Tokens (JWT) for authentication.
- Tokens securely bundle claims (such as the `user_id` and expiration timestamp) into a verifiable signature. 
- Because tokens are self-contained and stateless, the API does not need to execute a database look-up to verify authenticity per-request, drastically improving performance.
- We utilize **client-side logout**. A dedicated `/logout` endpoint is deliberately *not* implemented to maintain the stateless architecture. Clients are simply responsible for discarding the token from storage when a user logs out.

## 4. Testing & Validation

Two shell scripts have been created in the `/scripts` directory, fully equipped to easily test the implementation on a local development server without needing a frontend client:
- **`register.sh`**: Customizes the necessary `username` and `password` variables and pushes a `POST` request to register a new user in the local database.
- **`login.sh`**: Authenticates an existing user and prints the resulting JSON Web Token back to the terminal, allowing developers to manually inspect or pass the token to subsequent curl requests.