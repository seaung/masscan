# Masscan

<p align="center">
	<img with="350" src="images/gopher.gif" />
</p>


## What is masscan

> masscan is an Internet-scale port scanner, useful for large scale surveys of the Internet, or of internal networks. While the default transmit rate is only 100 packets/second, it can optional go as fast as 25 million packets/second, a rate sufficient to scan the Internet in 3 minutes for one port.

## Installation

```bash
go get github.com/seaung/masscan

```

note: You need to install masscan on your computer in advance, otherwise you can't use this library

## Supported features

- [ x ] Support some parameters of masscan scanner

## TODO

- [] Constantly improve this library


## Simple example

```go
package main

import (
	"context"
	"fmt"
	"github.com/seaung/masscan"
	"time"
	"log"
)

func main() {
	cxt, cancel := context.WithTimeout(context.Background(), time.Minute * 5)
	defer cancel()

	masscanScanner, err := masscan.NewMasscanScanner(
		masscan.WithTargets("192.168.1.106"),
		masscan.WithPorts("0-8000"),
		masscan.WithBanners(),
		masscan.WithContext(cxt),
	)

	if err != nil {
		log.Fatalf("unable to create masscan scan: %v\n", err)
	}

	result, _, err := masscanScanner.Run()
	if err != nil {
		log.Fatalf("unable to run masscan scan: %v\n", err)
		return
	}
	fmt.Println(result)
}
```

## The development soul comes from
The development of this library is inspired by this library [Ullaakut/nmap](https://github1s.com/Ullaakut/nmap)

---
that's all
