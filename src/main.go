package main

import (
	"flag"
	"fmt"
)

func main() {
	var host string
	var to string
	var help string
	flag.StringVar(&host, "host", "localhost", "the host to expose")
	flag.StringVar(&to, "to", "server", "address of the remote server")
	flag.String("help", "", "print help information")
	flag.Parse()

	
	
	fmt.Println(host, to, help)
	fmt.Println(flag.Args())
}
