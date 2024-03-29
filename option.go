package goya

type Option struct {
	before []BeforeBuildFunc
	after  []AfterBuildFunc
	client []ClientBuildFunc
}

func NewOption(opts ...OptionFunc) *Option {
	opt := &Option{
		before: []BeforeBuildFunc{},
		after:  []AfterBuildFunc{},
		client: []ClientBuildFunc{},
	}
	for _, f := range opts {
		b, a, c := f()
		opt.before = append(opt.before, b)
		opt.after = append(opt.after, a)
		opt.client = append(opt.client, c)
	}
	return opt
}
