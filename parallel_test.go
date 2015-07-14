package parallel

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"testing"
	"time"
)

type TestStruct struct {
	length int
	url    string
}

func TestParallel(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	pt := New()
	pt.Run(4)

	//var wg sync.WaitGroup

	go func() {
		for {
			select {
			case re := <-pt.ResultChan:
				fmt.Println(re)
			case err := <-pt.ErrChan:
				fmt.Println(err)
			}
		}
	}()
	pt.Add(
		Worker{
			WorkerFunc: Task,
			Input:      "http://www.baidu.com",
		},
	)
	pt.Add(
		Worker{
			WorkerFunc: Task,
			Input:      "http://www.bilibili.com",
		},
	)
	pt.Add(
		Worker{
			WorkerFunc: Task,
			Input:      "https://v2ex.com",
		},
	)

	time.Sleep(time.Millisecond * 600)

}

func Task(input interface{}) (re interface{}, err error) {
	url := input.(string)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return TestStruct{len(d), url}, nil
}
