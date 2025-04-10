package document

type SectionOptions struct {
	delimiter string
}

type SectionOption func(options *SectionOptions)

func DefaultSectionOptions() *SectionOptions {
	return &SectionOptions{
		delimiter: ".",
	}
}

func (options *SectionOptions) Apply(opts ...SectionOption) {
	for _, opt := range opts {
		opt(options)
	}
}

func WithDelimiter(delimiter string) SectionOption {
	return func(options *SectionOptions) {
		options.delimiter = delimiter
	}
}
