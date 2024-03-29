package goya

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type testStruct struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func TestGet(t *testing.T) {
	// Pass the map as params
	// The map can be map[any]any, but the first 'any' will be changed to a string using fmt.Sprintf.
	resp := Get[*BasicGetResponse](getURL, map[string]string{"temp": "2"})
	// The response will be the zero value of the type you specified if some errors occur
	if resp == nil {
		t.Fatal("Get got nil")
	}
	if resp.URL != stringPlus(getURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, stringPlus(getURL, "?temp=2"))
	}

	req := testStruct{"Hello", 3306}
	// Pass the struct as params
	resp2 := Get[BasicGetResponse](getURL, req)
	// The order of the queries will be random
	if resp2.URL != stringPlus(getURL, "?name=Hello&id=3306") && resp2.URL != stringPlus(getURL, "?id=3306&name=Hello") {
		t.Errorf("resp.URL got %v but want %v", resp2.URL, stringPlus(getURL, "?name=Hello&id=3306"))
	}
	// You can also pass the pointer of the struct
	resp3 := Get[BasicGetResponse](getURL, &req)
	if resp3.URL != stringPlus(getURL, "?name=Hello&id=3306") && resp3.URL != stringPlus(getURL, "?id=3306&name=Hello") {
		t.Errorf("resp.URL got %v but want %v", resp3.URL, stringPlus(getURL, "?name=Hello&id=3306"))
	}
}

func TestGetOpts(t *testing.T) {
	// GetOpts can support more features
	resp := GetOpts[BasicGetResponse](getURL, NewOption(WithParams(map[string]string{"temp": "2"})))
	if resp.URL != stringPlus(getURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, stringPlus(getURL, "?temp=2"))
	}
}

func TestPost(t *testing.T) {
	// Pass the map as Json
	// The map can be map[any]any, but the first 'any' will be changed to a string using fmt.Sprintf.
	resp := Post[*BasicPostResponse](postURL, map[string]string{"temp": "2"})
	// The response will be the zero value of the type you specified if some errors occur
	if resp == nil {
		t.Fatal("Post got nil")
	}
	mp := map[string]string{}
	json.Unmarshal([]byte(resp.Data), &mp)
	if !reflect.DeepEqual(mp, map[string]string{"temp": "2"}) {
		t.Errorf("resp.Data got %v but want %v", mp, map[string]string{"temp": "2"})
	}

	req := testStruct{"Hello", 3306}
	// Pass the struct as Json
	// The struct can also be a pointer.
	resp2 := Post[BasicPostResponse](postURL, req)
	data := &testStruct{}
	json.Unmarshal([]byte(resp2.Data), data)

	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

func TestPostOpts(t *testing.T) {
	req := testStruct{"Hello", 3306}
	// PostOpts can support more features
	resp := PostOpts[BasicPostResponse](postURL, NewOption(WithParams(map[string]string{"temp": "2"}), WithJson(req)))
	if resp.URL != stringPlus(postURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, stringPlus(postURL, "?temp=2"))
	}
	data := &testStruct{}
	json.Unmarshal([]byte(resp.Data), data)
	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

func TestPut(t *testing.T) {
	// Pass the map as Json
	// The map can be map[any]any, but the first 'any' will be changed to a string using fmt.Sprintf.
	resp := Put[*BasicPostResponse](putURL, map[string]string{"temp": "2"})
	// The response will be the zero value of the type you specified if some errors occur
	if resp == nil {
		t.Fatal("Put got nil")
	}
	mp := map[string]string{}
	json.Unmarshal([]byte(resp.Data), &mp)
	if !reflect.DeepEqual(mp, map[string]string{"temp": "2"}) {
		t.Errorf("resp.Data got %v but want %v", mp, map[string]string{"temp": "2"})
	}

	req := testStruct{"Hello", 3306}
	// Pass the struct as Json
	// The struct can also be a pointer.
	resp2 := Put[BasicPostResponse](putURL, req)
	data := &testStruct{}
	json.Unmarshal([]byte(resp2.Data), data)

	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

func TestPutOpts(t *testing.T) {
	req := testStruct{"Hello", 3306}
	// PutOpts can support more features
	resp := PutOpts[BasicPostResponse](putURL, NewOption(WithParams(map[string]string{"temp": "2"}), WithJson(req)))
	if resp.URL != stringPlus(putURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, stringPlus(putURL, "?temp=2"))
	}
	data := &testStruct{}
	json.Unmarshal([]byte(resp.Data), data)
	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

func TestDel(t *testing.T) {
	// Pass the map as Json
	// The map can be map[any]any, but the first 'any' will be changed to a string using fmt.Sprintf.
	resp := Delete[*BasicPostResponse](delURL, map[string]string{"temp": "2"})
	// The response will be the zero value of the type you specified if some errors occur
	if resp == nil {
		t.Fatal("Del got nil")
	}
	mp := map[string]string{}
	json.Unmarshal([]byte(resp.Data), &mp)
	if !reflect.DeepEqual(mp, map[string]string{"temp": "2"}) {
		t.Errorf("resp.Data got %v but want %v", mp, map[string]string{"temp": "2"})
	}

	req := testStruct{"Hello", 3306}
	// Pass the struct as Json
	// The struct can also be a pointer.
	resp2 := Delete[BasicPostResponse](delURL, req)
	data := &testStruct{}
	json.Unmarshal([]byte(resp2.Data), data)

	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

func TestDeleteOpts(t *testing.T) {
	req := testStruct{"Hello", 3306}
	// DeleteOpts can support more features
	resp := DeleteOpts[BasicPostResponse](delURL, NewOption(WithParams(map[string]string{"temp": "2"}), WithJson(req)))
	if resp.URL != stringPlus(delURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, stringPlus(delURL, "?temp=2"))
	}
	data := &testStruct{}
	json.Unmarshal([]byte(resp.Data), data)
	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

type FormStruct struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	Numbers []string `json:"numbers"`
}

func TestFormData(t *testing.T) {
	// The Form only support string and []string
	req := FormStruct{"Hello", "3306", []string{"1", "2", "3"}}
	resp := PostOpts[BasicPostResponse](postURL, NewOption(WithForm(req)))
	data := &FormStruct{}
	bts, _ := json.Marshal(resp.Form)
	json.Unmarshal(bts, data)

	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Form got %v but want %v", data, req)
	}

	req2 := map[string]any{"name": "Hello", "numbers": []any{"1", "2", "3"}}
	resp2 := PostOpts[BasicPostResponse](postURL, NewOption(WithForm(req2)))
	data2 := map[string]any{}
	bts, _ = json.Marshal(resp2.Form)
	json.Unmarshal(bts, &data2)

	if !reflect.DeepEqual(data2, req2) {
		t.Errorf("resp.Form got %v but want %v", data2, req2)
	}
}

func TestForceHeaders(t *testing.T) {
	headers := map[string][]string{"Content-Type": {"123"}, "Test-Header": {"1", "2", "3"}}
	resp := GetOpts[BasicGetResponse](getURL, NewOption(WithForceHeaders(headers)))

	if resp.Headers.ContentType != "123" {
		t.Errorf("Content-Type got %v but want %v", resp.Headers.ContentType, "123")
	}
	if resp.Headers.TestHeader != "1,2,3" {
		t.Errorf("Content-Type got %v but want %v", resp.Headers.TestHeader, "1,2,3")
	}
}

func TestCookies(t *testing.T) {
	cookies := []*http.Cookie{{Name: "Test", Value: "123"}, {Name: "Test2", Value: "321"}, {Name: "Test", Value: "12345"}}
	resp := GetOpts[BasicGetResponse](getURL, NewOption(WithCookies(cookies)))

	if resp.Headers.Cookie != "Test=123; Test2=321; Test=12345" {
		t.Errorf("Cookies got %v but want %v", resp.Headers.Cookie, "Test=123; Test2=321; Test=12345")
	}
}

func TestTimeout(t *testing.T) {
	timeout := 1 * time.Millisecond
	resp := GetOpts[*BasicGetResponse](getURL, NewOption(WithTimeout(timeout)))

	if resp != nil {
		t.Error("resp should be nil")
	}

	timeout = 1 * time.Minute
	resp = GetOpts[*BasicGetResponse](getURL, NewOption(WithTimeout(timeout)))

	if resp == nil {
		t.Fatal("resp got nil")
	}
	if resp.URL != getURL {
		t.Errorf("url got %v but want %v", resp.URL, getURL)
	}
}
