package table

import (
	"io"
	"os"
)

type Formatter func(string) string

type Options struct {
	out             io.Writer
	headerFormatter Formatter
}

type Option func(options *Options)

func DefaultOptions() *Options {
	return &Options{
		out:             os.Stdout,
		headerFormatter: nil,
	}
}

func WithOut(out io.Writer) Option {
	return func(options *Options) { options.out = out }
}

func WithHeaderFormatter(formatter Formatter) Option {
	return func(options *Options) { options.headerFormatter = formatter }
}
