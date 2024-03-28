package goya

type OptionFunc func(opt *Option)

// WithJson will add data to the request body
// and set the Content-Type to application/json
func WithJson(data any) OptionFunc {
	return func(opt *Option) {
		opt.Json = data
	}
}

// WithParams will add params to the request URL
// params must be a structure or a map
func WithParams(params any) OptionFunc {
	return func(opt *Option) {
		opt.Params = params
	}
}

// WithForm will add data to the request body
// and set the Content-Type to multipart/form-data
func WithForm(data any) OptionFunc {
	return func(opt *Option) {
		opt.FormData = data

	}
}
