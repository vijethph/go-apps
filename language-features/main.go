// code standard order is like this:
// 1. package declaration
// 2. import statements
// 3. source code
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type logWriter struct{}

func main() {
	fmt.Println("Hi there!")

	resp, err := http.Get("https://google.com")

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// make() - built-in function that takes a type of a slice and the number of elements the slice should have (zero values - initialized)
	// bs := make([]byte, 99999)
	// resp.Body.Read(bs)
	// fmt.Println(string(bs))

	// io.Copy(os.Stdout, resp.Body)

	lw := logWriter{}
	io.Copy(lw, resp.Body)
}

func (logWriter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))

	return len(bs), nil
}

func statusCheck() {
	links := []string{
		"https://google.com",
		"https://github.com",
		"https://amazon.com",
	}

	// channel of type string
	c := make(chan string)

	for _, link := range links {
		go checkLink(link, c)
	}

	// fmt.Println(<-c) // this is a blocking call. Main routine terminates execution if there are no other statements / no additional channel data sent by other routines. Fastest child routine to resolve first ends this blocking call.

	// alternative for loop syntax
	// for i := 0; i < len(links); i++ {
	// 	go checkLink(<-c, c)
	// }

	// infinite loop
	// for {
	// 	go checkLink(<-c, c)
	// }

	// alternative syntax for above infinite loop
	for l := range c { // equivalent of waiting for channel to return some value and assign it to iterating variable
		// go checkLink(l, c)

		go func(link string) {
			// Sleep() pauses the current goroutine for duration specified
			time.Sleep(5 * time.Second)
			checkLink(link, c)

		}(l)
	}
}

func checkLink(link string, c chan string) {

	_, err := http.Get(link)

	if err != nil {
		fmt.Println(link, "might be down!")
		// c <- "might be really down"
		c <- link
	}

	fmt.Println(link, "is up")
	// c <- "it is up"
	c <- link
}
