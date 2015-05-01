package bazooka

import (
	"fmt"
	"net"
	"time"
)

func WaitForTcpConnection(host, port string, retryEvery, timeout time.Duration) error {
	url := fmt.Sprintf("%s:%s", host, port)
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
			return fmt.Errorf("Coudln't establish a connection to %s:%s after %v", host, port, timeout)
		}
	}
}
