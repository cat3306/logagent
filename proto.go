package main

import (
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
	LogLevel   uint16
	ServerName string
	Payload    []byte
}

func Decode(buf []byte) *Context {
	hLen := packetEndian.Uint32(buf[:headerLen])
	header := buf[headerLen : headerLen+hLen]
	payload := buf[hLen+headerLen:]
	logLevel := packetEndian.Uint16(header[:logLevelLen])
	serverName := header[logLevelLen:]
	copyPayload := make([]byte, len(payload))
	copy(copyPayload, payload)
	return &Context{
		LogLevel:   logLevel,
		ServerName: string(serverName), //copy
		Payload:    copyPayload,        //copy
	}
}
