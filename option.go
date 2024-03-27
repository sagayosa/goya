package goya

type Option struct {
	// Data can be map or struct
	// Data will convert into the body
	Data any

	// Params can be map or struct
	// Params will convert into the URL as the query argument
	Params any
}

func NewOption(opts ...OptionFunc) *Option {
	opt := &Option{}
	for _, f := range opts {
		f(opt)
	}
	return opt
}
