# Testing Phase 2 - Cloud Logbook Core

This guide helps you test the QSO logging functionality implemented in Phase 2.

## Prerequisites

1. Docker and Docker Compose installed
2. Backend and database running
3. Frontend dev server running

## Quick Start

```bash
# Terminal 1: Start database and backend
cd Skypig
docker-compose -f deployments/docker/docker-compose.yml up postgres

# Terminal 2: Start backend (in another terminal)
cd backend
go run cmd/api/main.go

# Terminal 3: Start frontend (in another terminal)
cd frontend
npm run dev
```

**Access the app:** http://localhost:5173

---

## Test Scenarios

### 1. Create a QSO Manually

**Steps:**
1. Navigate to http://localhost:5173/logbook
2. Click "Log QSO" button
3. Fill in the form:
   - Callsign: `K1ABC`
   - Frequency: `14.074`
   - Band: `20m`
   - Mode: `FT8`
   - RST Sent: `-10`
   - RST Received: `-12`
   - Grid Square: `FN42`
4. Click "Log QSO"
5. Verify success message appears
6. Verify QSO appears in the list

**Expected Result:**
- Success notification: "QSO logged successfully!"
- QSO appears in table with all fields populated
- Callsign is uppercased automatically

---

### 2. Filter QSOs

**Steps:**
1. Create multiple QSOs with different bands/modes
2. Use the filter dropdowns:
   - Select Band: `20m`
   - Leave other filters empty
3. Click outside the dropdown or press Enter

**Expected Result:**
- Only QSOs on 20m band are displayed
- Other bands are filtered out
- "Clear All" button clears filters

---

### 3. Edit a QSO

**Steps:**
1. Find a QSO in the list
2. Click "Edit" button
3. Modify a field (e.g., change comment)
4. Click "Update QSO"

**Expected Result:**
- Success notification: "QSO updated successfully!"
- Changes are reflected in the table
- Modal closes automatically

---

### 4. Delete a QSO

**Steps:**
1. Find a QSO in the list
2. Click "Delete" button
3. Confirm deletion in the browser prompt

**Expected Result:**
- Confirmation dialog appears
- After confirming: Success notification
- QSO is removed from the list
- Database record is deleted

---

### 5. Import ADIF File

**Preparation:**
Create a test ADIF file (`test.adi`):

```
ADIF Export Test
<ADIF_VER:5>3.1.0
<PROGRAMID:4>Test
<EOH>

<CALL:5>W1AW <QSO_DATE:8>20231120 <TIME_ON:6>143000 <BAND:3>20m <MODE:3>CW <RST_SENT:3>599 <RST_RCVD:3>599 <FREQ:8>14.050000 <GRIDSQUARE:4>FN31 <EOR>
<CALL:5>K2XYZ <QSO_DATE:8>20231120 <TIME_ON:6>150000 <BAND:3>40m <MODE:3>SSB <RST_SENT:2>59 <RST_RCVD:2>57 <FREQ:7>7.200000 <GRIDSQUARE:4>FN42 <EOR>
<CALL:6>N3TEST <QSO_DATE:8>20231121 <TIME_ON:6>120000 <BAND:3>15m <MODE:4>FT8 <RST_SENT:3>-10 <RST_RCVD:3>-12 <FREQ:8>21.074000 <EOR>
```

**Steps:**
1. Click "Import ADIF" button
2. Select the `test.adi` file
3. Wait for import to complete

**Expected Result:**
- Success message: "Imported 3 QSOs (0 skipped)"
- All 3 QSOs appear in the logbook
- Fields are parsed correctly:
  - Callsigns uppercase
  - Dates/times converted properly
  - Frequencies and bands match

---

### 6. Export ADIF File

**Steps:**
1. Create some QSOs or import test data
2. Optionally apply filters (e.g., Band: `20m`)
3. Click "Export ADIF" button

**Expected Result:**
- Browser downloads a file named `logbook_YYYY-MM-DD.adi`
- File contains only filtered QSOs (if filters applied)
- File format is valid ADIF 3.1.0
- Can be re-imported to verify roundtrip

---

### 7. Date Range Filtering

**Steps:**
1. Create QSOs on different dates
2. Set Start Date: `2023-11-01`
3. Set End Date: `2023-11-30`
4. Filter should apply automatically

**Expected Result:**
- Only QSOs within the date range are shown
- QSOs outside the range are hidden
- Combining with other filters (band/mode) works correctly

---

### 8. Callsign Search

**Steps:**
1. Create QSOs with different callsigns
2. In Callsign filter, type: `W1`
3. Press Enter or click outside

**Expected Result:**
- Only QSOs with callsigns containing "W1" are shown
- Search is case-insensitive
- Partial matches work (e.g., `W1ABC`, `W1TEST` both match)

---

## API Testing (curl)

### Create QSO
```bash
curl -X POST http://localhost:8080/api/v1/qso \
  -H "Content-Type: application/json" \
  -d '{
    "callsign": "W1AW",
    "frequency": 14.074,
    "band": "20m",
    "mode": "FT8",
    "rst_sent": "-10",
    "rst_received": "-12",
    "qso_date": "2023-11-20",
    "time_on": "2023-11-20T14:30:00Z",
    "grid_square": "FN31",
    "tx_power": 100
  }'
```

### List QSOs
```bash
curl http://localhost:8080/api/v1/qso
```

### Filter QSOs
```bash
curl "http://localhost:8080/api/v1/qso?band=20m&mode=FT8"
```

### Get Stats
```bash
curl http://localhost:8080/api/v1/qso/stats
```

### Export ADIF
```bash
curl http://localhost:8080/api/v1/qso/export/adif > logbook.adi
```

### Import ADIF
```bash
curl -X POST http://localhost:8080/api/v1/qso/import/adif \
  -F "file=@test.adi"
```

---

## Database Validation

Connect to PostgreSQL:
```bash
docker exec -it hamradio-db psql -U postgres -d hamradio
```

Check QSOs:
```sql
SELECT callsign, band, mode, qso_date, time_on
FROM qsos
ORDER BY time_on DESC
LIMIT 10;
```

Count QSOs:
```sql
SELECT COUNT(*) FROM qsos;
```

Filter by band:
```sql
SELECT * FROM qsos WHERE band = '20m';
```

---

## Known Issues / Limitations

1. **No Authentication Yet:** All QSOs use test-user-id
2. **No LoTW Sync:** LoTW status indicators are placeholders
3. **No QSO Limits:** Tier enforcement not implemented
4. **Time Zones:** All times are UTC (as per amateur radio standard)

---

## Success Criteria

Phase 2 is successful if:

- ✅ Can create, read, update, delete QSOs
- ✅ ADIF import works with valid files
- ✅ ADIF export produces valid files
- ✅ Filtering works for callsign, band, mode, dates
- ✅ UI is responsive and provides feedback
- ✅ No console errors in browser
- ✅ API returns proper HTTP status codes
- ✅ Database schema handles data correctly

---

## Troubleshooting

### "Failed to load QSOs"
- Check backend is running on port 8080
- Check database is running
- Check browser console for errors
- Verify database connection in backend logs

### "Import failed"
- Verify ADIF file format is correct
- Check for required fields (CALL, QSO_DATE, TIME_ON)
- Check backend logs for parsing errors

### QSOs not appearing
- Check database: `SELECT * FROM qsos;`
- Verify user_id matches (currently `test-user-id`)
- Check filters are cleared

### Backend won't start
- Check database is running
- Verify DATABASE_URL in .env or environment
- Run `go mod download` in backend directory

---

*Last Updated: 2025-11-23*
*Phase: 2 (Cloud Logbook Core)*
