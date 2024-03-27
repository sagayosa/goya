package goya

import (
	"testing"
)

func TestGet(t *testing.T) {
	resp := Get[*BasicGetResponse](testURL, nil)
	if resp == nil {
		t.Fatal("Get got nil")
	}
	if resp.URL != testURL {
		t.Errorf("resp.URL got %v but want %v", resp.URL, testURL)
	}
}
