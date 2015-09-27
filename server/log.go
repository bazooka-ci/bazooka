package main

import (
	"bufio"
	"net"

	"io"

	log "github.com/Sirupsen/logrus"
	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/commons/syslogparser"
)

func (c *context) startLogServer(iface string) {
	server, err := net.Listen("tcp", iface)
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Errorf("Cannot accept log conn: %v", err)
		}

		go c.handleLogConn(conn)
	}
}

func (c *context) handleLogConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Errorf("Error reading log: %v\n", err)
			}
			return
		}

		p, err := syslogparser.Parse(line)
		if err != nil {
			return
		}

		template := lib.LogEntry{
			ProjectID: p.Meta["project"],
			JobID:     p.Meta["job"],
			VariantID: p.Meta["variant"],
			Image:     p.Meta["image"],
			Time:      p.Timestamp,
		}
		entry := lib.ConstructLog(p.Content, template)
		if err := c.connector.AddLog(&entry); err != nil {
			log.Errorf("Error adding log entry %v: %v", entry, err)
		}
	}
}
