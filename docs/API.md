# Ham Radio Cloud API Documentation

**Version:** 1.0.0
**Base URL:** `http://localhost:8080/api/v1`
**Protocol:** REST
**Authentication:** JWT Bearer Token (OAuth 2.0 via QRZ.com)

---

## Table of Contents

1. [Authentication](#authentication)
2. [Health & Status](#health--status)
3. [QSO Endpoints](#qso-endpoints)
4. [Propagation Endpoints](#propagation-endpoints)
5. [User Endpoints](#user-endpoints)
6. [Error Responses](#error-responses)

---

## Authentication

### OAuth Flow (QRZ.com)

#### 1. Initiate Login
```http
GET /api/v1/auth/login/qrz
```

Redirects user to QRZ.com OAuth consent page.

#### 2. OAuth Callback
```http
GET /api/v1/auth/callback/qrz?code={auth_code}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "callsign": "W1AW",
    "email": "user@example.com",
    "name": "Test User",
    "tier": "operator"
  }
}
```

#### 3. Using the Token

Include the JWT token in all authenticated requests:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## Health & Status

### Health Check
```http
GET /health
GET /api/v1/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-11-23T10:30:00Z",
  "service": "ham-radio-cloud-api",
  "version": "1.0.0"
}
```

---

## QSO Endpoints

### List QSOs
```http
GET /api/v1/qso
```

**Query Parameters:**
- `callsign` (string): Filter by callsign
- `band` (string): Filter by band (e.g., "20m", "40m")
- `mode` (string): Filter by mode (e.g., "SSB", "CW", "FT8")
- `start_date` (ISO 8601): Filter by start date
- `end_date` (ISO 8601): Filter by end date
- `limit` (int): Results per page (default: 100, max: 1000)
- `offset` (int): Pagination offset

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "callsign": "K1ABC",
      "frequency": 14.074,
      "band": "20m",
      "mode": "FT8",
      "rst_sent": "-10",
      "rst_received": "-12",
      "qso_date": "2025-11-23",
      "time_on": "2025-11-23T14:30:00Z",
      "time_off": "2025-11-23T14:31:00Z",
      "grid_square": "FN42",
      "country": "USA",
      "state": "MA",
      "comment": "Nice contact!",
      "lotw_sent": false,
      "lotw_confirmed": false,
      "created_at": "2025-11-23T14:32:00Z",
      "updated_at": "2025-11-23T14:32:00Z"
    }
  ],
  "total": 1,
  "limit": 100,
  "offset": 0
}
```

### Create QSO
```http
POST /api/v1/qso
Content-Type: application/json
```

**Request Body:**
```json
{
  "callsign": "K1ABC",
  "frequency": 14.074,
  "band": "20m",
  "mode": "FT8",
  "rst_sent": "-10",
  "rst_received": "-12",
  "qso_date": "2025-11-23",
  "time_on": "2025-11-23T14:30:00Z",
  "time_off": "2025-11-23T14:31:00Z",
  "grid_square": "FN42",
  "country": "USA",
  "state": "MA",
  "comment": "Nice contact!",
  "tx_power": 100
}
```

**Response:** 201 Created
```json
{
  "id": "uuid",
  "message": "QSO created successfully"
}
```

### Update QSO
```http
PUT /api/v1/qso/{id}
Content-Type: application/json
```

**Request Body:** Same as Create QSO

**Response:** 200 OK

### Delete QSO
```http
DELETE /api/v1/qso/{id}
```

**Response:** 204 No Content

### Import ADIF
```http
POST /api/v1/qso/import/adif
Content-Type: multipart/form-data
```

**Request Body:**
- `file`: ADIF file (.adi or .adif)

**Response:**
```json
{
  "imported": 150,
  "skipped": 5,
  "errors": []
}
```

### Export ADIF
```http
GET /api/v1/qso/export/adif
```

**Query Parameters:** Same as List QSOs (for filtering)

**Response:** ADIF file download

---

## Propagation Endpoints

### Current Propagation Data
```http
GET /api/v1/propagation/current
```

**Response:**
```json
{
  "timestamp": "2025-11-23T10:00:00Z",
  "solar_flux": 150,
  "sunspot_number": 85,
  "a_index": 5,
  "k_index": 2,
  "xray_flux": "B1.2",
  "helium_line": 0.0,
  "proton_flux": 0,
  "electron_flux": 0,
  "source": "NOAA SWPC"
}
```

### Band Conditions
```http
GET /api/v1/propagation/bands
```

**Query Parameters:**
- `grid_square` (string): Observer grid square (e.g., "FN42")

**Response:**
```json
{
  "timestamp": "2025-11-23T10:00:00Z",
  "conditions": [
    {
      "band": "160m",
      "condition": "poor",
      "score": 2,
      "day_night": "day"
    },
    {
      "band": "80m",
      "condition": "fair",
      "score": 5,
      "day_night": "day"
    },
    {
      "band": "40m",
      "condition": "good",
      "score": 8,
      "day_night": "day"
    },
    {
      "band": "20m",
      "condition": "good",
      "score": 9,
      "day_night": "day"
    },
    {
      "band": "15m",
      "condition": "fair",
      "score": 6,
      "day_night": "day"
    },
    {
      "band": "10m",
      "condition": "poor",
      "score": 3,
      "day_night": "day"
    }
  ]
}
```

---

## User Endpoints

### Get Profile
```http
GET /api/v1/user/profile
```

**Response:**
```json
{
  "id": "uuid",
  "callsign": "W1AW",
  "email": "user@example.com",
  "name": "Test User",
  "tier": "operator",
  "created_at": "2025-11-01T00:00:00Z",
  "updated_at": "2025-11-23T10:00:00Z"
}
```

### Get User Settings
```http
GET /api/v1/user/settings
```

**Response:**
```json
{
  "user_id": "uuid",
  "lotw_enabled": true,
  "lotw_username": "w1aw",
  "propagation_alerts": true,
  "preferred_bands": ["20m", "40m"],
  "grid_square": "FN42",
  "time_zone": "America/New_York"
}
```

### Update User Settings
```http
PUT /api/v1/user/settings
Content-Type: application/json
```

**Request Body:**
```json
{
  "lotw_enabled": true,
  "lotw_username": "w1aw",
  "propagation_alerts": true,
  "preferred_bands": ["20m", "40m", "15m"],
  "grid_square": "FN42",
  "time_zone": "America/New_York"
}
```

**Response:** 200 OK

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "error_type",
  "message": "Human-readable error message",
  "code": 400
}
```

### Common HTTP Status Codes

- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `204 No Content` - Success with no response body
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

### Example Error Response

```json
{
  "error": "validation_error",
  "message": "Invalid callsign format",
  "code": 400
}
```

---

## Rate Limiting

- **Free Tier:** 100 requests/hour
- **Operator Tier:** 1,000 requests/hour
- **Contester Tier:** 10,000 requests/hour

Rate limit headers:
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1637683200
```

---

## Webhooks (Future)

Webhook support will be added in Phase 5 for:
- New QSO confirmations (LoTW)
- Propagation alerts
- DX cluster spots

---

## GraphQL (Future)

A GraphQL endpoint will be available at `/api/graphql` in a future release for more flexible querying.

---

*Last Updated: 2025-11-23*
*API Version: 1.0.0*
