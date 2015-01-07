package logs

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type BzkFormatter struct {
}

func (f *BzkFormatter) Format(entry *log.Entry) ([]byte, error) {

	b := &bytes.Buffer{}

	fmt.Fprintf(b, "[%s] %s ", strings.ToUpper(entry.Level.String()), entry.Message)

	for k, v := range entry.Data {
		fmt.Fprintf(b, "%v=%s ", k, v)
	}
	b.WriteByte('\n')

	return b.Bytes(), nil
}
