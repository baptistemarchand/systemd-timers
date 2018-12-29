package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"sort"

	"github.com/baptistemarchand/systemd-timers/systemd"
	"github.com/dustin/go-humanize"
	"github.com/reconquest/loreley"
)

const (
	oneSecond = 1000 * 1000
	oneMinute = 60 * oneSecond
)

var (
	headers = []string{
		"UNIT",
		"LAST",
		"RESULT",
		"TIME",
		"NEXT",
		"SCHEDULE",
	}
)

func generateTable(timers []*systemd.Timer) (string, error) {
	buf := &bytes.Buffer{}

	w := tabwriter.NewWriter(buf, 0, 0, 2, ' ', tabwriter.FilterHTML)
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	sort.Slice(timers, func (i, j int) bool {
		return timers[i].LastTriggered.Before(timers[j].LastTriggered)
	})

	for _, timer := range timers {
		var lastTriggered, result, nextElapse string

		if timer.LastTriggered.IsZero() {
			lastTriggered = ""
			result = ""
		} else {
			lastTriggered = fmt.Sprintf("%s (%s)", timer.LastTriggered.Local().String(), humanize.Time(timer.LastTriggered))

			if timer.Result == "success" {
				result = "<fg 2>✔<reset>"
			} else {
				result = "<fg 1>✘<reset>"
			}
		}

		if timer.NextElapse.IsZero() {
			nextElapse = ""
		} else {
			nextElapse = timer.NextElapse.Local().String()
		}

		fmt.Fprintln(w, strings.Join([]string{
			timer.Name,
			lastTriggered,
			result,
			formatExecutionTime(timer.LastExecutionTime),
			nextElapse,
			timer.Schedule,
		}, "\t"))
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
