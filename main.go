package main

import (
	"fmt"
	"os"
	"flag"

	"github.com/baptistemarchand/systemd-timers/systemd"
)

func main() {
	conn, err := systemd.NewConn()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	client := systemd.NewClient(conn)

	timers, err := client.ListTimers()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	verbose := flag.Bool("v", false, "Verbose mode: show schedule")

	flag.Parse()

	table, err := generateTable(timers, flag.Args(), *verbose)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(table)
}
