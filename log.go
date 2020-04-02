package chromemobile

import (
	"log"
)

type LogWriter interface {
	WriteLog(s string)
}

func SetLogOutput(w LogWriter) {
	log.SetOutput(writerFunc(
		func(p []byte) (n int, err error) {
			w.WriteLog(string(p))
			return len(p), nil
		},
	))
}

type writerFunc func(p []byte) (n int, err error)

func (f writerFunc) Write(p []byte) (n int, err error) {
	return f(p)
}
