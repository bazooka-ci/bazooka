package syslog

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestParse(t *testing.T) {
	msg, err := Parse([]byte("<27>2015-06-07T16:12:49Z jessie-amd64 docker-1.7.0-dev/image=bazooka/scm-git;project=05d8477d809d1d0ec642676d4e9a7575;job=7c811d7e89a83dbb3576a2cda8b3cbe2[4432]: Warning: Permanently added 'bitbucket.org,131.103.20.16"))
	require.NoError(t, err, "Should parse")
	require.Equal(t, 27, msg.Priority)
}
