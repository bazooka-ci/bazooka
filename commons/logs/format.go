package logs

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type BzkFormatter struct {
}

func (f *BzkFormatter) Format(entry *log.Entry) ([]byte, error) {

	b := &bytes.Buffer{}

	fmt.Fprintf(b, "%s [%s] %s ", entry.Time.Format(time.RFC3339), strings.ToUpper(entry.Level.String()), entry.Message)

	for k, v := range entry.Data {
		fmt.Fprintf(b, "%v=%s ", k, v)
	}
	b.WriteByte('\n')

	return b.Bytes(), nil
}
