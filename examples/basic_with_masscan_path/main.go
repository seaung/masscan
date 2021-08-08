package main

import (
	"context"
	"fmt"
	"github.com/seaung/masscan"
	"time"
)

func main() {
	cxt, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)
	defer cancel()

	var path string = "/home/to/path/masscan/bin/masscan"

	MasscaScanner, err := masscan.NewMasscanScannerWithBinaryPath(
		path,
		masscan.WithTargets("192.168.1.106"),
		masscan.WithPorts("22,80"),
		masscan.WithTTL(5),
		masscan.WithContext(cxt),
	)

	if err != nil {
		fmt.Printf("error : %v\n", err)
	}

	result, _, err := MasscaScanner.Run()
	if err != nil {
		fmt.Printf("Can't run masscan : %+v\n", err.Error())
		return
	}

	fmt.Println(result)

}
