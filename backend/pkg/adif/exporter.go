package adif

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/nitefawkes/ham-radio-cloud/internal/models"
)

type Exporter struct {
	writer io.Writer
}

func NewExporter(writer io.Writer) *Exporter {
	return &Exporter{writer: writer}
}

func (e *Exporter) Export(qsos []models.QSO, programName, programVersion string) error {
	// Write ADIF header
	header := fmt.Sprintf("ADIF Export from %s v%s\n", programName, programVersion)
	header += fmt.Sprintf("Generated: %s\n", time.Now().Format(time.RFC3339))
	header += "<ADIF_VER:5>3.1.0\n"
	header += fmt.Sprintf("<PROGRAMID:%d>%s\n", len(programName), programName)
	header += fmt.Sprintf("<PROGRAMVERSION:%d>%s\n", len(programVersion), programVersion)
	header += "<EOH>\n\n"

	if _, err := e.writer.Write([]byte(header)); err != nil {
		return fmt.Errorf("failed to write ADIF header: %w", err)
	}

	// Write each QSO record
	for _, qso := range qsos {
		record := e.formatQSO(qso)
		if _, err := e.writer.Write([]byte(record)); err != nil {
			return fmt.Errorf("failed to write QSO record: %w", err)
		}
	}

	return nil
}

func (e *Exporter) formatQSO(qso models.QSO) string {
	var fields []string

	// Required fields
	fields = append(fields, e.formatField("CALL", qso.Callsign))

	// Date and Time
	qsoDate := qso.QSODate.Format("20060102")
	timeOn := qso.TimeOn.Format("150405")
	fields = append(fields, e.formatField("QSO_DATE", qsoDate))
	fields = append(fields, e.formatField("TIME_ON", timeOn))

	if qso.TimeOff != nil {
		timeOff := qso.TimeOff.Format("150405")
		fields = append(fields, e.formatField("TIME_OFF", timeOff))
	}

	// Band and Frequency
	if qso.Band != "" {
		fields = append(fields, e.formatField("BAND", qso.Band))
	}
	if qso.Frequency > 0 {
		freq := fmt.Sprintf("%.6f", qso.Frequency)
		fields = append(fields, e.formatField("FREQ", freq))
	}

	// Mode
	if qso.Mode != "" {
		fields = append(fields, e.formatField("MODE", qso.Mode))
	}

	// RST
	if qso.RST_Sent != "" {
		fields = append(fields, e.formatField("RST_SENT", qso.RST_Sent))
	}
	if qso.RST_Received != "" {
		fields = append(fields, e.formatField("RST_RCVD", qso.RST_Received))
	}

	// Location
	if qso.GridSquare != "" {
		fields = append(fields, e.formatField("GRIDSQUARE", qso.GridSquare))
	}
	if qso.Country != "" {
		fields = append(fields, e.formatField("COUNTRY", qso.Country))
	}
	if qso.State != "" {
		fields = append(fields, e.formatField("STATE", qso.State))
	}
	if qso.County != "" {
		fields = append(fields, e.formatField("CNTY", qso.County))
	}

	// Additional fields
	if qso.Comment != "" {
		fields = append(fields, e.formatField("COMMENT", qso.Comment))
	}
	if qso.TXPower > 0 {
		power := fmt.Sprintf("%.0f", qso.TXPower)
		fields = append(fields, e.formatField("TX_PWR", power))
	}

	// Contest and propagation
	if qso.ContestID != "" {
		fields = append(fields, e.formatField("CONTEST_ID", qso.ContestID))
	}
	if qso.PropagationMode != "" {
		fields = append(fields, e.formatField("PROP_MODE", qso.PropagationMode))
	}
	if qso.SatelliteName != "" {
		fields = append(fields, e.formatField("SAT_NAME", qso.SatelliteName))
	}

	// LoTW status
	if qso.LoTWConfirmed {
		fields = append(fields, e.formatField("LOTW_QSLRDATE", qso.UpdatedAt.Format("20060102")))
	}

	record := strings.Join(fields, " ")
	record += " <EOR>\n"

	return record
}

func (e *Exporter) formatField(name, value string) string {
	if value == "" {
		return ""
	}
	return fmt.Sprintf("<%s:%d>%s", strings.ToUpper(name), len(value), value)
}
