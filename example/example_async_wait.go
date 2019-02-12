package main

import (
	"fmt"
	"github.com/GaruGaru/duty/duty"
	"github.com/GaruGaru/duty/pool"
	"github.com/GaruGaru/duty/storage"
	"sync"
)

type PrintTask struct {
	Output string
}

func (t PrintTask) Run() error {
	fmt.Println(t.Output)
	return nil
}

func (PrintTask) Type() string {
	return "print"
}

func exampleAsyncWait() {

	wg := sync.WaitGroup{}
	wg.Add(1)

	d := duty.New(storage.NewMemoryStorage(), duty.Options{
		ResultCallback: func(result pool.ScheduledTaskResult) {
			fmt.Println("Got task status " + result.ScheduledTask.ID + ": " + result.Status.State)
			if result.Status.Completed {
				wg.Done()
			}
		},
	})

	d.Init()

	d.Enqueue(PrintTask{
		Output: "test",
	})

	wg.Wait()

}
