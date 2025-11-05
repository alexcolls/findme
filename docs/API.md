# FindMe API Documentation

Complete REST API reference for the FindMe backend service.

## Base URL

```
Development: http://localhost:8080/api/v1
Production:  https://api.findme.ai/api/v1
```

## Authentication

FindMe uses JWT (JSON Web Tokens) for authentication. Include the access token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Token Lifecycle

- **Access Token**: 15 minutes expiry
- **Refresh Token**: 7 days expiry
- Tokens are returned on login/registration

## Response Format

### Success Response

```json
{
  "success": true,
  "data": { ... },
  "message": "Operation successful"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": { ... }
  }
}
```

## Error Codes

| Code | HTTP Status | Description |
|------|------------|-------------|
| `UNAUTHORIZED` | 401 | Invalid or expired token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 422 | Invalid input data |
| `CONFLICT` | 409 | Resource already exists |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Authentication Endpoints

### Register User

Create a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe",
  "date_of_birth": "1995-06-15",
  "gender": "male"
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "full_name": "John Doe",
      "verified": false
    },
    "tokens": {
      "access_token": "eyJhbGc...",
      "refresh_token": "eyJhbGc...",
      "expires_in": 900
    }
  }
}
```

### Login

Authenticate existing user.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "full_name": "John Doe",
      "verified": true,
      "video_profile_url": "https://..."
    },
    "tokens": {
      "access_token": "eyJhbGc...",
      "refresh_token": "eyJhbGc...",
      "expires_in": 900
    }
  }
}
```

### Refresh Token

Get new access token using refresh token.

**Endpoint:** `POST /auth/refresh`

**Request Body:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "expires_in": 900
  }
}
```

### Verify Email

Verify user email with token.

**Endpoint:** `POST /auth/verify-email`

**Request Body:**
```json
{
  "token": "verification_token_from_email"
}
```

**Response:** `200 OK`

### Request Password Reset

**Endpoint:** `POST /auth/forgot-password`

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response:** `200 OK`

### Reset Password

**Endpoint:** `POST /auth/reset-password`

**Request Body:**
```json
{
  "token": "reset_token_from_email",
  "password": "NewSecurePass123!"
}
```

**Response:** `200 OK`

---

## User Profile Endpoints

### Get Current User

Get authenticated user's profile.

**Endpoint:** `GET /users/me`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "full_name": "John Doe",
    "date_of_birth": "1995-06-15",
    "gender": "male",
    "bio": "Looking for meaningful connections...",
    "video_profile": {
      "id": "uuid",
      "url": "https://cdn.findme.ai/videos/...",
      "thumbnail_url": "https://cdn.findme.ai/thumbnails/...",
      "duration": 87,
      "verified": true
    },
    "verified": true,
    "active": true,
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

### Update Profile

Update user profile information.

**Endpoint:** `PATCH /users/me`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "full_name": "John Smith",
  "bio": "Updated bio text..."
}
```

**Response:** `200 OK`

### Delete Account

Permanently delete user account.

**Endpoint:** `DELETE /users/me`

**Headers:** `Authorization: Bearer <token>`

**Response:** `204 No Content`

---

## Video Profile Endpoints

### Upload Video Profile

Upload a new video profile.

**Endpoint:** `POST /videos/upload`

**Headers:** 
- `Authorization: Bearer <token>`
- `Content-Type: multipart/form-data`

**Request Body (Multipart):**
```
video: <video_file> (max 50MB, mp4/mov/avi, max 2 minutes)
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "video_id": "uuid",
    "status": "processing",
    "upload_url": "https://...",
    "estimated_processing_time": 60
  }
}
```

### Get Video Status

Check video processing status.

**Endpoint:** `GET /videos/{video_id}/status`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "verified",
    "url": "https://cdn.findme.ai/videos/...",
    "thumbnail_url": "https://cdn.findme.ai/thumbnails/...",
    "duration": 87,
    "verification_score": 0.95,
    "processed_at": "2025-01-15T10:35:00Z"
  }
}
```

**Status Values:**
- `uploading` - Upload in progress
- `processing` - Video being processed
- `verifying` - AI verification in progress
- `verified` - Video approved
- `rejected` - Video rejected (inappropriate content)
- `error` - Processing error

### Delete Video Profile

Delete current video profile.

**Endpoint:** `DELETE /videos/{video_id}`

**Headers:** `Authorization: Bearer <token>`

**Response:** `204 No Content`

---

## Matching Endpoints

### Get Weekly Matches

Get this week's AI-recommended matches.

**Endpoint:** `GET /matches/weekly`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "week_number": 5,
    "year": 2025,
    "matches": [
      {
        "id": "uuid",
        "user": {
          "id": "uuid",
          "full_name": "Jane Doe",
          "age": 28,
          "bio": "Adventure seeker...",
          "video_profile": {
            "thumbnail_url": "https://...",
            "duration": 95
          }
        },
        "match_score": 0.87,
        "status": "pending",
        "expires_at": "2025-01-20T23:59:59Z"
      }
    ]
  }
}
```

### Accept Match

Accept a weekly match recommendation.

**Endpoint:** `POST /matches/{match_id}/accept`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "match_id": "uuid",
    "status": "accepted",
    "mutual_match": false,
    "message": "Waiting for other user to accept"
  }
}
```

### Reject Match

Reject a match recommendation.

**Endpoint:** `POST /matches/{match_id}/reject`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`

### Get Active Matches

Get all active mutual matches.

**Endpoint:** `GET /matches/active`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "matches": [
      {
        "id": "uuid",
        "user": {
          "id": "uuid",
          "full_name": "Jane Doe",
          "video_profile": { ... }
        },
        "matched_at": "2025-01-15T14:00:00Z",
        "can_video_call": true,
        "calls_count": 3,
        "last_call_at": "2025-01-17T19:30:00Z"
      }
    ]
  }
}
```

### End Match

End an active match (weekly decision).

**Endpoint:** `POST /matches/{match_id}/end`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "reason": "not_compatible",
  "feedback": "Optional feedback text"
}
```

**Response:** `200 OK`

---

## Video Call Endpoints

### Initiate Call

Start a video call with a match.

**Endpoint:** `POST /calls/initiate`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "match_id": "uuid"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "call_session_id": "uuid",
    "status": "initiated",
    "ice_servers": [
      {
        "urls": ["stun:stun.l.google.com:19302"]
      },
      {
        "urls": ["turn:turn.findme.ai:3478"],
        "username": "user",
        "credential": "pass"
      }
    ]
  }
}
```

### Accept Call

Accept an incoming call.

**Endpoint:** `POST /calls/{session_id}/accept`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "call_session_id": "uuid",
    "status": "active",
    "started_at": "2025-01-15T20:00:00Z"
  }
}
```

### Reject Call

Reject an incoming call.

**Endpoint:** `POST /calls/{session_id}/reject`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`

### End Call

End an active call.

**Endpoint:** `POST /calls/{session_id}/end`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "call_session_id": "uuid",
    "duration": 1847,
    "ended_at": "2025-01-15T20:30:47Z"
  }
}
```

### Get Call History

Get call history for a match.

**Endpoint:** `GET /calls/history?match_id={match_id}`

**Headers:** `Authorization: Bearer <token>`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "calls": [
      {
        "id": "uuid",
        "started_at": "2025-01-15T20:00:00Z",
        "ended_at": "2025-01-15T20:30:47Z",
        "duration": 1847,
        "status": "completed"
      }
    ],
    "total_duration": 5632,
    "calls_count": 3
  }
}
```

---

## WebSocket Events

Connect to WebSocket for real-time events:

```
ws://localhost:8080/ws?token=<access_token>
```

### Client → Server Events

#### Subscribe to Notifications
```json
{
  "type": "subscribe",
  "channels": ["matches", "calls", "notifications"]
}
```

#### WebRTC Signaling
```json
{
  "type": "webrtc_signal",
  "call_session_id": "uuid",
  "signal": {
    "type": "offer|answer|ice-candidate",
    "data": { ... }
  }
}
```

### Server → Client Events

#### New Match Available
```json
{
  "type": "new_match",
  "data": {
    "match_id": "uuid",
    "user": { ... }
  }
}
```

#### Incoming Call
```json
{
  "type": "incoming_call",
  "data": {
    "call_session_id": "uuid",
    "caller": {
      "id": "uuid",
      "full_name": "Jane Doe",
      "thumbnail_url": "https://..."
    }
  }
}
```

#### Call Status Update
```json
{
  "type": "call_status",
  "data": {
    "call_session_id": "uuid",
    "status": "active|ended|rejected"
  }
}
```

#### WebRTC Signal
```json
{
  "type": "webrtc_signal",
  "data": {
    "call_session_id": "uuid",
    "signal": { ... }
  }
}
```

---

## Rate Limits

| Endpoint Category | Limit |
|------------------|-------|
| Authentication | 10 requests/hour |
| Video Upload | 3 uploads/day |
| API Requests | 60 requests/minute |
| WebSocket Messages | 100 messages/minute |

**Rate Limit Headers:**
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1642348800
```

---

## Pagination

For endpoints returning lists, use pagination parameters:

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

**Response Headers:**
```
X-Total-Count: 150
X-Page: 1
X-Per-Page: 20
X-Total-Pages: 8
```

---

## Testing

### Postman Collection

Import our Postman collection:
```bash
curl -o findme-api.postman_collection.json \
  https://api.findme.ai/docs/postman-collection.json
```

### Example cURL Requests

**Register:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!",
    "full_name": "Test User",
    "date_of_birth": "1995-01-01",
    "gender": "male"
  }'
```

**Get Profile:**
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## Changelog

### v1.0.0 (Planned)
- Initial API release
- Authentication endpoints
- User profile management
- Video upload and processing
- Matching algorithm
- Video call signaling

---

For questions or support, contact: api@findme.ai
