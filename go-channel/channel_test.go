package main

import (
	"fmt"
	"testing"
)

func TestAddUser(t *testing.T) {
	server := NewServer()
	server.Run()
	for i := 0; i < 10; i++ {
		// go server.AddUser(fmt.Sprintf("user_%d", i))
		// go func(i int) {
		server.userch <- fmt.Sprintf("user_%d", i)
		// }(i)

	}
	fmt.Println("Done")
}
