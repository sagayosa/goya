package goya

type OptionFunc func(opt *Option)

// WithJson will add data to the request body
func WithJson(data any) OptionFunc {
	return func(opt *Option) {
		opt.Json = data
		if data != nil {
			opt.Headers[ContentType] = ContentTypeJSON
		}
	}
}

// WithParams will add params to the request URL
// params must be a structure or a map[string]any
func WithParams(params any) OptionFunc {
	return func(opt *Option) {
		opt.Params = params
	}
}
