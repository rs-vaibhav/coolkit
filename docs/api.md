# API Reference

**Base URL**: `/api/v1`

## Authentication
Protected routes require a JWT token in the `Authorization` header:
`Authorization: Bearer <your_token>`

## Standard Response Format
All API responses follow this standard JSON structure:

**Success:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error:**
```json
{
  "success": false,
  "error": "Error message description"
}
```

---

## Endpoints

### 1. Health Check
Check if the API is running.

- **Method**: `GET`
- **Path**: `/health`
- **Auth**: 🔓 Public

**Success Response (200 OK)**
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "version": "0.1.0"
  }
}
```

### 2. Database Health Check
Check if the API can connect to the database.

- **Method**: `GET`
- **Path**: `/health/db`
- **Auth**: 🔓 Public

**Success Response (200 OK)**
```json
{
  "success": true,
  "data": {
    "database": "connected"
  }
}
```

### 3. Register User
Create a new user account.

- **Method**: `POST`
- **Path**: `/auth/register`
- **Auth**: 🔓 Public

**Request Body**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "securepassword123"
}
```

**Success Response (201 Created)**
```json
{
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Jane Doe",
    "email": "jane@example.com"
  }
}
```

**Error Response (400 Bad Request)**
```json
{
  "success": false,
  "error": "email is already taken"
}
```

### 4. Login
Authenticate and get a JWT token.

- **Method**: `POST`
- **Path**: `/auth/login`
- **Auth**: 🔓 Public

**Request Body**
```json
{
  "email": "jane@example.com",
  "password": "securepassword123"
}
```

**Success Response (200 OK)**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
  }
}
```

**Error Response (401 Unauthorized)**
```json
{
  "success": false,
  "error": "invalid email or password"
}
```

### 5. Get Current User
Get the profile of the currently authenticated user.

- **Method**: `GET`
- **Path**: `/auth/me`
- **Auth**: 🔒 Protected

**Success Response (200 OK)**
```json
{
  "success": true,
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### 6. Create Club
Create a new club. The creating user automatically becomes the owner.

- **Method**: `POST`
- **Path**: `/clubs`
- **Auth**: 🔒 Protected

**Request Body**
```json
{
  "name": "Computer Science Society",
  "description": "A club for CS enthusiasts."
}
```

**Success Response (201 Created)**
```json
{
  "success": true,
  "data": {
    "id": "987fcdeb-51a2-43d7-9012-345678901234",
    "name": "Computer Science Society",
    "description": "A club for CS enthusiasts.",
    "owner_id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

### 7. List My Clubs
Get a list of all clubs the current user is a member or owner of.

- **Method**: `GET`
- **Path**: `/clubs`
- **Auth**: 🔒 Protected

**Success Response (200 OK)**
```json
{
  "success": true,
  "data": [
    {
      "id": "987fcdeb-51a2-43d7-9012-345678901234",
      "name": "Computer Science Society",
      "description": "A club for CS enthusiasts."
    }
  ]
}
```

### 8. Get Club Details
Retrieve details for a specific club.

- **Method**: `GET`
- **Path**: `/clubs/:id`
- **Auth**: 🔒 Protected

**Success Response (200 OK)**
```json
{
  "success": true,
  "data": {
    "id": "987fcdeb-51a2-43d7-9012-345678901234",
    "name": "Computer Science Society",
    "description": "A club for CS enthusiasts.",
    "owner_id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2024-01-02T10:00:00Z"
  }
}
```

**Error Response (404 Not Found)**
```json
{
  "success": false,
  "error": "club not found"
}
```

### 9. Join Club
Join an existing club as a regular member.

- **Method**: `POST`
- **Path**: `/clubs/:id/join`
- **Auth**: 🔒 Protected

**Success Response (200 OK)**
```json
{
  "success": true,
  "data": {
    "message": "successfully joined club"
  }
}
```

**Error Response (400 Bad Request)**
```json
{
  "success": false,
  "error": "user is already a member of this club"
}
```

### 10. List Club Members
Get a list of all members in a specific club.

- **Method**: `GET`
- **Path**: `/clubs/:id/members`
- **Auth**: 🔒 Protected

**Success Response (200 OK)**
```json
{
  "success": true,
  "data": [
    {
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Jane Doe",
      "role": "owner",
      "joined_at": "2024-01-02T10:00:00Z"
    }
  ]
}
```
