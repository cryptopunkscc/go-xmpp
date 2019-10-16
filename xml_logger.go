package xmpp

import (
	"encoding/xml"
	"io"
	"log"
)

type ReadWriteLogger struct {
	target   io.ReadWriter
	readLog  io.Writer
	writeLog io.Writer
}

func (rwl *ReadWriteLogger) Write(p []byte) (n int, err error) {
	rwl.writeLog.Write(p)
	return rwl.target.Write(p)
}

func (rwl *ReadWriteLogger) Read(p []byte) (n int, err error) {
	n, err = rwl.target.Read(p)
	rwl.readLog.Write(p[:n])
	return n, err
}

type XMLLogger struct {
	piper   *io.PipeReader
	pipew   *io.PipeWriter
	decoder *xml.Decoder
	encoder *xml.Encoder
}

func NewXMLLogger(output io.Writer, prefix string) *XMLLogger {
	logger := &XMLLogger{}
	logger.piper, logger.pipew = io.Pipe()
	logger.decoder = xml.NewDecoder(logger.piper)
	logger.encoder = xml.NewEncoder(output)
	logger.encoder.Indent(prefix, "  ")
	go logger.rewriter()
	return logger
}

func (logger *XMLLogger) rewriter() {
	for {
		token, err := logger.decoder.Token()
		if err != nil {
			log.Printf("XMLLogger: error decoding token", err)
			return
		}
		if err := logger.encoder.EncodeToken(token); err != nil {
			log.Println("Error encoding xml log:", err)
		}
		logger.encoder.Flush()
	}
}

func (logger *XMLLogger) Write(p []byte) (n int, err error) {
	return logger.pipew.Write(p)
}
