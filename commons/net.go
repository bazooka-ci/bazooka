package bazooka

import (
	"fmt"
	"net"
	"time"
)

func WaitForTcpConnection(url string, retryEvery, timeout time.Duration) error {
	giveUp := time.After(timeout)
	for {
		select {
		case <-time.After(retryEvery):
			conn, err := net.DialTimeout("tcp", url, 100*time.Millisecond)
			if err == nil {
				conn.Close()
				return nil
			}
		case <-giveUp:
			return fmt.Errorf("Coudln't establish a connection to %s after %v", url, timeout)
		}
	}
}
