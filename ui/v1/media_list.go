package v1

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/dubeyKartikay/lazyspotify/core/logger"
)

type state int

const (
	initilized state = iota
	loading
	ready
)
type mediaList struct {
	kind   ListKind
	items  []Entity
	list   list.Model
	width  int
	height int
	title  string
	state  state
}


func newMediaList(kind ListKind) mediaList {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("252")).
		PaddingLeft(1)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(lipgloss.Color("245")).
		PaddingLeft(1)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("14")).
		Bold(true).
		BorderLeft(false).
		PaddingLeft(1)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("117")).
		BorderLeft(false).
		PaddingLeft(1)
	delegate.SetSpacing(1)
	listModel := list.New(nil, delegate, 0, 0)

	styles := listModel.Styles
	styles.Title = styles.Title.MarginLeft(1)

	styles.TitleBar = lipgloss.NewStyle().MarginBottom(1)
	styles.NoItems = styles.NoItems.Foreground(lipgloss.Color("8"))
	listModel.Styles = styles
	listModel.SetShowHelp(false)
	listModel.SetShowStatusBar(false)
	listModel.SetShowFilter(false)
	listModel.SetShowPagination(false)
	listModel.InfiniteScrolling = true
	listModel.Title = listTitle(kind)

	return mediaList{
		kind:   kind,
		list:   listModel,
		state:  initilized,
	}
}

func (m mediaList) View() string {
	listWidth := m.width - 4
	listHeight := m.height - 2
	m.list.SetSize(listWidth, listHeight)
	if(m.state == loading) {
		m.list.Title = "Loading..."
		m.list.SetItems(nil)
	}
	return m.list.View()
}

func (m *mediaList) SetTitle(title string) {
	m.list.Title = title
}

func (m *mediaList) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *mediaList) SetSize(width, height int) {
	m.width = width
	m.height = height

	listStyles := m.list.Styles
	m.list.Styles = listStyles
}
func (m *mediaList) StartLoading() tea.Cmd {
	m.state = loading
	return tea.Batch(
		m.list.StartSpinner(),
	)
}

func (m *mediaList) StopLoading() {
	m.state = ready
	m.list.StopSpinner()
}

func (m *mediaList) SetContent(entities []Entity, kind ListKind) tea.Cmd {
	items := make([]list.Item, 0, len(entities))
	for _, entity := range entities {
		if entity.Name == "" {
			continue
		}
		items = append(items, mediaListItem{entity: entity})
	}
	m.kind = kind
	m.items = entities
	setItemsCmd := m.list.SetItems(items)
	logger.Log.Info().Any("items", entities).Int("kind", int(kind)).Msg("set content")
	m.StopLoading()
	return setItemsCmd
}

func (m *mediaList) StopSpinner() {
	m.list.StopSpinner()
}

func (m *mediaList) SetStatus(message string) tea.Cmd {
	return m.list.NewStatusMessage(message)
}

type mediaListItem struct {
	entity Entity
}

func (i mediaListItem) Title() string {
	return i.entity.Name
}

func (i mediaListItem) Description() string {
	return i.entity.Desc
}

func (i mediaListItem) FilterValue() string {
	return fmt.Sprintf("%s %s", i.entity.Name, i.entity.Desc)
}

func listTitle(kind ListKind) string {
	switch kind {
	case Albums:
		return "Albums"
	case Artists:
		return "Artists"
	case Playlists:
		return "Playlists"
	case Tracks:
		return "Tracks"
	case Shows:
		return "Shows"
	case Episodes:
		return "Episodes"
	case AudioBooks:
		return "Audiobooks"
	default:
		return "Media"
	}
}

func listTitleAbbr(kind ListKind) string {
	switch kind {
	case Albums:
		return "AL"
	case Artists:
		return "AR"
	case Playlists:
		return "PL"
	case Tracks:
		return "TR"
	case Shows:
		return "SH"
	case Episodes:
		return "EP"
	case AudioBooks:
		return "AB"
	default:
		return "Media"
	}
}

func GenerateListTitle(kinds []ListKind) string {
	var parts []string
	for i := range len(kinds) - 1 {
		parts = append(parts, listTitleAbbr(kinds[i]))
	}
	parts = append(parts, listTitle(kinds[len(kinds)-1]))
	return strings.Join(parts, ">")
}

func (m *mediaList) IsEmpty() bool {
	return len(m.items) == 0
}
