package goroutine

import "fmt"

func Go(f func()) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("panic error." + err.(string))
		}
	}()
	go f()
}
