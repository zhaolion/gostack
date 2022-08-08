package main

import (
	"fmt"
	"time"

	"github.com/zhaolion/gostack/util/svcutil"
)

func main() {
	svcutil.WaitFor(":8086", func(stop <-chan struct{}) error {
		for {
			select {
			case <-stop:
				return nil
			default:
				fmt.Println(time.Now())
			}
		}
	})
}
