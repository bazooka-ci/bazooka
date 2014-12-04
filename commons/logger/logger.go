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

	Trace = log.New(traceHandle, "[Trace] ", log.Ldate|log.Ltime)
	Info = log.New(infoHandle, "[Info] ", log.Ldate|log.Ltime)
	Warning = log.New(warningHandle, "[Warning] ", log.Ldate|log.Ltime)
	Error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime)
	Docker = log.New(dockerHandle, "[Docker] ", 0)
}
