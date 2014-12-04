package logger

import (
	"io"
	"log"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Docker  *log.Logger
)

func Init(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer, dockerHandle io.Writer) {

	Trace = log.New(traceHandle,
		"[TRACE]: ",
		log.Ldate|log.Ltime)

	Info = log.New(infoHandle,
		"[INFO]: ",
		log.Ldate|log.Ltime)

	Warning = log.New(warningHandle,
		"[WARNING]: ",
		log.Ldate|log.Ltime)

	Error = log.New(errorHandle,
		"[ERROR]: ",
		log.Ldate|log.Ltime)

	Docker = log.New(dockerHandle,
		"[DOCKER]: ",
		log.Ldate|log.Ltime)
}
