package goya

import (
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func TestConvertStructToMap(t *testing.T) {
	p := Person{Name: "John", Age: 30}
	mapResult := ConvertStructToMap(p)
	if mapResult == nil {
		t.Errorf("ConvertStructToMap returned nil for struct input")
	}
	if mapResult["Name"] != "John" || mapResult["Age"] != 30 {
		t.Errorf("ConvertStructToMap did not convert struct correctly")
	}

	ptrResult := ConvertStructToMap(&p)
	if ptrResult == nil {
		t.Errorf("ConvertStructToMap returned nil for pointer to struct input")
	}
	if ptrResult["Name"] != "John" || ptrResult["Age"] != 30 {
		t.Errorf("ConvertStructToMap did not convert struct pointer correctly")
	}

	nonStructResult := ConvertStructToMap(42)
	if nonStructResult != nil {
		t.Errorf("ConvertStructToMap should return nil for non-struct inputs")
	}

	nilResult := ConvertStructToMap(nil)
	if nilResult != nil {
		t.Errorf("ConvertStructToMap should return nil for nil inputs")
	}
}
