package spass

import (
	"testing"
)

var (
	shadow = &Shadow{}
	s      *[]byte
)

func ExampleRead_Shadow() {
	shdw := &Shadow{}

	// do not reset my password
	err := shdw.Read(0)
	if err != nil {
		panic(err)
	}
	println("My password is:", *shdw.Get())
}

func TestStoreString_Shadow(t *testing.T) {
	shadow.StoreString("this is my password")
}

func TestGet_Shadow(t *testing.T) {
	s = shadow.Get()
	println(string(*s))
}

func TestClean_Shadow(t *testing.T) {
	shadow.Clean()
	println(string(*s))
}
