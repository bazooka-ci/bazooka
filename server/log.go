package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	lib "github.com/bazooka-ci/bazooka/commons"
	"github.com/bazooka-ci/bazooka/commons/syslog"
)

func (c *context) startLogServer() {
	server, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalf("Cannot accept: %v", err)
		}

		go c.handleLogConn(conn)
	}
}

func (c *context) handleLogConn(conn net.Conn) {
	fmt.Printf("Accepted %v\n", conn.RemoteAddr())
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error reading: %v\n", err)
			return
		}

		p, err := syslog.Parse(line)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("Parsed: %#v\n", p)
		template := lib.LogEntry{
			ProjectID: p.Meta["project"],
			JobID:     p.Meta["job"],
			VariantID: p.Meta["variant"],
			Image:     p.Meta["image"],
			Time:      p.Timestamp,
		}
		entry := lib.ConstructLog(p.Content, template)
		c.connector.AddLog(&entry)
	}
}
