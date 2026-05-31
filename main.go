package main

import (
	"fmt"
	"net"
	"sort"
	"time"
	"os"
	"sync"

	"github.com/pterm/pterm"
)

func worker(ports, results chan int, target string, wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		fmt.Printf("Checked port: %d\n", p) //QA

		address := fmt.Sprintf("%s:%d", target, p)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err != nil {
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	var wg sync.WaitGroup
	if len(os.Args) < 2 {
		pterm.FgRed.Printf("Usage: ./mantis <IP or Domain>")
		return
	}

	target := os.Args[1]

	var workers int
	for {
		pterm.FgCyan.Printf("How many workers would you like to use?\n")
		_, err := fmt.Scanln(&workers)
		if err == nil {
			break
		} else {
			pterm.FgYellow.Printf("Please input a valid number")
		}
	}

	ports := make(chan int, workers)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		wg.Add(1)
		go worker(ports, results, target, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for i := 1; i <= 65535; i++ {
			ports <- i
		}
		close(ports)
	}()

	for port := range results {
		openports = append(openports, port)
	}

	fmt.Println("--------------------------------")

	if len(openports) == 0 {
		fmt.Println("No open ports available")
	} else {
		sort.Ints(openports)
		for _, port := range openports {
			coloredPort := pterm.FgCyan.Sprintf("%d", port)
			fmt.Printf("%s open\n", coloredPort)
		}
	}
}
