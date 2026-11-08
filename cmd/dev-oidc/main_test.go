package main

import (
	"reflect"
	"testing"
)

func TestSplitOrigins(t *testing.T) {
	got := splitOrigins(" http://localhost:4321/, http://localhost:4322,,", "http://localhost:4321")
	if want := []string{"http://localhost:4321", "http://localhost:4322"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("%v", got)
	}
	if got := splitOrigins("", "http://localhost:3000"); !reflect.DeepEqual(got, []string{"http://localhost:3000"}) {
		t.Fatalf("%v", got)
	}
}
