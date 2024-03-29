package goya

import (
	"testing"
)

type person struct {
	Name string
	Age  int
}

func TestConvertStructToMap(t *testing.T) {
	p := person{Name: "John", Age: 30}
	mapResult := convertStructToMap(p)
	if mapResult == nil {
		t.Errorf("convertStructToMap returned nil for struct input")
	}
	if mapResult["Name"] != "John" || mapResult["Age"] != 30 {
		t.Errorf("convertStructToMap did not convert struct correctly")
	}

	ptrResult := convertStructToMap(&p)
	if ptrResult == nil {
		t.Errorf("convertStructToMap returned nil for pointer to struct input")
	}
	if ptrResult["Name"] != "John" || ptrResult["Age"] != 30 {
		t.Errorf("convertStructToMap did not convert struct pointer correctly")
	}

	nonStructResult := convertStructToMap(42)
	if nonStructResult != nil {
		t.Errorf("convertStructToMap should return nil for non-struct inputs")
	}

	nilResult := convertStructToMap(nil)
	if nilResult != nil {
		t.Errorf("convertStructToMap should return nil for nil inputs")
	}
}
