package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/t-shah02/pacview/internal/utils"
)

const (
	titleLipglossColor = "99"
	helpLipglossColor  = "241"
	styleEdgeMargin    = 1

	defaultTermWidth   = 100
	defaultTermHeight  = 24
	initialTableHeight = 18

	minTermWidthForClamp  = 40
	fallbackTermWidth     = 80
	tableHorizontalGutter = 8
	minUsableInnerWidth   = 48
	minColumnWidth        = 6

	colWeightName        = 14
	colWeightDescription = 22
	colWeightVersion     = 12
	colWeightInstalledAt = 18
	colWeightDependsOn   = 18
	colWeightRequiredBy  = 16
	descriptionColumnIdx = 1

	// Title, search row (or bordered search when focused), help; keep slack for lipgloss borders.
	layoutChromeLines = 8
	minTableBodyHeight = 5

	searchInputWidthPadding = 4
	searchCharLimit         = 256

	keyQuit       = "q"
	keyQuitCtrl   = "ctrl+c"
	keySearchBlur = "esc"

	keyFocusSearch    = "/"
	keyFocusRequired  = "f"
	keyScopeBack      = "b"
)

var columnWidthWeights = [...]int{
	colWeightName,
	colWeightDescription,
	colWeightVersion,
	colWeightInstalledAt,
	colWeightDependsOn,
	colWeightRequiredBy,
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(titleLipglossColor)).
			MarginBottom(styleEdgeMargin)
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(helpLipglossColor)).
			MarginTop(styleEdgeMargin)
	searchBarFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(titleLipglossColor)).
				Padding(0, 1)
)

type model struct {
	all        []utils.PacmanPackage
	byName     map[string]utils.PacmanPackage
	scopeStack [][]utils.PacmanPackage
	displayed  []utils.PacmanPackage

	table         table.Model
	search        textinput.Model
	searchFocused bool

	w, h int
}

func Run(packages []utils.PacmanPackage) error {
	m := newModel(packages)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newModel(packages []utils.PacmanPackage) *model {
	ti := textinput.New()
	ti.Prompt = "search: "
	ti.Placeholder = "filter by name, version, description, deps…"
	ti.CharLimit = searchCharLimit
	ti.Width = defaultTermWidth - searchInputWidthPadding

	tbl := table.New(
		table.WithColumns(columnsForWidth(defaultTermWidth)),
		table.WithRows(nil),
		table.WithFocused(true),
		table.WithHeight(initialTableHeight),
		table.WithWidth(defaultTermWidth),
	)

	m := &model{
		all:        packages,
		byName:     indexPackagesByName(packages),
		scopeStack: [][]utils.PacmanPackage{packages},
		table:      tbl,
		search:     ti,
		w:          defaultTermWidth,
		h:          defaultTermHeight,
	}
	m.refreshTableRows()
	return m
}

func indexPackagesByName(pkgs []utils.PacmanPackage) map[string]utils.PacmanPackage {
	m := make(map[string]utils.PacmanPackage, len(pkgs))
	for _, p := range pkgs {
		m[p.Name] = p
	}
	return m
}

func packagesToRows(pkgs []utils.PacmanPackage) []table.Row {
	rows := make([]table.Row, len(pkgs))
	for i, p := range pkgs {
		rows[i] = table.Row{
			p.Name,
			p.Description,
			p.Version,
			p.InstalledAt,
			strings.Join(p.DependsOn, ", "),
			strings.Join(p.RequiredBy, ", "),
		}
	}
	return rows
}

func filterPackagesByQuery(pkgs []utils.PacmanPackage, query string) []utils.PacmanPackage {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		out := make([]utils.PacmanPackage, len(pkgs))
		copy(out, pkgs)
		return out
	}
	var out []utils.PacmanPackage
	for _, p := range pkgs {
		if packageMatchesQuery(p, q) {
			out = append(out, p)
		}
	}
	return out
}

func packageMatchesQuery(p utils.PacmanPackage, qLower string) bool {
	if strings.Contains(strings.ToLower(p.Name), qLower) {
		return true
	}
	if strings.Contains(strings.ToLower(p.Description), qLower) {
		return true
	}
	if strings.Contains(strings.ToLower(p.Version), qLower) {
		return true
	}
	if strings.Contains(strings.ToLower(p.InstalledAt), qLower) {
		return true
	}
	if strings.Contains(strings.ToLower(strings.Join(p.DependsOn, " ")), qLower) {
		return true
	}
	if strings.Contains(strings.ToLower(strings.Join(p.RequiredBy, " ")), qLower) {
		return true
	}
	return false
}

func packagesListedInRequiredBy(parent utils.PacmanPackage, byName map[string]utils.PacmanPackage) []utils.PacmanPackage {
	if len(parent.RequiredBy) == 0 {
		return nil
	}
	var out []utils.PacmanPackage
	for _, name := range parent.RequiredBy {
		if pkg, ok := byName[name]; ok {
			out = append(out, pkg)
		}
	}
	return out
}

func (m *model) currentScope() []utils.PacmanPackage {
	return m.scopeStack[len(m.scopeStack)-1]
}

func (m *model) refreshTableRows() {
	scope := m.currentScope()
	filtered := filterPackagesByQuery(scope, m.search.Value())
	m.displayed = filtered
	m.table.SetRows(packagesToRows(filtered))
	if m.table.Cursor() >= len(filtered) {
		m.table.SetCursor(max(0, len(filtered)-1))
	}
}

func (m *model) focusRequiredBySubset() {
	if len(m.displayed) == 0 {
		return
	}
	i := m.table.Cursor()
	if i < 0 || i >= len(m.displayed) {
		return
	}
	next := packagesListedInRequiredBy(m.displayed[i], m.byName)
	if len(next) == 0 {
		return
	}
	m.scopeStack = append(m.scopeStack, next)
	m.search.SetValue("")
	m.refreshTableRows()
	m.table.SetCursor(0)
}

func (m *model) popScope() {
	if len(m.scopeStack) <= 1 {
		return
	}
	m.scopeStack = m.scopeStack[:len(m.scopeStack)-1]
	m.refreshTableRows()
	m.table.SetCursor(0)
}

func columnsForWidth(termW int) []table.Column {
	if termW < minTermWidthForClamp {
		termW = fallbackTermWidth
	}
	usable := max(termW-tableHorizontalGutter, minUsableInnerWidth)

	sumW := 0
	for _, w := range columnWidthWeights {
		sumW += w
	}
	widths := make([]int, len(columnWidthWeights))
	total := 0
	for i, w := range columnWidthWeights {
		widths[i] = usable * w / sumW
		if widths[i] < minColumnWidth {
			widths[i] = minColumnWidth
		}
		total += widths[i]
	}
	if total < usable {
		widths[descriptionColumnIdx] += usable - total
	}
	titles := []string{
		"Name",
		"Description",
		"Version",
		"Installed at",
		"Depends on",
		"Required by",
	}
	cols := make([]table.Column, len(titles))
	for i := range titles {
		cols[i] = table.Column{Title: titles[i], Width: widths[i]}
	}
	return cols
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case keyQuit, keyQuitCtrl:
			return m, tea.Quit
		}

		if m.searchFocused {
			if msg.String() == keySearchBlur {
				m.searchFocused = false
				m.search.Blur()
				m.table.Focus()
				return m, nil
			}
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			m.refreshTableRows()
			return m, cmd
		}

		switch msg.String() {
		case keyFocusSearch:
			m.searchFocused = true
			m.table.Blur()
			return m, m.search.Focus()
		case keyFocusRequired:
			m.focusRequiredBySubset()
			return m, nil
		case keyScopeBack:
			m.popScope()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.search.Width = max(m.w-searchInputWidthPadding, minUsableInnerWidth/2)
		m.table.SetColumns(columnsForWidth(m.w))
		m.table.SetWidth(m.w)
		bodyH := max(m.h-layoutChromeLines, minTableBodyHeight)
		m.table.SetHeight(bodyH)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) helpText() string {
	narrowed := ""
	if len(m.scopeStack) > 1 {
		narrowed = " · narrowed (b back)"
	}
	return fmt.Sprintf(
		"%d shown · / search · f required-by of row · b back · ↑/↓ scroll · esc leave search · q quit%s",
		len(m.displayed),
		narrowed,
	)
}

func (m *model) View() string {
	title := titleStyle.Render("pacview — installed packages")
	searchLine := m.search.View()
	if m.searchFocused {
		searchLine = searchBarFocusedStyle.Render(searchLine)
	}
	help := helpStyle.Render(m.helpText())
	return lipgloss.JoinVertical(lipgloss.Left, title, searchLine, m.table.View(), help)
}
