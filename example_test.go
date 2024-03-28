package goya

import (
	"testing"
)

func TestGet(t *testing.T) {
	// Pass the map as params
	// The map can be map[any]any, but the first 'any' will be changed to a string using fmt.Sprintf.
	resp := Get[*BasicGetResponse](testURL, map[string]string{"temp": "2"})
	// The response will be the zero value of the type you specified if some errors occur
	if resp == nil {
		t.Fatal("Get got nil")
	}
	if resp.URL != StringPlus(testURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, StringPlus(testURL, "?temp=2"))
	}

	req := struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}{"Hello", 3306}
	// Pass the struct as params
	resp2 := Get[BasicGetResponse](testURL, req)
	// The order of the queries will be random
	if resp2.URL != StringPlus(testURL, "?name=Hello&id=3306") && resp2.URL != StringPlus(testURL, "?id=3306&name=Hello") {
		t.Errorf("resp.URL got %v but want %v", resp2.URL, StringPlus(testURL, "?name=Hello&id=3306"))
	}
	// You can also pass the pointer of the struct
	resp3 := Get[BasicGetResponse](testURL, &req)
	if resp3.URL != StringPlus(testURL, "?name=Hello&id=3306") && resp3.URL != StringPlus(testURL, "?id=3306&name=Hello") {
		t.Errorf("resp.URL got %v but want %v", resp3.URL, StringPlus(testURL, "?name=Hello&id=3306"))
	}
}

func TestGetOpts(t *testing.T) {
	// GetOpts can support more features
	resp := GetOpts[BasicGetResponse](testURL, NewOption(WithParams(map[string]string{"temp": "2"})))
	if resp.URL != StringPlus(testURL, "?temp=2") {
		t.Errorf("resp.URL got %v but want %v", resp.URL, StringPlus(testURL, "?temp=2"))
	}
}
