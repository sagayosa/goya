package goya

type Option struct {
	// Json can be used then you want to send JSON within the request body
	// Json should be map or struct
	// Json will convert into the body
	// Json any

	// Params must be map or struct
	// Params will convert into the URL as the query argument
	// Params any

	// FormData must be map or struct
	// FormData will convert into the body
	// FormData any

	// Headers map[string]string
	before []BeforeBuildFunc
	after  []AfterBuildFunc
}

func NewOption(opts ...OptionFunc) *Option {
	opt := &Option{
		before: []BeforeBuildFunc{},
		after:  []AfterBuildFunc{},
	}
	for _, f := range opts {
		b, a := f(opt)
		opt.before = append(opt.before, b)
		opt.after = append(opt.after, a)
	}
	return opt
}
