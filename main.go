package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/frort/ksk/internal/search"
	"github.com/frort/ksk/internal/tui"
)

func main() {
	engine := flag.String("e", "duckduckgo", "search engine (duckduckgo, brave)")
	region := flag.String("r", "", "region/country code (e.g. jp, us, de)")
	flag.Parse()

	var backend search.Backend
	switch *engine {
	case "duckduckgo", "ddg":
		backend = &search.DuckDuckGo{Region: *region}
	case "brave", "b":
		backend = &search.Brave{Region: *region}
	default:
		fmt.Fprintf(os.Stderr, "Unknown engine: %s (use duckduckgo or brave)\n", *engine)
		os.Exit(1)
	}

	query := strings.Join(flag.Args(), " ")

	m := tui.NewModel(query, backend)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
