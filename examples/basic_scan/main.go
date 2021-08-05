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
	}

	if result != nil {
		for _, host := range result.Hosts {
			fmt.Printf("address : %s\n address type : %s\n", host.Address.Addr, host.Address.AddrType)
		}
	}
}
