package syslogparser

import (
	"github.com/stretchr/testify/require"

	"log/syslog"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	msg, err := Parse([]byte("<27>2015-06-07T16:12:49Z jessie-amd64 docker-1.7.0-dev/image=bazooka/scm-git;project=pid;job=jid[4432]: Warning: Permanently added 'bitbucket.org,131.103.20.16"))
	require.NoError(t, err, "Should parse")

	require.Equal(t, syslog.Priority(27), msg.Priority)
	require.Equal(t, syslog.Priority(24), msg.Facility)
	require.Equal(t, syslog.Priority(3), msg.Severity)
	ts, _ := time.Parse("2006-01-02T15:04:05", "2015-06-07T16:12:49")
	require.Equal(t, ts, msg.Timestamp)
	require.Equal(t, "jessie-amd64", msg.Host)
	require.Equal(t, map[string]string{
		"project": "pid",
		"job":     "jid",
		"image":   "bazooka/scm-git",
	}, msg.Meta)
	require.Equal(t, 4432, msg.Pid)
	require.Equal(t, "Warning: Permanently added 'bitbucket.org,131.103.20.16", msg.Content)
}
