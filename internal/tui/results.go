package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/frort/ksk/internal/search"
)

type resultsModel struct {
	results []search.Result
	cursor  int
	offset  int // viewport scroll offset
	pageNum int
	hasMore bool
	width   int
	height  int
}

func newResultsModel() resultsModel {
	return resultsModel{}
}

func (m *resultsModel) SetResults(results []search.Result, pageNum int, hasMore bool) {
	m.results = results
	m.cursor = 0
	m.offset = 0
	m.pageNum = pageNum
	m.hasMore = hasMore
}

func (m *resultsModel) CursorDown() {
	if m.cursor < len(m.results)-1 {
		m.cursor++
		m.ensureVisible()
	}
}

func (m *resultsModel) CursorUp() {
	if m.cursor > 0 {
		m.cursor--
		m.ensureVisible()
	}
}

func (m *resultsModel) CursorTop() {
	m.cursor = 0
	m.offset = 0
}

func (m *resultsModel) CursorBottom() {
	if len(m.results) > 0 {
		m.cursor = len(m.results) - 1
		m.ensureVisible()
	}
}

func (m *resultsModel) SelectedResult() *search.Result {
	if len(m.results) == 0 {
		return nil
	}
	return &m.results[m.cursor]
}

func (m *resultsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// visibleFrom counts how many results fit in the viewport starting from startIdx.
func (m *resultsModel) visibleFrom(startIdx int) int {
	if m.height <= 0 || len(m.results) == 0 {
		return len(m.results)
	}
	totalHeight := 0
	count := 0
	for i := startIdx; i < len(m.results); i++ {
		block := m.renderBlock(i)
		blockH := lipgloss.Height(block) + 1
		if totalHeight+blockH > m.height && count > 0 {
			break
		}
		totalHeight += blockH
		count++
	}
	return max(1, count)
}

func (m *resultsModel) ensureVisible() {
	if m.cursor < m.offset {
		m.offset = m.cursor
		return
	}
	// Scroll down: shrink offset until cursor fits in viewport
	for m.offset < m.cursor {
		visible := m.visibleFrom(m.offset)
		if m.cursor < m.offset+visible {
			break
		}
		m.offset++
	}
}

func (m *resultsModel) renderBlock(i int) string {
	r := m.results[i]
	selected := i == m.cursor

	contentWidth := m.contentWidth()
	// Text width = content width minus padding (1 left + 1 right)
	textWidth := contentWidth - 2

	title := truncate(r.Title, textWidth)
	url := truncate(r.URL, textWidth)
	snippet := truncate(r.Snippet, textWidth)

	var titleRendered, urlRendered, snippetRendered string
	if selected {
		titleRendered = selectedTitleStyle.Render(title)
	} else {
		titleRendered = titleStyle.Render(title)
	}
	urlRendered = urlStyle.Render(url)
	snippetRendered = snippetStyle.Render(snippet)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleRendered,
		urlRendered,
		snippetRendered,
	)

	blockStyle := resultBlock.Width(contentWidth)
	if selected {
		blockStyle = selectedBlock.Width(contentWidth)
	}
	return blockStyle.Render(content)
}

func (m *resultsModel) contentWidth() int {
	// Account for border (2) + outer margin/space (2)
	cw := m.width - 4
	if cw < 20 {
		cw = 60
	}
	return cw
}

func (m *resultsModel) View() string {
	if len(m.results) == 0 {
		return "\n  No results found.\n"
	}

	var b strings.Builder
	totalHeight := 0

	for i := m.offset; i < len(m.results); i++ {
		block := m.renderBlock(i)
		blockH := lipgloss.Height(block) + 1 // +1 for separator newline

		if totalHeight+blockH > m.height && totalHeight > 0 {
			break
		}

		b.WriteString(block)
		b.WriteString("\n")
		totalHeight += blockH
	}

	return b.String()
}

func (m *resultsModel) StatusView(engineName string) string {
	if len(m.results) == 0 {
		return ""
	}
	return fmt.Sprintf("[%s] Page %d | %d/%d | j/k:move h/l:page Enter:open /:search q:quit",
		engineName, m.pageNum, m.cursor+1, len(m.results))
}

func truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) > maxWidth {
		if maxWidth > 3 {
			return string(runes[:maxWidth-3]) + "..."
		}
		return string(runes[:maxWidth])
	}
	return s
}
