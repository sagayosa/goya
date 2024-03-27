package goya

type Option struct {
	// Json can be used then you want to send JSON within the request body
	// Json should be map or struct
	// Json will convert into the body
	Json any

	// Params must be map or struct
	// Params will convert into the URL as the query argument
	Params any

	Headers map[string]string
}

func NewOption(opts ...OptionFunc) *Option {
	opt := &Option{
		Headers: map[string]string{},
	}
	for _, f := range opts {
		f(opt)
	}
	return opt
}
