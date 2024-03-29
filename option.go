package goya

type Option struct {
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
