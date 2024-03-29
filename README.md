## goya

goya is currently a personal HTTP request library for go

It currently only provides some basic functions for HTTP requests and needs further improvement



## Features

- Response can be deserialized through generics, without any error handling to obtain the final instance
- Currently, only providing sugar for GET POST PUT DELETE, but other verbs can also be specified
- Provides free request configuration and a easy way to write custom options



## Install

`go get -u github.com/sagayosa/goya`



## Basic Examples

**Sample GET**

``````go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sagayosa/goya"
)

type BasicGetResponse struct {
	Args any    `json:"args"`
	URL  string `json:"url"`
}

func main() {
	// sample get
	resp := goya.Get[*BasicGetResponse]("http://httpbin.org/get", nil)
    // if some errors occurred during the build or request
    // it will return a zero value of the type you specified
	if resp == nil {
		log.Fatalln("resp is nil")
	}

	// http://httpbin.org/get
	fmt.Println(resp.URL)

	// with params
	// params can be struct or map
	resp = goya.Get[*BasicGetResponse]("http://httpbin.org/get", map[string]string{"arg1": "1", "arg2": "2"})
	if resp == nil {
		log.Fatalln("resp is nil")
	}

	// map[arg1:1 arg2:2]
	fmt.Println(resp.Args)
	// http://httpbin.org/get?arg1=1&arg2=2
	fmt.Println(resp.URL)
}
``````

**Enhance Use**

``````go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sagayosa/goya"
)

type BasicPostResponse struct {
	Headers any    `json:"headers"`
	URL     string `json:"url"`
	Args    any    `json:"args"`
	Data    string `json:"data"`
}

func main() {
	// enhance use
	resp := goya.PostOpts[*BasicPostResponse]("http://httpbin.org/post", goya.NewOption(
		goya.WithParams(struct {
			Arg1 string `json:"arg1"`
			Arg2 string `json:"arg2"`
		}{"1", "2"}),
		goya.WithForceHeaders(http.Header{
			"Test-Header": []string{"Header-Test"},
		}),
		goya.WithJson(map[string]any{
			"arg3": struct {
				Arg4 string `json:"arg4"`
				Arg5 []int  `json:"arg5"`
			}{"4gra", []int{1, 2, 3}},
			"arg6": "6gra",
		}),
		goya.WithTimeout(5*time.Second),
		goya.WithCookies([]*http.Cookie{{Name: "Goya!!!", Value: "!!!ayoG"}}),
		// goya.WithForm(map[string]string{}),
	))
	if resp == nil {
		log.Fatalln("resp is nil")
	}
	// http://httpbin.org/post?arg1=1&arg2=2
	fmt.Println(resp.URL)

	// map[arg1:1 arg2:2]
	fmt.Println(resp.Args)

	// map[Accept-Encoding:gzip Content-Length:53 Content-Type:application/json Cookie:Goya!!!=!!!ayoG
	// Host:httpbin.org Test-Header:Header-Test User-Agent:Go-http-client/1.1
	// X-Amzn-Trace-Id:Root=1-66066c9e-446da9385bfb210076950b18]
	fmt.Println(resp.Headers)

	// {"arg3":{"arg4":"4gra","arg5":[1,2,3]},"arg6":"6gra"}
	fmt.Println(resp.Data)
}

``````

**Custom Option**

I only provided some basic options, but customizing options is also easy, such as WithCustomXML

``````go
package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sagayosa/goya"
)

type BasicPostResponse struct {
	Headers any    `json:"headers"`
	URL     string `json:"url"`
	Args    any    `json:"args"`
	Data    string `json:"data"`
}

func WithCustomXML(data any) goya.OptionFunc {
	return func() (goya.BeforeBuildFunc, goya.AfterBuildFunc, goya.ClientBuildFunc) {
		// BeforeBuildFunc allows you to change fields like url and body before the *http.Request is built
		return func(b *goya.RequestBuider) {
				bts, err := xml.Marshal(data)
				if err != nil {
					// it will add the err to b.errs
					// you can use b.Errors() retrieves all errors that occurred during the build
					b.ErrHappen(err)
				}
				// set the xml data to the body
				b.Body = bts
			},
			// AfterBuildFunc allows you to change fields like headers after the *http.Request is built
			func(req *http.Request) {
				req.Header.Set("Content-Type", "application/xml")
			},
			// ClientBuildFunc allows you to change the client that will ultimately be used for client.Do(request)
			func(client *http.Client) {
				// same as WithTimeout(5 * time.Second)
				client.Timeout = 5 * time.Second
			}
	}
}

func main() {
	resp := goya.PostOpts[*BasicPostResponse]("http://httpbin.org/post", goya.NewOption(WithCustomXML(
		struct {
			XMLName xml.Name `xml:"person"`
			Name    string   `xml:"name"`
			Age     int      `xml:"age"`
		}{xml.Name{}, "goya", 114514},
	)))
	if resp == nil {
		log.Fatalln("resp is nil")
	}

	// map[Accept-Encoding:gzip Content-Length:51 Content-Type:application/xml Host:httpbin.org
	// User-Agent:Go-http-client/1.1 X-Amzn-Trace-Id:Root=1-660670b9-5edbb9e566d1ace01214a675]
	fmt.Println(resp.Headers)

	// <person><name>goya</name><age>114514</age></person>
	fmt.Println(resp.Data)
}

``````



## License

This project is licensed under the [MIT license](https://github.com/sagayosa/goya/blob/master/LICENSE)

