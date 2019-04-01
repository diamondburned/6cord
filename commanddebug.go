package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"text/tabwriter"
	"time"
)

func commandDebug(text []string) {
	s := strings.Builder{}
	w := tabwriter.NewWriter(&s, 0, 0, 1, ' ', 0)

	fmt.Fprintf(w, "Channel ID:\t%d\n", Channel.ID)
	fmt.Fprintf(w, "Channel Icon:\t%s\n", Channel.Icon)
	fmt.Fprintf(w, "Guild ID:\t%d\n", Channel.GuildID)

	if g, _ := d.State.Guild(Channel.GuildID); g != nil {
		fmt.Fprintf(w,
			"Guild Icon:\thttps://cdn.discordapp.com/icons/%d/%s.png\n",
			g.ID, g.Icon,
		)
	}

	fmt.Fprintf(w, "Number of goroutines:\t%d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "GOMAXPROCS:\t%d\n", runtime.GOMAXPROCS(-1))
	fmt.Fprintf(w, "GOOS:\t%s\n", runtime.GOOS)
	fmt.Fprintf(w, "GOARCH:\t%s\n", runtime.GOARCH)
	fmt.Fprintf(w, "Go version:\t%s", runtime.Version())

	var gc = &debug.GCStats{}
	debug.ReadGCStats(gc)

	if gc != nil {
		fmt.Fprintf(w, "\nLast garbage collection:\t%s\n", gc.LastGC.Format(time.Kitchen))
		fmt.Fprintf(w, "Total garbage collection:\t%d\n", gc.NumGC)
		fmt.Fprintf(w, "Total pause for GC:\t%d", gc.PauseTotal)
	}

	var mem = &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	if mem != nil {
		fmt.Fprintf(w,
			"\nTotal RAM usage:\t%.2f MB\n",
			float64(mem.Alloc)/1000000,
		)
		fmt.Fprintf(w,
			"Total heap allocated:\t%.2f MB\n",
			float64(mem.HeapAlloc)/1000000,
		)
	}

	if err := w.Flush(); err != nil {
		Warn(err.Error())
		return
	}

	Message(s.String())

	runtime.GC()
}
