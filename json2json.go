package json2json

import (
	"bytes"
	"io"
)

type Json2Json struct {
	inputReader  io.Reader
	outputWriter io.Writer

	inputSampleReader io.Reader
	processReader     io.Reader

	fn func(r io.Reader, w io.Writer)
}

type Opt func(*Json2Json)

func New(r io.Reader, w io.Writer, opts ...Opt) *Json2Json {
	j := Json2Json{
		inputReader:  r,
		outputWriter: w,

		inputSampleReader: bytes.NewReader([]byte{}),
		processReader:     bytes.NewReader([]byte{}),
	}
	for _, opt := range opts {
		opt(&j)
	}
	return &j
}

func WithMiddlewareFn(fn func(r io.Reader, w io.Writer)) Opt {
	return func(j *Json2Json) {
		j.fn = fn
	}
}

func (j *Json2Json) ReadInput(b []byte) *Json2Json {
	return j
}

func (j *Json2Json) ReadInputFile(filepath string) *Json2Json {
	return j
}

func (j *Json2Json) ReadConfig(b []byte) *Json2Json {
	return j
}

func (j *Json2Json) ReadConfigFile(filepath string) *Json2Json {
	return j
}

func (j *Json2Json) WriteOutput() *Json2Json {
	if j.fn != nil {
		j.fn(j.inputReader, j.outputWriter)
	}
	return j
}
