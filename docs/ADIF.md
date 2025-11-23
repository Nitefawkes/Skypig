# ADIF Import/Export Documentation

Complete guide to ADIF (Amateur Data Interchange Format) import and export functionality in Ham-Radio Cloud.

## Overview

ADIF is the standard interchange format for amateur radio logging data. Ham-Radio Cloud provides full support for importing and exporting logs in ADIF format, enabling:

- **Migration** from other logging software
- **Backup** of your QSO data
- **Synchronization** with LoTW and other services
- **Data sharing** with contest organizers and clubs

**Supported ADIF Version**: 3.1.4

---

## Import

### Import ADIF File

Import QSOs from an ADIF file into your logbook.

**Endpoint**: `POST /api/v1/qsos/import`

**Headers**:
```
Content-Type: text/plain
```

**Request Body**: Raw ADIF content

**Query Parameters**:
- `strict` (boolean, default=false) - If true, import fails on first error

**Example**:
```bash
curl -X POST http://localhost:8080/api/v1/qsos/import \
  -H "Content-Type: text/plain" \
  --data-binary @mylog.adi
```

**Response**: `200 OK` or `206 Partial Content`
```json
{
  "data": {
    "total_records": 100,
    "imported_records": 95,
    "failed_records": 3,
    "skipped_records": 2,
    "errors": [
      "Record 12: callsign is required",
      "Record 45: QSO limit reached (500/500)"
    ]
  },
  "meta": {
    "message": "Imported 95 of 100 records"
  }
}
```

**Status Codes**:
- `200` - All records imported successfully
- `206` - Partial success (some records failed/skipped)
- `400` - Import failed (strict mode only)

### Validation

Validate an ADIF file without importing.

**Endpoint**: `POST /api/v1/qsos/validate`

**Request Body**: Raw ADIF content

**Example**:
```bash
curl -X POST http://localhost:8080/api/v1/qsos/validate \
  -H "Content-Type: text/plain" \
  --data-binary @mylog.adi
```

**Response**: `200 OK` (valid) or `400 Bad Request` (invalid)
```json
{
  "data": {
    "total_records": 100,
    "imported_records": 0,
    "failed_records": 2,
    "errors": [
      "Record 12 (K1ABC): missing required field: TIME_ON",
      "Record 45 (W1AW): invalid date format"
    ]
  },
  "meta": {
    "valid": false,
    "message": "Invalid ADIF: 2 errors found"
  }
}
```

---

## Export

### Export QSOs to ADIF

Export your QSOs in ADIF format.

**Endpoint**: `GET /api/v1/qsos/export`

**Query Parameters** (all optional):
- `callsign` - Filter by callsign (partial match)
- `band` - Filter by band (exact match)
- `mode` - Filter by mode (exact match)
- `country` - Filter by country (partial match)
- `date_from` - Start date (YYYY-MM-DD)
- `date_to` - End date (YYYY-MM-DD)
- `limit` - Maximum number of QSOs to export

**Response Headers**:
```
Content-Type: text/plain; charset=utf-8
Content-Disposition: attachment; filename="hamradio_cloud_export_20251123_140530.adi"
```

**Response Body**: ADIF formatted content

**Examples**:

```bash
# Export all QSOs
curl http://localhost:8080/api/v1/qsos/export -o mylog.adi

# Export only 20m QSOs
curl "http://localhost:8080/api/v1/qsos/export?band=20m" -o 20m_log.adi

# Export QSOs from specific date range
curl "http://localhost:8080/api/v1/qsos/export?date_from=2025-01-01&date_to=2025-12-31" -o 2025_log.adi

# Export QSOs for specific callsign
curl "http://localhost:8080/api/v1/qsos/export?callsign=W1AW" -o w1aw_contacts.adi
```

**Sample Export**:
```
ADIF Export from Ham-Radio Cloud
Generated: 2025-11-23 14:30:00
<adif_ver:5>3.1.4
<programid:17>Ham-Radio Cloud
<programversion:5>1.0.0
<eoh>

<CALL:5>W1AW <QSO_DATE:8>20251123 <TIME_ON:6>140000 <BAND:3>20m <MODE:3>SSB <FREQ:7>14.250 <RST_SENT:2>59 <RST_RCVD:2>57 <NAME:4>John <QTH:10>Boston, MA <GRIDSQUARE:4>FN42 <TX_PWR:3>100 <eor>

<CALL:5>K1ABC <QSO_DATE:8>20251123 <TIME_ON:6>150000 <BAND:3>40m <MODE:2>CW <FREQ:5>7.030 <RST_SENT:3>599 <RST_RCVD:3>579 <eor>
```

---

## Supported Fields

### Required Fields
- `CALL` - Contact's callsign
- `QSO_DATE` - Date in YYYYMMDD format
- `TIME_ON` - Start time in HHMM or HHMMSS format

### Frequency & Band
- `FREQ` - Frequency in MHz
- `FREQ_RX` - Receive frequency (split operation)
- `BAND` - Band (e.g., "20m", "40m")
- `BAND_RX` - Receive band (split operation)

### Mode
- `MODE` - Operating mode (SSB, CW, FT8, etc.)
- `SUBMODE` - Sub-mode (USB, LSB, etc.)

### Signal Reports
- `RST_SENT` - Signal report sent
- `RST_RCVD` - Signal report received

### Contact Information
- `NAME` - Contact's name
- `QTH` - Contact's location
- `GRIDSQUARE` - Maidenhead grid square
- `OPERATOR` - Your callsign
- `STATION_CALLSIGN` - Station callsign

### Location
- `COUNTRY` - Country name
- `DXCC` - DXCC entity number
- `STATE` - State/province
- `COUNTY` - County

### Power & Propagation
- `TX_PWR` - Transmit power in watts
- `RX_PWR` - Receive power in watts
- `PROP_MODE` - Propagation mode (F2, Es, SAT, etc.)

### Satellite
- `SAT_NAME` - Satellite name
- `SAT_MODE` - Satellite mode

### Contest
- `CONTEST_ID` - Contest identifier
- `STX` - Serial number transmitted
- `SRX` - Serial number received

### QSL
- `LOTW_QSL_SENT` - LoTW QSL sent (Y/N/R)
- `LOTW_QSL_RCVD` - LoTW QSL received (Y/N)
- `EQSL_QSL_SENT` - eQSL sent (Y/N)
- `EQSL_QSL_RCVD` - eQSL received (Y/N)

### Comments
- `COMMENT` - Comment field
- `NOTES` - Notes field

---

## Import Behavior

### Validation
All imports are validated before being saved:
- Required fields checked (CALL, QSO_DATE, TIME_ON)
- Callsign format (3-20 characters)
- Date/time parsing
- Band validation
- Frequency ranges
- Power limits

### QSO Limits
Import respects user tier limits:
- **Free tier**: 500 QSOs
- **Operator tier**: 20,000 QSOs
- **Contester tier**: Unlimited

When limit is reached:
- **Non-strict mode**: Skips remaining records, returns partial success
- **Strict mode**: Fails entire import, no records saved

### Error Handling

**Non-strict mode** (default):
- Imports valid records
- Skips invalid records
- Returns summary with errors
- HTTP 206 Partial Content

**Strict mode** (`?strict=true`):
- Stops at first error
- No records imported
- Returns error immediately
- HTTP 400 Bad Request

---

## Export Options

### Filtering
All standard QSO filters apply:
- Callsign (partial match)
- Band (exact match)
- Mode (exact match)
- Country (partial match)
- Date range

### File Format
Exports include:
- ADIF header with version
- Program identification
- Generation timestamp
- All QSO records
- End-of-record markers

---

## Testing

### Test Script
Use the provided test script to verify ADIF functionality:

```bash
./scripts/test-adif.sh
```

This tests:
- ✅ ADIF validation
- ✅ Import (successful)
- ✅ Import (partial/errors)
- ✅ Export (all records)
- ✅ Export (filtered)
- ✅ Error handling

### Sample Data
A sample ADIF file is provided at:
```
backend/testdata/sample.adi
```

Contains 5 sample QSOs for testing.

---

## Common Use Cases

### 1. Migrate from Another Logging Software

```bash
# Export from old software to ADIF
# Then import to Ham-Radio Cloud
curl -X POST http://localhost:8080/api/v1/qsos/import \
  -H "Content-Type: text/plain" \
  --data-binary @export_from_old_software.adi
```

### 2. Backup Your Logs

```bash
# Export all QSOs
curl http://localhost:8080/api/v1/qsos/export -o backup_$(date +%Y%m%d).adi
```

### 3. Prepare for LoTW Upload

```bash
# Export recent QSOs for LoTW
curl "http://localhost:8080/api/v1/qsos/export?date_from=2025-01-01" -o lotw_upload.adi
```

### 4. Extract Contest Log

```bash
# Export QSOs from specific contest
curl "http://localhost:8080/api/v1/qsos/export?date_from=2025-11-23&date_to=2025-11-24" \
  -o contest_log.adi
```

### 5. Share DXpedition Log

```bash
# Export QSOs from specific location
curl "http://localhost:8080/api/v1/qsos/export?gridsquare=FN42" -o dxpedition.adi
```

---

## Validation Rules

### Callsign
- 3-20 characters
- Alphanumeric
- Automatically converted to uppercase

### Date/Time
- `QSO_DATE`: YYYYMMDD format
- `TIME_ON`: HHMM or HHMMSS format
- `TIME_OFF`: Optional, must be after TIME_ON

### Band
- Valid amateur radio bands
- Examples: 160m, 80m, 40m, 20m, 10m, 2m, 70cm

### Frequency
- 0 - 300,000 MHz (0-300 GHz)
- Decimal format supported

### Power
- 0 - 10,000 watts (10 kW maximum)

---

## Error Codes

| Code | Description |
|------|-------------|
| `EMPTY_CONTENT` | Request body is empty |
| `IMPORT_FAILED` | Import failed (strict mode) |
| `EXPORT_FAILED` | Export operation failed |
| `INVALID_ADIF` | ADIF parsing error |

---

## Limitations

### Current
- Maximum file size: Limited by server configuration
- Batch size: All records processed in single transaction
- Format: ADIF 3.1.4 only (no ADX/XML support yet)

### Future Enhancements
- [ ] Chunked import for large files
- [ ] Progress tracking for long imports
- [ ] ADX (ADIF XML) format support
- [ ] Cabrillo format import/export
- [ ] Duplicate detection
- [ ] Merge strategies (update vs skip vs error)

---

## Implementation Details

### Parser
- **Location**: `backend/pkg/adif/parser.go`
- **Strategy**: Regex-based field extraction
- **Modes**: Strict and non-strict
- **Encoding**: UTF-8

### Generator
- **Location**: `backend/pkg/adif/generator.go`
- **Format**: ADIF 3.1.4
- **Header**: Includes program ID and version
- **Fields**: All non-empty fields included

### Service
- **Location**: `backend/internal/services/adif_service.go`
- **Validation**: Comprehensive field checking
- **Limits**: Tier-based QSO limits enforced
- **Transactions**: Atomic per-QSO (no batch rollback)

---

## Troubleshooting

### Import Fails with "missing required field"
- Ensure CALL, QSO_DATE, and TIME_ON are present
- Check field format (YYYYMMDD for dates, HHMM for times)

### Some Records Skipped
- Check QSO limit for your tier
- Review error messages in response
- Validate ADIF file with `/validate` endpoint

### Export is Empty
- Verify QSOs exist with `GET /api/v1/qsos`
- Check filter parameters
- Ensure authentication (when implemented)

### Invalid ADIF Format
- Use ADIF 3.1.4 specification
- Ensure `<eor>` markers are present
- Check field length specifications match actual length

---

## Further Reading

- [ADIF Specification](http://adif.org/)
- [QSO API Documentation](./QSO_API.md)
- [Development Guide](./DEVELOPMENT.md)

---

**Next Steps**: See [LOTW.md](./LOTW.md) (future) for LoTW synchronization using ADIF.
