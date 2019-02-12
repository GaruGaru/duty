# Duty
### Go Library for managing long running persistent tasks

[![Build Status](https://travis-ci.org/GaruGaru/duty.svg?branch=master)](https://travis-ci.org/GaruGaru/duty)
[![Go Report Card](https://goreportcard.com/badge/github.com/GaruGaru/duty)](https://goreportcard.com/report/github.com/GaruGaru/duty)
![license](https://img.shields.io/github/license/GaruGaru/duty.svg)



**Task definition** 

```go

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

```


**Task scheduling**

```go
d := duty.New(storage.NewMemoryStorage(), duty.Default)
d.Enqueue(PrintTask{Output: "Hello World!"})
```



**Task monitoring (sync)**

```go
d := duty.New(storage.NewMemoryStorage(), duty.Default)

task, _ := d.Enqueue(PrintTask{Output: "Hello World!"})

task, _ = d.Get(task.ID)

fmt.Println(task.Status.State)

```


**Task monitoring (async)**

```go
	
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

```

