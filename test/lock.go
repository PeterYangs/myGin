package main

import (
	"fmt"
	"github.com/PeterYangs/tools/http"
	"sync"
)

func main() {

	client := http.Client()

	wait := sync.WaitGroup{}

	for i := 0; i < 4; i++ {

		wait.Add(1)

		go func(index int) {

			defer wait.Done()

			str, e := client.Request().GetToString("http://127.0.0.1:8080")

			fmt.Println("请求：", index, ",", str, e)

		}(i)

	}

	wait.Wait()

	fmt.Println("finish")

}
