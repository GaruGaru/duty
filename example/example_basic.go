package main

import (
	"github.com/GaruGaru/duty/duty"
	"github.com/GaruGaru/duty/storage"
)

func main() {

	d := duty.New(storage.NewMemoryStorage(), duty.Default)

	d.Enqueue(PrintTask{Output: "Hello World!"})

}
