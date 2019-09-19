package gosignal_test

import (
	"fmt"
	"testing"

	gs "github.com/RussellLuo/gosignal"
)

type receiver struct {
	Name string
	C    chan gs.Arguments
}

func (r *receiver) Receive(args gs.Arguments) {
	fmt.Println(r.Name)
	fmt.Printf("arguments: %+v\n", args)
}

func TestNew(t *testing.T) {
	sig1 := gs.New("")
	sig2 := gs.New("")
	sig3 := gs.New("ready")
	sig4 := gs.New("ready")

	if sig1 == sig2 {
		t.Error("Got: sig1 == sig2, Want: sig1 != sig2")
	}

	if sig2 == sig3 {
		t.Error("Got: sig2 == sig3, Want: sig2 != sig3")
	}

	if sig3 != sig4 {
		t.Error("Got: sig3 != sig4, Want: sig3 == sig4")
	}
}

func TestSignal_ConnectAndSend(t *testing.T) {
	r1 := &receiver{Name: "receiver1"}
	r2 := &receiver{Name: "receiver2"}

	sig := new(gs.Signal)

	sig.Connect(r1, r2)
	sig.Send(gs.Arguments{
		"greeting": "hello world",
	})

	sig.Disconnect()
	sig.Send(nil)
}
