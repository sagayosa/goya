package goya

import (
	"testing"
)

func TestGet(t *testing.T) {
	resp := Get[*BasicGetResponse](testURL, map[string]string{"temp": "2"})
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
	resp2 := Get[BasicGetResponse](testURL, req)
	if resp2.URL != StringPlus(testURL, "?name=Hello&id=3306") && resp2.URL != StringPlus(testURL, "?id=3306&name=Hello") {
		t.Errorf("resp.URL got %v but want %v", resp2.URL, StringPlus(testURL, "?name=Hello&id=3306"))
	}
}
