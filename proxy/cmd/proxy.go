package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/intex-software/goutils/proxy"
)

func main() {
	flags := flag.NewFlagSet("proxy", flag.ExitOnError)
	target := flags.String("target", "https://www.example.com", "URL des Proxys")
	addr := flags.String("addr", ":8080", "Port des Proxys")
	_ = flags.Parse(os.Args[1:])

	fmt.Println("Starting proxy server")
	proxy.Proxy(*target, *addr)
	fmt.Println("Stopped proxy server")
}
