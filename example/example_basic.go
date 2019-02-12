package main

import (
	"github.com/GaruGaru/duty/duty"
	"github.com/GaruGaru/duty/storage"
)

func exampleBasic() {

	d := duty.New(storage.NewMemoryStorage(), duty.Default)

	d.Enqueue(PrintTask{Output: "Hello World!"})

}
