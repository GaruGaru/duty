package main

import (
	"fmt"
	"github.com/GaruGaru/duty/duty"
	"github.com/GaruGaru/duty/storage"
)

func exampleBasicWithStatus() {

	d := duty.New(storage.NewMemoryStorage(), duty.Default)

	_ = d.Init()

	task, err := d.Enqueue(PrintTask{Output: "Hello World!"})

	if err != nil {
		panic(err)
	}

	task, _ = d.Get(task.ID)

	fmt.Println(task.Status.State)

}
