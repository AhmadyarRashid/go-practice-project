# Postman Collection Guide

This folder contains Postman collection and environment files for testing the Go Enterprise API.

## Files

| File | Description |
|------|-------------|
| `Go_Enterprise_API.postman_collection.json` | All API endpoints organized by category |
| `Go_Enterprise_API.postman_environment.json` | Environment variables for local testing |

## How to Import

### Step 1: Import Collection

1. Open Postman
2. Click **Import** button (top left)
3. Drag and drop `Go_Enterprise_API.postman_collection.json`
4. Or click **Upload Files** and select the file

### Step 2: Import Environment

1. Click **Import** button again
2. Drag and drop `Go_Enterprise_API.postman_environment.json`
3. Or click **Upload Files** and select the file

### Step 3: Select Environment

1. Look at top right corner of Postman
2. Click the environment dropdown
3. Select **"Go Enterprise API - Local"**

## How to Test

### 1. Start the Server

```bash
cd /Users/ahmadyar/Desktop/go-enterprise-api
go run cmd/api/main.go
```

### 2. Test Health Endpoint (No Auth Required)

1. In Postman, expand **Health** folder
2. Click on **Health Check**
3. Click **Send**
4. You should see: `{"success":true,"data":{"status":"healthy"}}`

### 3. Register a User

1. Expand **Authentication** folder
2. Click **Register**
3. Modify the email if needed
4. Click **Send**
5. Tokens are automatically saved!

### 4. Test Protected Endpoints

1. Click **Get Current User (Me)**
2. Click **Send**
3. Works because token is automatically used!

## API Testing Flow

```
1. Register or Login
        │
        ▼ (tokens saved automatically)
2. Access Protected Endpoints
        │
        ▼ (token expires after 24 hours)
3. Use Refresh Token endpoint
        │
        ▼ (new tokens saved automatically)
4. Continue using API
```

## Automatic Token Management

The collection has **pre-request scripts** that automatically:

- Save `access_token` after login/register
- Save `refresh_token` for refreshing
- Save `user_id` for user-related requests
- Save `post_id` after creating a post
- Clear tokens after logout

## Environment Variables

| Variable | Description | Auto-Saved |
|----------|-------------|------------|
| `baseUrl` | API base URL | No (preset) |
| `access_token` | JWT access token | Yes |
| `refresh_token` | JWT refresh token | Yes |
| `user_id` | Current user's ID | Yes |
| `post_id` | Last created post ID | Yes |
| `post_slug` | Last created post slug | Yes |

## Folder Structure

```
Collection
├── Health (No Auth)
│   ├── Health Check
│   ├── Readiness Check
│   └── Liveness Check
│
├── Authentication
│   ├── Register
│   ├── Login
│   ├── Get Current User (Me)
│   ├── Refresh Tokens
│   ├── Change Password
│   └── Logout
│
├── Users (Auth Required)
│   ├── Get All Users
│   ├── Get User by ID
│   ├── Search Users
│   ├── Update User
│   ├── Delete User (Admin)
│   ├── Update User Status (Admin)
│   └── Update User Role (Admin)
│
├── Posts
│   ├── Get All Posts (Public)
│   ├── Get Post by ID
│   ├── Get Post by Slug
│   ├── Search Posts (Public)
│   ├── Create Post
│   ├── Get My Posts
│   ├── Update Post
│   └── Delete Post
│
└── Admin
    └── System Info (Admin Only)
```

## Tips

1. **Test Health First**: Always check `/health` endpoint to ensure server is running

2. **Register Before Login**: You need to create an account first

3. **Check Console**: Postman console (View > Show Postman Console) shows token saves

4. **Admin Access**: To test admin endpoints, manually update user role in database:
   ```sql
   UPDATE users SET role = 'admin' WHERE email = 'your@email.com';
   ```

5. **Reset Environment**: If tokens become invalid, clear them:
   - Click the eye icon next to environment dropdown
   - Delete values for `access_token` and `refresh_token`
   - Login again

## Troubleshooting

### "Connection refused"
- Server not running. Start with `go run cmd/api/main.go`

### "Unauthorized"
- Token expired. Use **Refresh Tokens** endpoint
- Or login again

### "Token saved" not appearing
- Check Postman Console for errors
- Make sure environment is selected

### Variables not working
- Ensure environment is selected (top right)
- Check variable names match exactly
