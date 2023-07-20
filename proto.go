package main

import (
	"cloud/logagent/internal"
	"encoding/binary"
)

//header

const (
	headerLen   = uint32(4)
	logLevelLen = uint32(2)
)

var (
	packetEndian = binary.LittleEndian
)

type Context struct {
	LogLevel     uint16
	ServerName   string
	Payload      []byte
	PayloadLen   int
}

func Decode(buf []byte) *Context {
	hLen := packetEndian.Uint32(buf[:headerLen])
	header := buf[headerLen : headerLen+hLen]
	payload := buf[hLen+headerLen:]
	logLevel := packetEndian.Uint16(header[:logLevelLen])
	serverName := header[logLevelLen:]

	copyPayload := internal.BUFFERPOOL.Get(uint32(len(payload)))
	copy(*copyPayload, payload)
	return &Context{
		LogLevel:   logLevel,
		ServerName: string(serverName), //copy
		Payload:    *copyPayload,       //copy
		PayloadLen: len(payload),
	}
}
