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

	flag.Parse()

	table, err := generateTable(timers, flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(table)
}
