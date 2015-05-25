package syslog

import (
	"fmt"
	"log/syslog"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Priority  syslog.Priority
	Facility  syslog.Priority
	Severity  syslog.Priority
	Timestamp time.Time
	Host      string
	Meta      map[string]string
	Pid       int
	Content   string
}

const (
	severityMask = 0x07
	facilityMask = 0xf8
)

func Parse(line []byte) (*Message, error) {
	return (&parser{string(line), 0}).parse()
}

type parser struct {
	line string
	pos  int
}

func (p *parser) parse() (msg *Message, err error) {
	defer func() {
		if e := recover(); e != nil {
			msg = nil
			err = fmt.Errorf("%s\n%s^ %v\n", p.line, strings.Repeat(" ", p.pos), e)
		}
	}()
	p.line = strings.TrimSuffix(p.line, "\n")
	err = nil
	msg = &Message{}
	p.expect("<")
	{
		rawPri := p.until(">", "priority")
		pri, err := strconv.Atoi(rawPri)
		if err != nil {
			return nil, fmt.Errorf("Invalid priority %s", rawPri)
		}

		msg.Priority = syslog.Priority(pri)
		msg.Facility = syslog.Priority(pri & facilityMask)
		msg.Severity = syslog.Priority(pri & severityMask)
	}

	{
		ts := p.until(" ", "timestamp")
		parsedTimestamp, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return nil, err
		}
		msg.Timestamp = parsedTimestamp
	}

	{
		msg.Host = p.until(" ", "host")
	}

	{
		msg.Meta = map[string]string{}

		tag := p.until("[", "tag")
		elfAndMeta := strings.SplitN(tag, "/", 2)
		meta := elfAndMeta[1]
		kvs := strings.Split(meta, ";")
		fmt.Printf("TAG: %s\n", tag)
		for _, kv := range kvs {
			a := strings.Split(kv, "=")
			msg.Meta[a[0]] = a[1]
		}
	}

	{
		rawPid := p.until("]", "pid")
		pid, err := strconv.Atoi(rawPid)
		if err != nil {
			return nil, fmt.Errorf("Invalid pid: %s", rawPid)
		}
		msg.Pid = pid
	}

	p.expect(": ")
	msg.Content = p.line[p.pos:]
	return
}

func (p *parser) until(end, name string) string {
	pos0 := p.pos
	for !p.eof() {
		if !strings.HasPrefix(p.line[p.pos:], end) {
			p.pos++
			continue
		}
		break
	}
	if pos0 == p.pos {
		panic(fmt.Sprintf("Missing %s", name))
	}

	res := p.line[pos0:p.pos]
	if !p.eof() {
		p.pos++
	}
	return res
}

func (p *parser) eof() bool {
	return p.pos >= len(p.line)
}

func (p *parser) found(s string) bool {
	if strings.HasPrefix(p.line[p.pos:], s) {
		p.pos += len(s)
		return true
	}
	return false
}

func (p *parser) expect(s string) {
	if p.found(s) {
		return
	}
	panic(fmt.Sprintf("Was expecting '%s' but got '%s", s, p.line[p.pos:]))
}
