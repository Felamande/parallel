package parallel

import "runtime"

type Worker struct {
	WorkerFunc func(interface{}) (interface{}, error)
	Input      interface{}
}

type ParallelTask struct {
	inputChan  chan Worker
	ResultChan chan interface{}
	ErrChan    chan error
	closeChan  chan bool
}

func New() *ParallelTask {
	CPUs := runtime.NumCPU()
	return &ParallelTask{
		inputChan:  make(chan Worker, CPUs),
		ResultChan: make(chan interface{}, CPUs),
		ErrChan:    make(chan error, CPUs),
		closeChan:  make(chan bool),
	}
}

func (pt *ParallelTask) Add(w Worker) {
	go func() {
		pt.inputChan <- w
	}()
}

func (pt *ParallelTask) not_impl_Close() {
	go func() {
		pt.closeChan <- true
	}()
}

func (pt *ParallelTask) Run(MaxParallelingTask int) {

	tokens := make(chan int, MaxParallelingTask)

	go func() {
		for {
			tokens <- 1
			work := <-pt.inputChan
			re, err := work.WorkerFunc(work.Input)
			if err != nil {
				pt.ErrChan <- err
				return
			}
			pt.ResultChan <- re
			<-tokens
		}
	}()
}
