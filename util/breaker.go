package util

import (
	"io"
	"regexp"
)

// BytesLineBreaker ...
type BytesLineBreaker struct {
	reader         io.ReadCloser
	tempBuf        []byte
	lineBuf        chan []byte
	errChan        chan error
	backgroundRead bool
	re             *regexp.Regexp
}

// NewBytesLineBreaker ...
func NewBytesLineBreaker(reader io.ReadCloser) *BytesLineBreaker {
	return &BytesLineBreaker{
		reader:  reader,
		lineBuf: make(chan []byte, 100),
		errChan: make(chan error, 1),
		re:      regexp.MustCompile(`.*(\n|$)`),
	}
}

// LinesChan ...
func (b *BytesLineBreaker) LinesChan() <-chan []byte {
	if !b.backgroundRead {
		b.backgroundRead = true
		go b.start()
	}

	return b.lineBuf
}

// ErrorChan ...
func (b *BytesLineBreaker) ErrorChan() <-chan error {
	return b.errChan
}

func (b *BytesLineBreaker) start() {
	var buf []byte
	var err error
	for buf, err = b.fetchMore(); err == nil || err == io.EOF; buf, err = b.fetchMore() {
		buf = append(b.tempBuf, buf...)

		lines, leftover := b.breakLines(buf)
		if err == io.EOF {
			lines = append(lines, leftover)
		} else {
			b.tempBuf = leftover
		}

		for _, line := range lines {
			b.lineBuf <- line
		}

		if err == io.EOF {
			close(b.lineBuf)
			return
		}
	}

	if err != nil {
		close(b.lineBuf)
		b.errChan <- err
	}
	return
}

func (b *BytesLineBreaker) fetchMore() ([]byte, error) {
	buf := make([]byte, 4096) // 4096 bytes per fetch
	readN, err := b.reader.Read(buf)

	return buf[:readN], err
}

func (b *BytesLineBreaker) breakLines(buf []byte) ([][]byte, []byte) {
	lines := b.re.FindAll(buf, -1)

	return lines[:len(lines)-1], lines[len(lines)-1]
}

// Close ...
func (b *BytesLineBreaker) Close() error {
	b.reader.Close()
	return nil
}
