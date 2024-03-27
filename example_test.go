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
		t.Errorf("resp.URL got %v but want %v", resp.URL, testURL)
	}
}
