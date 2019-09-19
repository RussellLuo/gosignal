package gosignal_test

import (
	"fmt"

	gs "github.com/RussellLuo/gosignal"
)

func Example_funcReceiver() {
	sig := gs.New("")
	sig.Connect(gs.FuncReceiver(func(args gs.Arguments) {
		fmt.Println("Signal received")
	}))
	defer sig.Disconnect()

	sig.Send(nil)

	// Output:
	// Signal received
}

type structReceiver struct{}

func (r *structReceiver) Receive(args gs.Arguments) {
	fmt.Printf("arguments: %+v\n", args)
}

func Example_structReceiver() {
	sig := gs.New("")
	sig.Connect(&structReceiver{})
	defer sig.Disconnect()

	sig.Send(gs.Arguments{
		"greeting": "hello world",
	})

	// Output:
	// arguments: map[greeting:hello world]
}

func Example_namedSignal() {
	sig := gs.New("ready")
	sig.Connect(gs.FuncReceiver(func(args gs.Arguments) {
		fmt.Println("Signal received")
	}))
	defer sig.Disconnect()

	exitC := make(chan struct{}, 1)
	go func() {
		sig := gs.New("ready")
		sig.Send(nil)

		exitC <- struct{}{}
	}()

	<-exitC

	// Output:
	// Signal received
}
