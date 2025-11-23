package adif

import (
	"fmt"
	"strings"
	"time"

	"github.com/Nitefawkes/Skypig/backend/internal/models"
)

// Generator handles ADIF format generation
type Generator struct {
	includeHeader bool
	appName       string
	appVersion    string
}

// NewGenerator creates a new ADIF generator
func NewGenerator() *Generator {
	return &Generator{
		includeHeader: true,
		appName:       "Ham-Radio Cloud",
		appVersion:    "1.0.0",
	}
}

// Generate converts QSOs to ADIF format
func (g *Generator) Generate(qsos []*models.QSO) string {
	var sb strings.Builder

	// Write header
	if g.includeHeader {
		g.writeHeader(&sb)
	}

	// Write each QSO record
	for _, qso := range qsos {
		g.writeRecord(&sb, qso)
	}

	return sb.String()
}

// writeHeader writes the ADIF header
func (g *Generator) writeHeader(sb *strings.Builder) {
	sb.WriteString("ADIF Export from ")
	sb.WriteString(g.appName)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("<adif_ver:5>3.1.4\n")
	sb.WriteString(fmt.Sprintf("<programid:%d>%s\n", len(g.appName), g.appName))
	sb.WriteString(fmt.Sprintf("<programversion:%d>%s\n", len(g.appVersion), g.appVersion))
	sb.WriteString("<eoh>\n\n")
}

// writeRecord writes a single QSO record
func (g *Generator) writeRecord(sb *strings.Builder, qso *models.QSO) {
	// Required fields
	g.writeField(sb, "CALL", qso.Callsign)
	g.writeField(sb, "QSO_DATE", formatADIFDate(qso.QSODate))
	g.writeField(sb, "TIME_ON", formatADIFTime(qso.TimeOn))

	// Optional time_off
	if !qso.TimeOff.IsZero() {
		g.writeField(sb, "TIME_OFF", formatADIFTime(qso.TimeOff))
		if !qso.TimeOff.Truncate(24*time.Hour).Equal(qso.TimeOn.Truncate(24 * time.Hour)) {
			g.writeField(sb, "QSO_DATE_OFF", formatADIFDate(qso.TimeOff))
		}
	}

	// Operator and station
	if qso.OperatorCall != "" {
		g.writeField(sb, "OPERATOR", qso.OperatorCall)
	}
	if qso.StationCallsign != "" {
		g.writeField(sb, "STATION_CALLSIGN", qso.StationCallsign)
	}

	// Frequency and band
	if qso.Freq != 0 {
		g.writeField(sb, "FREQ", fmt.Sprintf("%.6f", qso.Freq))
	}
	if qso.FreqRX != 0 {
		g.writeField(sb, "FREQ_RX", fmt.Sprintf("%.6f", qso.FreqRX))
	}
	if qso.Band != "" {
		g.writeField(sb, "BAND", qso.Band)
	}
	if qso.BandRX != "" {
		g.writeField(sb, "BAND_RX", qso.BandRX)
	}

	// Mode
	if qso.Mode != "" {
		g.writeField(sb, "MODE", qso.Mode)
	}
	if qso.Submode != "" {
		g.writeField(sb, "SUBMODE", qso.Submode)
	}

	// Signal reports
	if qso.RSTSent != "" {
		g.writeField(sb, "RST_SENT", qso.RSTSent)
	}
	if qso.RSTRcvd != "" {
		g.writeField(sb, "RST_RCVD", qso.RSTRcvd)
	}

	// Contact info
	if qso.Name != "" {
		g.writeField(sb, "NAME", qso.Name)
	}
	if qso.QTH != "" {
		g.writeField(sb, "QTH", qso.QTH)
	}
	if qso.GridSquare != "" {
		g.writeField(sb, "GRIDSQUARE", qso.GridSquare)
	}

	// Location
	if qso.Country != "" {
		g.writeField(sb, "COUNTRY", qso.Country)
	}
	if qso.DXCC != 0 {
		g.writeField(sb, "DXCC", fmt.Sprintf("%d", qso.DXCC))
	}
	if qso.State != "" {
		g.writeField(sb, "STATE", qso.State)
	}
	if qso.County != "" {
		g.writeField(sb, "COUNTY", qso.County)
	}

	// Power
	if qso.TXPower != 0 {
		g.writeField(sb, "TX_PWR", fmt.Sprintf("%.0f", qso.TXPower))
	}
	if qso.RXPower != 0 {
		g.writeField(sb, "RX_PWR", fmt.Sprintf("%.0f", qso.RXPower))
	}

	// Propagation
	if qso.PropagationMode != "" {
		g.writeField(sb, "PROP_MODE", qso.PropagationMode)
	}
	if qso.SatName != "" {
		g.writeField(sb, "SAT_NAME", qso.SatName)
	}
	if qso.SatMode != "" {
		g.writeField(sb, "SAT_MODE", qso.SatMode)
	}

	// Contest
	if qso.Contest != "" {
		g.writeField(sb, "CONTEST_ID", qso.Contest)
	}
	if qso.STX != 0 {
		g.writeField(sb, "STX", fmt.Sprintf("%d", qso.STX))
	}
	if qso.SRX != 0 {
		g.writeField(sb, "SRX", fmt.Sprintf("%d", qso.SRX))
	}

	// QSL
	if qso.LoTWQSLSent != "" {
		g.writeField(sb, "LOTW_QSL_SENT", qso.LoTWQSLSent)
	}
	if qso.LoTWQSLRcvd != "" {
		g.writeField(sb, "LOTW_QSL_RCVD", qso.LoTWQSLRcvd)
	}
	if qso.EQSLQSLSent != "" {
		g.writeField(sb, "EQSL_QSL_SENT", qso.EQSLQSLSent)
	}
	if qso.EQSLQSLRcvd != "" {
		g.writeField(sb, "EQSL_QSL_RCVD", qso.EQSLQSLRcvd)
	}

	// Comments and notes
	if qso.Comment != "" {
		g.writeField(sb, "COMMENT", qso.Comment)
	}
	if qso.Notes != "" {
		g.writeField(sb, "NOTES", qso.Notes)
	}

	// End of record
	sb.WriteString("<eor>\n\n")
}

// writeField writes an ADIF field
func (g *Generator) writeField(sb *strings.Builder, name, value string) {
	if value == "" {
		return
	}
	sb.WriteString(fmt.Sprintf("<%s:%d>%s ", name, len(value), value))
}

// formatADIFDate formats a time as YYYYMMDD
func formatADIFDate(t time.Time) string {
	return t.UTC().Format("20060102")
}

// formatADIFTime formats a time as HHMMSS
func formatADIFTime(t time.Time) string {
	return t.UTC().Format("150405")
}

// Export is a convenience function to generate ADIF from QSOs
func Export(qsos []*models.QSO) string {
	gen := NewGenerator()
	return gen.Generate(qsos)
}
