# QSO API Documentation

Complete API documentation for QSO (contact log) endpoints in Ham-Radio Cloud.

## Overview

The QSO API allows ham radio operators to create, read, update, and delete contact logs. All endpoints require authentication (to be implemented with OAuth).

**Base URL**: `http://localhost:8080/api/v1`

## Endpoints

### 1. Create QSO

Create a new contact log entry.

**Endpoint**: `POST /qsos`

**Request Body**:
```json
{
  "callsign": "K1ABC",
  "time_on": "2025-11-23T14:30:00Z",
  "time_off": "2025-11-23T14:45:00Z",
  "band": "20m",
  "mode": "SSB",
  "freq": 14.250,
  "rst_sent": "59",
  "rst_rcvd": "57",
  "name": "John",
  "qth": "Boston, MA",
  "gridsquare": "FN42",
  "country": "USA",
  "state": "MA",
  "tx_pwr": 100,
  "comment": "Great QSO!"
}
```

**Required Fields**:
- `callsign` (string, 3-20 chars) - Contact's callsign
- `time_on` (ISO 8601 timestamp) - Start time of contact

**Optional Fields**:
- `time_off` - End time of contact
- `band` - Amateur radio band (e.g., "20m", "40m")
- `mode` - Operating mode (e.g., "SSB", "CW", "FT8")
- `freq` - Frequency in MHz
- `rst_sent/rst_rcvd` - Signal reports
- `name`, `qth`, `gridsquare` - Contact's info
- `tx_pwr` - Transmit power in watts
- And many more ADIF-compliant fields...

**Response**: `201 Created`
```json
{
  "data": {
    "id": 1,
    "user_id": 1,
    "callsign": "K1ABC",
    "time_on": "2025-11-23T14:30:00Z",
    // ... all fields
    "created_at": "2025-11-23T14:30:05Z",
    "updated_at": "2025-11-23T14:30:05Z"
  }
}
```

**Errors**:
- `400` - Invalid request body or validation failed
- `400` - QSO limit reached for user's tier

---

### 2. List QSOs

Retrieve a list of QSOs with optional filtering.

**Endpoint**: `GET /qsos`

**Query Parameters**:
- `callsign` (string) - Filter by callsign (partial match)
- `band` (string) - Filter by band (exact match)
- `mode` (string) - Filter by mode (exact match)
- `country` (string) - Filter by country (partial match)
- `date_from` (YYYY-MM-DD) - Start date filter
- `date_to` (YYYY-MM-DD) - End date filter
- `limit` (int, default=50) - Number of results per page
- `offset` (int, default=0) - Pagination offset

**Examples**:
```bash
# Get all QSOs
GET /qsos

# Filter by band
GET /qsos?band=20m

# Filter by callsign
GET /qsos?callsign=K1ABC

# Date range
GET /qsos?date_from=2025-01-01&date_to=2025-12-31

# Pagination
GET /qsos?limit=25&offset=50
```

**Response**: `200 OK`
```json
{
  "data": [
    {
      "id": 1,
      "callsign": "K1ABC",
      // ... all QSO fields
    },
    // ... more QSOs
  ],
  "meta": {
    "total": 150,
    "limit": 50,
    "offset": 0,
    "returned": 50
  }
}
```

---

### 3. Get QSO by ID

Retrieve a specific QSO.

**Endpoint**: `GET /qsos/:id`

**Parameters**:
- `id` (path, int) - QSO ID

**Response**: `200 OK`
```json
{
  "data": {
    "id": 1,
    "user_id": 1,
    "callsign": "K1ABC",
    // ... all fields
  }
}
```

**Errors**:
- `400` - Invalid ID
- `404` - QSO not found or doesn't belong to user

---

### 4. Update QSO

Update an existing QSO.

**Endpoint**: `PUT /qsos/:id`

**Parameters**:
- `id` (path, int) - QSO ID

**Request Body**: Same as Create QSO (all fields)

**Response**: `200 OK`
```json
{
  "data": {
    "id": 1,
    // ... updated fields
    "updated_at": "2025-11-23T15:00:00Z"
  }
}
```

**Errors**:
- `400` - Invalid request body or validation failed
- `404` - QSO not found

---

### 5. Delete QSO

Delete a QSO.

**Endpoint**: `DELETE /qsos/:id`

**Parameters**:
- `id` (path, int) - QSO ID

**Response**: `204 No Content`

**Errors**:
- `400` - Invalid ID
- `404` - QSO not found

**Note**: Deleting a QSO decrements the user's QSO count.

---

### 6. Get QSO Statistics

Get user's QSO statistics and limits.

**Endpoint**: `GET /qsos/stats`

**Response**: `200 OK`
```json
{
  "data": {
    "total_qsos": 42,
    "qso_limit": 20000,
    "remaining_qsos": 19958
  }
}
```

**Note**: For "contester" tier, `remaining_qsos` returns -1 (unlimited).

---

## Validation Rules

### Callsign
- Required
- 3-20 characters
- Automatically converted to uppercase
- Alphanumeric characters

### Time
- `time_on` is required
- `time_off` must be after `time_on` (if provided)
- ISO 8601 format

### Band
- Must be a valid amateur radio band
- Supported: 160m, 80m, 60m, 40m, 30m, 20m, 17m, 15m, 12m, 10m, 6m, 2m, 70cm, etc.

### Frequency
- 0 - 300000 MHz (0-300 GHz)
- Decimal format

### Power
- 0 - 10000 watts (10 kW max)
- Decimal format

### Grid Square
- 4-8 characters
- Maidenhead grid square format
- Automatically converted to uppercase

---

## Error Responses

All errors follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": "Optional additional details"
  }
}
```

### Common Error Codes

- `INVALID_REQUEST` - Malformed request body
- `INVALID_ID` - Invalid QSO ID
- `NOT_FOUND` - QSO not found
- `CREATE_FAILED` - Failed to create QSO
- `UPDATE_FAILED` - Failed to update QSO
- `DELETE_FAILED` - Failed to delete QSO
- `STATS_FAILED` - Failed to retrieve stats

---

## Rate Limiting

(To be implemented)

Tier-based rate limiting:
- **Free**: 100 requests/hour
- **Operator**: 500 requests/hour
- **Contester**: 1000 requests/hour

---

## Authentication

(To be implemented - Phase 2)

All QSO endpoints will require authentication via JWT token:

```bash
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/qsos
```

---

## Testing

Use the provided test script to test all endpoints:

```bash
./scripts/test-qso-api.sh
```

Or manually with curl:

```bash
# Create a QSO
curl -X POST http://localhost:8080/api/v1/qsos \
  -H "Content-Type: application/json" \
  -d '{
    "callsign": "K1ABC",
    "time_on": "2025-11-23T14:30:00Z",
    "band": "20m",
    "mode": "SSB"
  }'

# List QSOs
curl http://localhost:8080/api/v1/qsos

# Get stats
curl http://localhost:8080/api/v1/qsos/stats
```

---

## Implementation Details

### Database
- QSOs stored in TimescaleDB hypertable optimized for time-series queries
- Automatic triggers update user's QSO count
- Indexes on commonly queried fields (callsign, band, mode, time_on)

### Architecture
- **Repository Layer**: `internal/database/qso_repository.go`
- **Service Layer**: `internal/services/qso_service.go`
- **Handler Layer**: `internal/handlers/qso_handler.go`

### ADIF Compliance
The QSO model supports all major ADIF (Amateur Data Interchange Format) fields, making it compatible with standard amateur radio logging software.

---

## Next Steps

1. **ADIF Import/Export** - Bulk import/export in ADIF format
2. **OAuth Authentication** - QRZ.com integration
3. **LoTW Sync** - Automatic Logbook of the World synchronization
4. **Advanced Filtering** - More complex queries (e.g., awards tracking)
5. **WebSocket Updates** - Real-time QSO updates

---

For more information, see:
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System design
- [DEVELOPMENT.md](./DEVELOPMENT.md) - Development guide
