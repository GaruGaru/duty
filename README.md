# Duty
### Go Library for managing long running persistent tasks

[![Build Status](https://travis-ci.org/GaruGaru/duty.svg?branch=master)](https://travis-ci.org/GaruGaru/duty)
[![Go Report Card](https://goreportcard.com/badge/github.com/GaruGaru/duty)](https://goreportcard.com/report/github.com/GaruGaru/duty)
![license](https://img.shields.io/github/license/GaruGaru/duty.svg)



**Define a task** 

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


**Schedule**

```


```