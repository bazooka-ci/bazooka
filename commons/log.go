package bazooka

import (
	"regexp"
	"strings"
	"time"
)

var (
	regLogLevel = regexp.MustCompile(`^\s*\[(\S+)\].*$`)      // Eg. [INFO] My message
	regMeta     = regexp.MustCompile(`^\s*\<(\S+):(.*)>\s*$`) // Eg. <CMD:go test -v ./...>
)

func ConstructLog(message string, template LogEntry) LogEntry {
	message = strings.TrimSpace(message)

	switch {
	case regLogLevel.MatchString(message):
		submatchs := regLogLevel.FindStringSubmatch(message)
		logLevel := submatchs[len(submatchs)-1]
		template.Level = logLevel
		template.Message = strings.TrimSpace(message[len(logLevel)+2:])
	case regMeta.MatchString(message):
		submatchs := regMeta.FindStringSubmatch(message)
		instructionType := submatchs[1]
		instructionValue := submatchs[2]
		switch instructionType {
		case "CMD":
			template.Command = instructionValue
			template.Message = ""
		case "PHASE":
			template.Phase = instructionValue
			template.Message = ""
		default:
			template.Message = message
		}
	default:
		template.Message = message
	}

	if template.Time.IsZero() {
		template.Time = time.Now()
	}
	return template
}
