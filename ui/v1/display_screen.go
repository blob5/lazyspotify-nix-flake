package v1

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/dubeyKartikay/lazyspotify/core/ticker"
	ansi "github.com/charmbracelet/x/ansi"
)

type displayScreen struct {
	width          int
	height         int
	display string
	scrollOffset   int
	styles         displayStyles
}

type displayStyles struct {
	panel   lipgloss.Style
	primary lipgloss.Style
	accent  lipgloss.Style
	muted   lipgloss.Style
	marquee lipgloss.Style
}

func newDisplayScreen() displayScreen {
	return displayScreen{
		display: "Lazyspotify: The cutest terminal music player",
		styles: displayStyles{
			panel: lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("30")).
				Foreground(lipgloss.Color("229")),
			primary: lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Bold(true),
			accent:  lipgloss.NewStyle().Foreground(lipgloss.Color("50")),
			muted:   lipgloss.NewStyle().Foreground(lipgloss.Color("151")),
			marquee: lipgloss.NewStyle().Foreground(lipgloss.Color("195")).Bold(true),
		},
	}
}

func (d *displayScreen) SetDisplayFromSong(songInfo SongInfo) {
	if songInfo.title == "" {
		return
	}
	separator := " • "
	styled := lipgloss.JoinHorizontal(
		lipgloss.Left,
		d.styles.primary.Render(songInfo.title),
		d.styles.accent.Render(separator),
		d.styles.muted.Render(songInfo.artist),
		d.styles.accent.Render(separator),
		d.styles.muted.Render(songInfo.album),
	)
	d.display = styled
}

func (d *displayScreen) SetDisplay(s string) {
	d.display = s
}

func (d *displayScreen) View() string {
	raw := d.display
	contentWidth := max(0, d.width-2)
	styled := d.styles.muted.Render(raw)
	if contentWidth > 0 {
		if lipgloss.Width(raw) > contentWidth {
			styled = d.styles.marquee.Render(d.scrollText(raw, contentWidth))
		}
		styled = lipgloss.NewStyle().Width(contentWidth).Align(lipgloss.Center).Render(styled)
	}
	panel := d.styles.panel.Width(d.width).Height(d.height).Render(styled)
	return panel
}

func (d *displayScreen) SetSize(width int, height int) {
	d.width = width
	d.height = height
}

func (d *displayScreen) NextFrame() tea.Cmd {
	d.scrollOffset = (d.scrollOffset + 1)%len(d.display)
	return ticker.DoTick()
}

func (d *displayScreen) scrollText(text string, width int) string {
	if width <= 0 {
		return ""
	}

	if lipgloss.Width(text) <= width {
		return text
	}

	const gap = "   "
	base := text + gap
	track := base + base
	if len(base) == 0 {
		return strings.Repeat(" ", width)
	}

	start := d.scrollOffset
	end := min(len(track), start + width)
	visible := ansi.Cut(track, start, end)
	return visible
}
