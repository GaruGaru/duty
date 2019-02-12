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


func main() {

	wg := sync.WaitGroup{}

	wg.Add(1)

	d := duty.New(storage.NewMemoryStorage(), duty.Options{
		ResultCallback: func(result pool.ScheduledTaskResult) {

			fmt.Println("Got status " + result.Status.State)
			if result.Status.Completed {
				wg.Done()
			}

		},
	})

	err := d.Init()

	if err != nil {
		panic("unable to schedule task")
	}

	task, err := d.Enqueue(PrintTask{
		Output: "test",
	})

	if err != nil {
		panic("unable to schedule task")
	}

	fmt.Println("task scheduled with id " + task.ID)

	wg.Wait()

}
