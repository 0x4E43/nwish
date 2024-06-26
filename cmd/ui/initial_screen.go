package ui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().PaddingTop(3).MarginTop(5).BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()

	navStyle = func() lipgloss.Style {
		// b := lipgloss.NormalBorder().Bottom
		return lipgloss.NewStyle().MarginRight(10).Foreground(lipgloss.Color("#B2BEB5"))
	}()

	navStyleHl = func() lipgloss.Style {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5733"))
	}()

	tBarStyle = func() lipgloss.Style {
		return lipgloss.NewStyle().MarginTop(5)
	}()
)

type BaseScreenModel struct {
	term     string
	content  string
	ready    bool
	viewport viewport.Model
}

func (m BaseScreenModel) Init() tea.Cmd {
	return nil
}

func (m BaseScreenModel) View() string {
	if !m.ready {
		log.Println("Not Ready")
		return "Initializing, Please wait!"
	}
	log.Println(" Ready")

	return fmt.Sprintf("%s\n\n%s\n\n\n%s\n%s", m.headerView(), m.navView(), m.viewport.View(), m.footerView())
}

func (m BaseScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

		if k == "a" || k == "A" {

			log.Println("A Pressed")
			m.viewport.SetContent("A Pressed")
			return m, tea.ClearScrollArea
		}
		if k == "b" || k == "B" {
			log.Println("B Pressed")
			m.viewport.SetContent("B Pressed")
			return m, tea.ClearScrollArea
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			log.Println("inside WindowSizeMsg: NOT READY")
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			// m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			log.Println("inside WindowSizeMsg: READY")

			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m BaseScreenModel) headerView() string {
	title := titleStyle.Render("Nimai C. (Dev)")
	line := tBarStyle.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m BaseScreenModel) navView() string {
	aboutLnk := navStyleHl.Render("(a)")
	aboutTxt := navStyle.Render(" About")

	blogLnk := navStyleHl.Render("(b)")
	blogTxt := navStyle.Render(" Blogs")

	projLnk := navStyleHl.Render("(p)")
	projTxt := navStyle.Render(" Projects")

	repoLnk := navStyleHl.Render("(g)")
	repoTxt := navStyle.Render(" Active Repos")

	return lipgloss.JoinHorizontal(lipgloss.Center, aboutLnk, aboutTxt, blogLnk, blogTxt, projLnk, projTxt, repoLnk, repoTxt)
}

func (m BaseScreenModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := tBarStyle.Render(strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func RenderScreen(term string) BaseScreenModel {
	content, err := os.ReadFile("/home/nimai/nwish/md/info.md")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	scrn := BaseScreenModel{
		term:    term,
		content: string(content),
		ready:   true,
	}
	scrn.viewport.SetContent(scrn.content)
	return scrn
}
