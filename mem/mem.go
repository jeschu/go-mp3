package mem

import (
	"fmt"
	"golang.org/x/text/message"
	"runtime"
	"strings"
)

var p = message.NewPrinter(message.MatchLanguage("de_DE"))

func MemStats() string {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	s := strings.Builder{}
	s.WriteString(p.Sprintf("HeapAlloc : %s\n", memory(stats.HeapAlloc)))
	s.WriteString(p.Sprintf("HeapSys   : %s\n", memory(stats.HeapSys)))
	s.WriteString(p.Sprintf("StackInuse: %s\n", memory(stats.StackInuse)))
	s.WriteString(p.Sprintf("StackSys  : %s\n", memory(stats.StackSys)))
	s.WriteString(p.Sprintf("Sys       : %s\n", memory(stats.Sys)))
	s.WriteString(p.Sprintf("Objects created: %s\n", formatUint64(stats.Mallocs)))
	s.WriteString(p.Sprintf("Objects deleted: %s\n", formatUint64(stats.Frees)))
	s.WriteString(p.Sprintf("Objects in use : %s\n", formatUint64(stats.Mallocs-stats.Frees)))
	return s.String()
}

func formatFloat64(value float64) string {
	return fixFormat(p.Sprintf("%10.1f", value))
}
func formatUint64(value uint64) string {
	return fixFormat(p.Sprintf("%10d", value))
}

func fixFormat(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(s, ",", "§"),
			".", ","),
		"§", ".")
}

var units = []string{"B", "MB", "GB", "TB"}

func memory(value uint64) string {
	v := float64(value)
	i := 0
	for ; v > 1024 && i < len(units); i++ {
		v /= 1024
	}
	return fmt.Sprintf("%s%s", formatFloat64(v), units[i])
}
