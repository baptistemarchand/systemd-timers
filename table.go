package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"sort"
	"time"

	"github.com/baptistemarchand/systemd-timers/systemd"
	"github.com/dustin/go-humanize"
	"github.com/reconquest/loreley"
)

const (
	oneSecond = 1000 * 1000
	oneMinute = 60 * oneSecond
)

func generateTable(timers []*systemd.Timer, filters []string, verbose bool) (string, error) {
	buf := &bytes.Buffer{}

	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', tabwriter.FilterHTML)

	sort.Slice(timers, func (i, j int) bool {
		return timers[i].LastTriggered.Before(timers[j].LastTriggered)
	})

	now := time.Now()

	for _, timer := range timers {
		if !matchesFilters(timer.Name, filters) || timer.LastTriggered.IsZero() {
			continue
		}
		var lastTriggered, result string

		color := colorizeTime(timer.LastTriggered.Local(), now)
		lastTriggered = fmt.Sprintf("%s%s\t(%s)<reset>", color, timer.LastTriggered.Local().Format("15:04:05"), humanize.Time(timer.LastTriggered))

		if timer.Result == "success" {
			result = "<fg 2>✔<reset>"
		} else {
			result = "<fg 1>✘<reset>"
		}

		columns := []string{
			formatName(timer.Name),
			lastTriggered,
			result,
			formatExecutionTime(timer.LastExecutionTime),
		}

		fmt.Fprintln(w, strings.Join(columns, "\t"))
	}

	fmt.Fprintln(w, "--\t")
	fmt.Fprintln(w, fmt.Sprintf("now\t%s", now.Format("15:04:05")))
	fmt.Fprintln(w, "--\t")

	sort.Slice(timers, func (i, j int) bool {
		return timers[i].NextElapse.Before(timers[j].NextElapse)
	})

	for _, timer := range timers {
		if !matchesFilters(timer.Name, filters) || timer.NextElapse.IsZero() {
			continue
		}

		color := colorizeTime(timer.NextElapse.Local(), now)
		nextElapse := fmt.Sprintf("%s%s\t(%s)<reset>", color, timer.NextElapse.Local().Format("15:04:05"), humanize.Time(timer.NextElapse))

		columns := []string{
			formatName(timer.Name),
			nextElapse,
		}
		if verbose {
			columns = append(columns, timer.Schedule)
		}
		fmt.Fprintln(w, strings.Join(columns, "\t"))
	}

	w.Flush()

	loreley.DelimLeft = "<"
	loreley.DelimRight = ">"

	table, err := loreley.CompileAndExecuteToString(
		buf.String(),
		nil,
		nil,
	)
	if err != nil {
		return "", err
	}

	return table, nil
}

func matchesFilters(name string, filters []string) bool {
	for _, filter := range filters {
		if !strings.Contains(name, filter) {
			return false
		}
	}
	return true
}

func colorizeTime(t time.Time, now time.Time) string {
	var diff time.Duration

	if t.After(now) {
		diff = t.Sub(now)
	} else {
		diff = now.Sub(t)
	}

	if diff.Minutes() < 15 {
		return "<fg 1>"
	}
	if diff.Minutes() < 30 {
		return "<fg 3>"
	}
	if diff.Hours() < 1 {
		return "<fg 2>"
	}

	return ""
}

func formatName(name string) string {

	colors := map[string]string{
		"stats": "<fg 4>",
		"structure": "<fg 5>",
		"stripe": "<fg 6>",
		"system": "<fg 0><bg 6>",
	}

	for pattern, color := range colors {
		if strings.Contains(name, pattern) {
			return fmt.Sprintf("%s%s<reset>", color, name)
		}
	}

	return name
}

func formatExecutionTime(executionTime uint64) string {
	if executionTime == 0 {
		return ""
	}

	if executionTime < oneSecond {
		return "0s"
	}

	if executionTime < 60*oneSecond {
		return fmt.Sprintf("%ss", strconv.Itoa(int(executionTime/oneSecond)))
	}

	return fmt.Sprintf("<fg 1>%sm %ss<reset>", strconv.Itoa(int(executionTime/oneSecond/60)), strconv.Itoa(int((executionTime-oneMinute)/oneSecond%60)))
}
