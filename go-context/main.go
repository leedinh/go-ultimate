package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	ctx := context.Background()
	val, err := fetchdata(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Time taken: ", time.Since(start))
	fmt.Println(val)
}

type Response struct {
	val int
	err error
}

func fetchdata(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	respch := make(chan Response)
	go func() {
		val, err := longrunfunction()
		respch <- Response{val, err}
	}()

	select {
	case <-ctx.Done():
		return 0, fmt.Errorf("function takes too long to complete")
	case resp := <-respch:
		return resp.val, resp.err
	}

}

func longrunfunction() (int, error) {
	time.Sleep(time.Millisecond * 500)
	return 444, nil
}
