package adif

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nitefawkes/ham-radio-cloud/internal/models"
)

// ADIF field pattern: <FIELD:LENGTH[:TYPE]>VALUE
var fieldPattern = regexp.MustCompile(`<([A-Z_]+):(\d+)(?::([A-Z]))?>([^<]*)`)

type Parser struct {
	reader io.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{reader: reader}
}

func (p *Parser) Parse() ([]models.QSO, error) {
	scanner := bufio.NewScanner(p.reader)
	scanner.Split(bufio.ScanLines)

	var content strings.Builder
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read ADIF file: %w", err)
	}

	text := content.String()

	// Find end of header
	eohIndex := strings.Index(strings.ToUpper(text), "<EOH>")
	if eohIndex != -1 {
		text = text[eohIndex+5:] // Skip <EOH>
	}

	// Split records by <EOR>
	records := strings.Split(strings.ToUpper(text), "<EOR>")

	var qsos []models.QSO
	for _, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		qso, err := p.parseRecord(record)
		if err != nil {
			// Log error but continue parsing other records
			continue
		}

		qsos = append(qsos, qso)
	}

	return qsos, nil
}

func (p *Parser) parseRecord(record string) (models.QSO, error) {
	qso := models.QSO{}
	fields := make(map[string]string)

	// Find all fields in the record
	matches := fieldPattern.FindAllStringSubmatch(record, -1)
	for _, match := range matches {
		if len(match) >= 5 {
			fieldName := strings.ToUpper(match[1])
			length, _ := strconv.Atoi(match[2])
			value := match[4]

			// Extract only the specified length of value
			if len(value) > length {
				value = value[:length]
			}

			fields[fieldName] = strings.TrimSpace(value)
		}
	}

	// Map ADIF fields to QSO model
	qso.Callsign = fields["CALL"]
	qso.Band = fields["BAND"]
	qso.Mode = fields["MODE"]
	qso.RST_Sent = fields["RST_SENT"]
	qso.RST_Received = fields["RST_RCVD"]
	qso.GridSquare = fields["GRIDSQUARE"]
	qso.Country = fields["COUNTRY"]
	qso.State = fields["STATE"]
	qso.County = fields["CNTY"]
	qso.Comment = fields["COMMENT"]

	// Parse frequency
	if freq := fields["FREQ"]; freq != "" {
		if f, err := strconv.ParseFloat(freq, 64); err == nil {
			qso.Frequency = f
		}
	}

	// Parse TX power
	if power := fields["TX_PWR"]; power != "" {
		if p, err := strconv.ParseFloat(power, 64); err == nil {
			qso.TXPower = p
		}
	}

	// Parse date and time
	qsoDate := fields["QSO_DATE"]
	timeOn := fields["TIME_ON"]
	timeOff := fields["TIME_OFF"]

	if qsoDate != "" && timeOn != "" {
		// ADIF date format: YYYYMMDD
		// ADIF time format: HHMMSS or HHMM
		dateStr := qsoDate
		timeStr := timeOn

		// Pad time to 6 digits if needed
		if len(timeStr) == 4 {
			timeStr += "00"
		}

		timestamp, err := time.Parse("20060102150405", dateStr+timeStr)
		if err == nil {
			qso.QSODate = timestamp
			qso.TimeOn = timestamp
		}
	}

	if qsoDate != "" && timeOff != "" {
		timeStr := timeOff
		if len(timeStr) == 4 {
			timeStr += "00"
		}
		timestamp, err := time.Parse("20060102150405", qsoDate+timeStr)
		if err == nil {
			qso.TimeOff = &timestamp
		}
	}

	// Validate required fields
	if qso.Callsign == "" {
		return qso, fmt.Errorf("missing required field: CALL")
	}

	return qso, nil
}
