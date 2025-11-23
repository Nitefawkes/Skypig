package adif

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Nitefawkes/Skypig/backend/internal/models"
)

// Parser handles ADIF format parsing
type Parser struct {
	strict bool
}

// NewParser creates a new ADIF parser
func NewParser(strict bool) *Parser {
	return &Parser{strict: strict}
}

// Parse parses ADIF content and returns QSOs
func (p *Parser) Parse(content string) ([]*models.QSO, error) {
	// Split into records (separated by <eor>)
	records := strings.Split(content, "<eor>")

	var qsos []*models.QSO
	var errors []string

	for i, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		qso, err := p.parseRecord(record)
		if err != nil {
			errMsg := fmt.Sprintf("Record %d: %v", i+1, err)
			if p.strict {
				return nil, fmt.Errorf(errMsg)
			}
			errors = append(errors, errMsg)
			continue
		}

		if qso != nil {
			qsos = append(qsos, qso)
		}
	}

	if len(errors) > 0 && !p.strict {
		// Log warnings but don't fail
		fmt.Printf("ADIF parse warnings: %s\n", strings.Join(errors, "; "))
	}

	return qsos, nil
}

// parseRecord parses a single ADIF record
func (p *Parser) parseRecord(record string) (*models.QSO, error) {
	fields := p.extractFields(record)

	// Callsign is required
	callsign, ok := fields["call"]
	if !ok {
		return nil, fmt.Errorf("missing required field: CALL")
	}

	qso := &models.QSO{
		Callsign: strings.ToUpper(strings.TrimSpace(callsign)),
	}

	// Parse date and time
	if err := p.parseDateTime(fields, qso); err != nil {
		return nil, err
	}

	// Parse all other fields
	p.parseOptionalFields(fields, qso)

	return qso, nil
}

// extractFields extracts all ADIF fields from a record
func (p *Parser) extractFields(record string) map[string]string {
	fields := make(map[string]string)

	// ADIF field format: <FIELD_NAME:LENGTH>VALUE or <FIELD_NAME:LENGTH:TYPE>VALUE
	re := regexp.MustCompile(`<([^:>]+):(\d+)(?::([^>]+))?>([^<]*)`)
	matches := re.FindAllStringSubmatch(record, -1)

	for _, match := range matches {
		if len(match) >= 5 {
			fieldName := strings.ToLower(match[1])
			length, _ := strconv.Atoi(match[2])
			value := match[4]

			// Trim value to specified length
			if len(value) > length {
				value = value[:length]
			}

			fields[fieldName] = strings.TrimSpace(value)
		}
	}

	return fields
}

// parseDateTime parses QSO date and time fields
func (p *Parser) parseDateTime(fields map[string]string, qso *models.QSO) error {
	// QSO_DATE format: YYYYMMDD
	qsoDate, hasDate := fields["qso_date"]
	if !hasDate {
		return fmt.Errorf("missing required field: QSO_DATE")
	}

	// TIME_ON format: HHMM or HHMMSS
	timeOn, hasTime := fields["time_on"]
	if !hasTime {
		return fmt.Errorf("missing required field: TIME_ON")
	}

	// Parse date
	t, err := parseADIFDate(qsoDate, timeOn)
	if err != nil {
		return fmt.Errorf("invalid date/time: %w", err)
	}
	qso.TimeOn = t
	qso.QSODate = t

	// Parse TIME_OFF if present
	if timeOff, ok := fields["time_off"]; ok {
		qsoDateOff := qsoDate
		if dateOff, ok := fields["qso_date_off"]; ok {
			qsoDateOff = dateOff
		}

		tOff, err := parseADIFDate(qsoDateOff, timeOff)
		if err == nil {
			qso.TimeOff = tOff
		}
	}

	return nil
}

// parseADIFDate parses ADIF date and time
func parseADIFDate(date, timeStr string) (time.Time, error) {
	// Date format: YYYYMMDD
	// Time format: HHMM or HHMMSS

	if len(date) != 8 {
		return time.Time{}, fmt.Errorf("invalid date format (expected YYYYMMDD): %s", date)
	}

	// Pad time to 6 digits if needed
	if len(timeStr) == 4 {
		timeStr += "00"
	}
	if len(timeStr) != 6 {
		return time.Time{}, fmt.Errorf("invalid time format (expected HHMM or HHMMSS): %s", timeStr)
	}

	dateTime := date + timeStr
	t, err := time.Parse("20060102150405", dateTime)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// parseOptionalFields parses all optional ADIF fields
func (p *Parser) parseOptionalFields(fields map[string]string, qso *models.QSO) {
	// String fields
	if v, ok := fields["operator"]; ok {
		qso.OperatorCall = v
	}
	if v, ok := fields["station_callsign"]; ok {
		qso.StationCallsign = v
	}
	if v, ok := fields["band"]; ok {
		qso.Band = strings.ToLower(v)
	}
	if v, ok := fields["band_rx"]; ok {
		qso.BandRX = strings.ToLower(v)
	}
	if v, ok := fields["mode"]; ok {
		qso.Mode = strings.ToUpper(v)
	}
	if v, ok := fields["submode"]; ok {
		qso.Submode = strings.ToUpper(v)
	}
	if v, ok := fields["rst_sent"]; ok {
		qso.RSTSent = v
	}
	if v, ok := fields["rst_rcvd"]; ok {
		qso.RSTRcvd = v
	}
	if v, ok := fields["name"]; ok {
		qso.Name = v
	}
	if v, ok := fields["qth"]; ok {
		qso.QTH = v
	}
	if v, ok := fields["gridsquare"]; ok {
		qso.GridSquare = strings.ToUpper(v)
	}
	if v, ok := fields["country"]; ok {
		qso.Country = v
	}
	if v, ok := fields["state"]; ok {
		qso.State = strings.ToUpper(v)
	}
	if v, ok := fields["county"]; ok {
		qso.County = v
	}
	if v, ok := fields["comment"]; ok {
		qso.Comment = v
	}
	if v, ok := fields["notes"]; ok {
		qso.Notes = v
	}
	if v, ok := fields["prop_mode"]; ok {
		qso.PropagationMode = v
	}
	if v, ok := fields["sat_name"]; ok {
		qso.SatName = v
	}
	if v, ok := fields["sat_mode"]; ok {
		qso.SatMode = v
	}
	if v, ok := fields["contest_id"]; ok {
		qso.Contest = v
	}

	// Numeric fields
	if v, ok := fields["freq"]; ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			qso.Freq = f
		}
	}
	if v, ok := fields["freq_rx"]; ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			qso.FreqRX = f
		}
	}
	if v, ok := fields["tx_pwr"]; ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			qso.TXPower = f
		}
	}
	if v, ok := fields["rx_pwr"]; ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			qso.RXPower = f
		}
	}
	if v, ok := fields["dxcc"]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			qso.DXCC = i
		}
	}
	if v, ok := fields["stx"]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			qso.STX = i
		}
	}
	if v, ok := fields["srx"]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			qso.SRX = i
		}
	}

	// QSL fields
	if v, ok := fields["lotw_qsl_sent"]; ok {
		qso.LoTWQSLSent = strings.ToUpper(v)
	}
	if v, ok := fields["lotw_qsl_rcvd"]; ok {
		qso.LoTWQSLRcvd = strings.ToUpper(v)
	}
	if v, ok := fields["eqsl_qsl_sent"]; ok {
		qso.EQSLQSLSent = strings.ToUpper(v)
	}
	if v, ok := fields["eqsl_qsl_rcvd"]; ok {
		qso.EQSLQSLRcvd = strings.ToUpper(v)
	}
}

// ParseFile parses an ADIF file (convenience method)
func ParseFile(content string) ([]*models.QSO, error) {
	parser := NewParser(false) // Non-strict mode
	return parser.Parse(content)
}

// ParseStrict parses ADIF content in strict mode
func ParseStrict(content string) ([]*models.QSO, error) {
	parser := NewParser(true)
	return parser.Parse(content)
}
