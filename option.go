package goya

type Option struct {
	before []BeforeBuildFunc
	after  []AfterBuildFunc
	client []ClientBuildFunc
	done   []ClientDoneFunc
}

func NewOption(opts ...OptionFunc) *Option {
	opt := &Option{
		before: []BeforeBuildFunc{},
		after:  []AfterBuildFunc{},
		client: []ClientBuildFunc{},
		done:   []ClientDoneFunc{},
	}
	for _, f := range opts {
		b, a, c, d := f()
		if b != nil {
			opt.before = append(opt.before, b)
		}
		if a != nil {
			opt.after = append(opt.after, a)
		}
		if c != nil {
			opt.client = append(opt.client, c)
		}
		if d != nil {
			opt.done = append(opt.done, d)
		}
	}
	return opt
}
