package sip

import (
	"fmt"
	pjsua2 "sova-caller-backend/pjsua2"
	"strings"
)

type LogWriter struct {
	name string
}

func (l *LogWriter) Write(entry pjsua2.LogEntry) {
	msgRaw := entry.GetMsg()
	msg := strings.Replace(msgRaw, "\r", "", -1)

	if msg[len(msg)-1] == '\n' {
		msg = msg[37 : len(msg)-1]
	}

	fmt.Printf("[ SIP ] %s\n", msg)
}
