package client

import (
	"encoding/json"
	"io"
	"net/http"

	lib "github.com/bazooka-ci/bazooka/commons"
)

func streamLog(response http.Response) chan lib.LogEntry {
	sink := make(chan lib.LogEntry)

	go func() {
		defer response.Body.Close()
		decoder := json.NewDecoder(response.Body)

		var log lib.LogEntry

		for {
			err := decoder.Decode(&log)
			switch err {
			case nil:
				sink <- log
			case io.EOF:
				close(sink)
				return
			default:
				close(sink)
				return
			}
		}
	}()

	return sink
}
