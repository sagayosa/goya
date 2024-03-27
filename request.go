package goya

import (
	"encoding/json"
	"net/http"
)

type RequestBuider struct {
	Method string
	URL    string
	Opt    *Option

	buildErr  []error
	buildURL  string
	buildBody []byte
}

func (b *RequestBuider) Build() *http.Request {
	if b.Opt.Data != nil {
		b.buildData()
	}
	request, _ := http.NewRequest(b.Method, b.URL, nil)
	return request
}

func (b *RequestBuider) errHappen(err error) {
	b.buildErr = append(b.buildErr, err)
}

func (b *RequestBuider) buildData() {
	bts, err := json.Marshal(b.Opt.Data)
	if err != nil {
		b.errHappen(err)
	}
	b.buildBody = bts
}

func (b *RequestBuider) buildParams() {

}
