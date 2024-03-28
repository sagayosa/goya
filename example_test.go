package goya

import (
	"encoding/json"
	"reflect"
	"testing"
)

type TestStruct struct {
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
	if resp.URL != StringPlus(getURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, StringPlus(getURL, "?temp=2"))
	}

	req := TestStruct{"Hello", 3306}
	// Pass the struct as params
	resp2 := Get[BasicGetResponse](getURL, req)
	// The order of the queries will be random
	if resp2.URL != StringPlus(getURL, "?name=Hello&id=3306") && resp2.URL != StringPlus(getURL, "?id=3306&name=Hello") {
		t.Errorf("resp.URL got %v but want %v", resp2.URL, StringPlus(getURL, "?name=Hello&id=3306"))
	}
	// You can also pass the pointer of the struct
	resp3 := Get[BasicGetResponse](getURL, &req)
	if resp3.URL != StringPlus(getURL, "?name=Hello&id=3306") && resp3.URL != StringPlus(getURL, "?id=3306&name=Hello") {
		t.Errorf("resp.URL got %v but want %v", resp3.URL, StringPlus(getURL, "?name=Hello&id=3306"))
	}
}

func TestGetOpts(t *testing.T) {
	// GetOpts can support more features
	resp := GetOpts[BasicGetResponse](getURL, NewOption(WithParams(map[string]string{"temp": "2"})))
	if resp.URL != StringPlus(getURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, StringPlus(getURL, "?temp=2"))
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

	req := TestStruct{"Hello", 3306}
	// Pass the struct as Json
	// The struct can also be a pointer.
	resp2 := Post[BasicPostResponse](postURL, req)
	data := &TestStruct{}
	json.Unmarshal([]byte(resp2.Data), data)

	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

func TestPostOpts(t *testing.T) {
	req := TestStruct{"Hello", 3306}
	// PostOpts can support more features
	resp := PostOpts[BasicPostResponse](postURL, NewOption(WithParams(map[string]string{"temp": "2"}), WithJson(req)))
	if resp.URL != StringPlus(postURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, StringPlus(postURL, "?temp=2"))
	}
	data := &TestStruct{}
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

	req := TestStruct{"Hello", 3306}
	// Pass the struct as Json
	// The struct can also be a pointer.
	resp2 := Put[BasicPostResponse](putURL, req)
	data := &TestStruct{}
	json.Unmarshal([]byte(resp2.Data), data)

	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

func TestPutOpts(t *testing.T) {
	req := TestStruct{"Hello", 3306}
	// PutOpts can support more features
	resp := PutOpts[BasicPostResponse](putURL, NewOption(WithParams(map[string]string{"temp": "2"}), WithJson(req)))
	if resp.URL != StringPlus(putURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, StringPlus(putURL, "?temp=2"))
	}
	data := &TestStruct{}
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

	req := TestStruct{"Hello", 3306}
	// Pass the struct as Json
	// The struct can also be a pointer.
	resp2 := Delete[BasicPostResponse](delURL, req)
	data := &TestStruct{}
	json.Unmarshal([]byte(resp2.Data), data)

	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}

func TestDeleteOpts(t *testing.T) {
	req := TestStruct{"Hello", 3306}
	// DeleteOpts can support more features
	resp := DeleteOpts[BasicPostResponse](delURL, NewOption(WithParams(map[string]string{"temp": "2"}), WithJson(req)))
	if resp.URL != StringPlus(delURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, StringPlus(delURL, "?temp=2"))
	}
	data := &TestStruct{}
	json.Unmarshal([]byte(resp.Data), data)
	if !reflect.DeepEqual(*data, req) {
		t.Errorf("resp.Data got %v but want %v", data, req)
	}
}
