package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type palette struct {
	bgMain          color.Color
	bgDim           color.Color
	bgActive        color.Color
	bgInactive      color.Color
	accent          color.Color
	accentMuted     color.Color
	accentSurface   color.Color
	accentWarm      color.Color
	accentWarmMuted color.Color
	highlight       color.Color
	highlightMuted  color.Color
	hash            color.Color
	info            color.Color
	textStrong      color.Color
	textBase        color.Color
	textMuted       color.Color
	textSubtle      color.Color
	textFaint       color.Color
	textDim         color.Color
	textCool        color.Color
	textTitle       color.Color
	textTitleActive color.Color
	border          color.Color
	borderMuted     color.Color
	borderStrong    color.Color
	errorSurface    color.Color
	success         color.Color
	danger          color.Color
}

func defaultPalette() palette {
	return palette{
		bgMain:          lipgloss.Color("#000000"),
		bgDim:           lipgloss.Color("#1e1e1e"),
		bgActive:        lipgloss.Color("#535353"),
		bgInactive:      lipgloss.Color("#303030"),
		accent:          lipgloss.Color("#2fafff"),
		accentMuted:     lipgloss.Color("#79a8ff"),
		accentSurface:   lipgloss.Color("#535353"),
		accentWarm:      lipgloss.Color("#d0bc00"),
		accentWarmMuted: lipgloss.Color("#fec43f"),
		highlight:       lipgloss.Color("#c6daff"),
		highlightMuted:  lipgloss.Color("#82b0ec"),
		hash:            lipgloss.Color("#b6a0ff"),
		info:            lipgloss.Color("#00d3d0"),
		textStrong:      lipgloss.Color("#ffffff"),
		textBase:        lipgloss.Color("#ffffff"),
		textMuted:       lipgloss.Color("#989898"),
		textSubtle:      lipgloss.Color("#646464"),
		textFaint:       lipgloss.Color("#646464"),
		textDim:         lipgloss.Color("#989898"),
		textCool:        lipgloss.Color("#c6daff"),
		textTitle:       lipgloss.Color("#76afbf"),
		textTitleActive: lipgloss.Color("#00bcff"),
		border:          lipgloss.Color("#646464"),
		borderMuted:     lipgloss.Color("#303030"),
		borderStrong:    lipgloss.Color("#535353"),
		errorSurface:    lipgloss.Color("#ff5f59"),
		success:         lipgloss.Color("#44bc44"),
		danger:          lipgloss.Color("#ff5f59"),
	}
}

// XXX: For now, this is in its own package so that it can be shared between
// different packages without incurring an illegal import cycle.

// Styles defines styles for the UI.
type Styles struct {
	ActiveBorderColor   color.Color
	InactiveBorderColor color.Color

	App                  lipgloss.Style
	ServerName           lipgloss.Style
	TopLevelNormalTab    lipgloss.Style
	TopLevelActiveTab    lipgloss.Style
	TopLevelActiveTabDot lipgloss.Style

	MenuItem       lipgloss.Style
	MenuLastUpdate lipgloss.Style

	RepoSelector struct {
		Normal struct {
			Base    lipgloss.Style
			Title   lipgloss.Style
			Desc    lipgloss.Style
			Command lipgloss.Style
			Updated lipgloss.Style
		}
		Active struct {
			Base    lipgloss.Style
			Title   lipgloss.Style
			Desc    lipgloss.Style
			Command lipgloss.Style
			Updated lipgloss.Style
		}
	}

	Repo struct {
		Base       lipgloss.Style
		Title      lipgloss.Style
		Command    lipgloss.Style
		Body       lipgloss.Style
		Header     lipgloss.Style
		HeaderName lipgloss.Style
		HeaderDesc lipgloss.Style
	}

	Footer      lipgloss.Style
	Branch      lipgloss.Style
	HelpKey     lipgloss.Style
	HelpValue   lipgloss.Style
	HelpDivider lipgloss.Style
	URLStyle    lipgloss.Style

	Error      lipgloss.Style
	ErrorTitle lipgloss.Style
	ErrorBody  lipgloss.Style

	LogItem struct {
		Normal struct {
			Base    lipgloss.Style
			Hash    lipgloss.Style
			Title   lipgloss.Style
			Desc    lipgloss.Style
			Keyword lipgloss.Style
		}
		Active struct {
			Base    lipgloss.Style
			Hash    lipgloss.Style
			Title   lipgloss.Style
			Desc    lipgloss.Style
			Keyword lipgloss.Style
		}
	}

	Log struct {
		Commit         lipgloss.Style
		CommitHash     lipgloss.Style
		CommitAuthor   lipgloss.Style
		CommitDate     lipgloss.Style
		CommitBody     lipgloss.Style
		CommitStatsAdd lipgloss.Style
		CommitStatsDel lipgloss.Style
		Paginator      lipgloss.Style
	}

	Ref struct {
		Normal struct {
			Base     lipgloss.Style
			Item     lipgloss.Style
			ItemTag  lipgloss.Style
			ItemDesc lipgloss.Style
			ItemHash lipgloss.Style
		}
		Active struct {
			Base     lipgloss.Style
			Item     lipgloss.Style
			ItemTag  lipgloss.Style
			ItemDesc lipgloss.Style
			ItemHash lipgloss.Style
		}
		ItemSelector lipgloss.Style
		Paginator    lipgloss.Style
		Selector     lipgloss.Style
	}

	Tree struct {
		Normal struct {
			FileName lipgloss.Style
			FileDir  lipgloss.Style
			FileMode lipgloss.Style
			FileSize lipgloss.Style
		}
		Active struct {
			FileName lipgloss.Style
			FileDir  lipgloss.Style
			FileMode lipgloss.Style
			FileSize lipgloss.Style
		}
		Selector    lipgloss.Style
		FileContent lipgloss.Style
		Paginator   lipgloss.Style
		Blame       struct {
			Hash    lipgloss.Style
			Message lipgloss.Style
			Who     lipgloss.Style
		}
	}

	Stash struct {
		Normal struct {
			Message lipgloss.Style
		}
		Active struct {
			Message lipgloss.Style
		}
		Title    lipgloss.Style
		Selector lipgloss.Style
	}

	Spinner          lipgloss.Style
	SpinnerContainer lipgloss.Style

	NoContent lipgloss.Style

	StatusBar       lipgloss.Style
	StatusBarKey    lipgloss.Style
	StatusBarValue  lipgloss.Style
	StatusBarInfo   lipgloss.Style
	StatusBarBranch lipgloss.Style
	StatusBarHelp   lipgloss.Style

	Tabs         lipgloss.Style
	TabInactive  lipgloss.Style
	TabActive    lipgloss.Style
	TabSeparator lipgloss.Style

	Code struct {
		LineDigit lipgloss.Style
		LineBar   lipgloss.Style
	}
}

// DefaultStyles returns default styles for the UI.
func DefaultStyles() *Styles {
	p := defaultPalette()

	s := new(Styles)

	s.ActiveBorderColor = p.accent
	s.InactiveBorderColor = p.textSubtle

	s.App = lipgloss.NewStyle().
		Margin(1, 2)

	s.ServerName = lipgloss.NewStyle().
		Height(1).
		MarginLeft(1).
		MarginBottom(1).
		Padding(0, 1).
		Background(p.accentSurface).
		Foreground(p.textStrong).
		Bold(true)

	s.TopLevelNormalTab = lipgloss.NewStyle().
		MarginRight(2)

	s.TopLevelActiveTab = s.TopLevelNormalTab.
		Foreground(p.accent)

	s.TopLevelActiveTabDot = lipgloss.NewStyle().
		Foreground(p.accent)

	s.RepoSelector.Normal.Base = lipgloss.NewStyle().
		PaddingLeft(1).
		Border(lipgloss.Border{Left: " "}, false, false, false, true).
		Height(3)

	s.RepoSelector.Normal.Title = lipgloss.NewStyle().Bold(true)

	s.RepoSelector.Normal.Desc = lipgloss.NewStyle().
		Foreground(p.textMuted)

	s.RepoSelector.Normal.Command = lipgloss.NewStyle().
		Foreground(p.accentMuted)

	s.RepoSelector.Normal.Updated = lipgloss.NewStyle().
		Foreground(p.textMuted)

	s.RepoSelector.Active.Base = s.RepoSelector.Normal.Base.
		BorderStyle(lipgloss.Border{Left: "┃"}).
		BorderForeground(p.accent)

	s.RepoSelector.Active.Title = s.RepoSelector.Normal.Title.
		Foreground(p.highlight)

	s.RepoSelector.Active.Desc = s.RepoSelector.Normal.Desc.
		Foreground(p.textCool)

	s.RepoSelector.Active.Updated = s.RepoSelector.Normal.Updated.
		Foreground(p.highlight)

	s.RepoSelector.Active.Command = s.RepoSelector.Normal.Command.
		Foreground(p.accentMuted)

	s.MenuItem = lipgloss.NewStyle().
		PaddingLeft(1).
		Border(lipgloss.Border{
			Left: " ",
		}, false, false, false, true).
		Height(3)

	s.MenuLastUpdate = lipgloss.NewStyle().
		Foreground(p.textSubtle).
		Align(lipgloss.Right)

	s.Repo.Base = lipgloss.NewStyle()

	s.Repo.Title = lipgloss.NewStyle().
		Padding(0, 2)

	s.Repo.Command = lipgloss.NewStyle().
		Foreground(p.accentMuted)

	s.Repo.Body = lipgloss.NewStyle().
		Margin(1, 0)

	s.Repo.Header = lipgloss.NewStyle().
		MaxHeight(2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(p.border)

	s.Repo.HeaderName = lipgloss.NewStyle().
		Foreground(p.highlight).
		Bold(true)

	s.Repo.HeaderDesc = lipgloss.NewStyle().
		Foreground(p.textMuted)

	s.Footer = lipgloss.NewStyle().
		MarginTop(1).
		Padding(0, 1).
		Height(1)

	s.Branch = lipgloss.NewStyle().
		Foreground(p.accentWarmMuted).
		Background(p.border).
		Padding(0, 1)

	s.HelpKey = lipgloss.NewStyle().
		Foreground(p.textSubtle)

	s.HelpValue = lipgloss.NewStyle().
		Foreground(p.textFaint)

	s.HelpDivider = lipgloss.NewStyle().
		Foreground(p.borderMuted).
		SetString(" • ")

	s.URLStyle = lipgloss.NewStyle().
		MarginLeft(1).
		Foreground(p.accentMuted)

	s.Error = lipgloss.NewStyle().
		MarginTop(2)

	s.ErrorTitle = lipgloss.NewStyle().
		Foreground(p.textStrong).
		Background(p.errorSurface).
		Bold(true).
		Padding(0, 1)

	s.ErrorBody = lipgloss.NewStyle().
		Foreground(p.textBase).
		MarginLeft(2)

	s.LogItem.Normal.Base = lipgloss.NewStyle().
		Border(lipgloss.Border{
			Left: " ",
		}, false, false, false, true).
		PaddingLeft(1)

	s.LogItem.Active.Base = s.LogItem.Normal.Base.
		Border(lipgloss.Border{
			Left: "┃",
		}, false, false, false, true).
		BorderForeground(p.accent)

	s.LogItem.Active.Hash = lipgloss.NewStyle().
		Bold(true).
		Foreground(p.highlight)

	s.LogItem.Normal.Title = lipgloss.NewStyle().
		Foreground(p.textTitle)

	s.LogItem.Active.Title = lipgloss.NewStyle().
		Foreground(p.highlight).
		Bold(true)

	s.LogItem.Normal.Desc = lipgloss.NewStyle().
		Foreground(p.textCool)

	s.LogItem.Active.Desc = lipgloss.NewStyle().
		Foreground(p.textTitleActive)

	s.LogItem.Active.Keyword = s.LogItem.Active.Desc.
		Foreground(p.highlightMuted)

	s.LogItem.Normal.Hash = lipgloss.NewStyle().
		Foreground(p.hash)

	s.LogItem.Active.Hash = lipgloss.NewStyle().
		Foreground(p.highlight)

	s.Log.Commit = lipgloss.NewStyle().
		Margin(0, 2)

	s.Log.CommitHash = lipgloss.NewStyle().
		Foreground(p.hash).
		Bold(true)

	s.Log.CommitBody = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2)

	s.Log.CommitStatsAdd = lipgloss.NewStyle().
		Foreground(p.success).
		Bold(true)

	s.Log.CommitStatsDel = lipgloss.NewStyle().
		Foreground(p.danger).
		Bold(true)

	s.Log.Paginator = lipgloss.NewStyle().
		Margin(0).
		Align(lipgloss.Center)

	s.Ref.Normal.Item = lipgloss.NewStyle()

	s.Ref.ItemSelector = lipgloss.NewStyle().
		Foreground(p.accent).
		SetString("> ")

	s.Ref.Active.Item = lipgloss.NewStyle().
		Foreground(p.highlightMuted)

	s.Ref.Normal.Base = lipgloss.NewStyle()

	s.Ref.Active.Base = lipgloss.NewStyle()

	s.Ref.Normal.ItemTag = lipgloss.NewStyle().
		Foreground(p.info)

	s.Ref.Active.ItemTag = lipgloss.NewStyle().
		Bold(true).
		Foreground(p.highlight)

	s.Ref.Active.Item = lipgloss.NewStyle().
		Bold(true).
		Foreground(p.highlight)

	s.Ref.Normal.ItemDesc = lipgloss.NewStyle().
		Faint(true)

	s.Ref.Active.ItemDesc = lipgloss.NewStyle().
		Foreground(p.highlight).
		Faint(true)

	s.Ref.Normal.ItemHash = lipgloss.NewStyle().
		Foreground(p.hash).
		Bold(true)

	s.Ref.Active.ItemHash = lipgloss.NewStyle().
		Foreground(p.highlight).
		Bold(true)

	s.Ref.Paginator = s.Log.Paginator

	s.Ref.Selector = lipgloss.NewStyle()

	s.Tree.Selector = s.Tree.Normal.FileName.
		Width(1).
		Foreground(p.accent)

	s.Tree.Normal.FileName = lipgloss.NewStyle().
		MarginLeft(1)

	s.Tree.Active.FileName = s.Tree.Normal.FileName.
		Bold(true).
		Foreground(p.highlight)

	s.Tree.Normal.FileDir = lipgloss.NewStyle().
		Foreground(p.info)

	s.Tree.Active.FileDir = lipgloss.NewStyle().
		Foreground(p.highlight)

	s.Tree.Normal.FileMode = s.Tree.Active.FileName.
		Width(10).
		Foreground(p.textMuted)

	s.Tree.Active.FileMode = s.Tree.Normal.FileMode.
		Foreground(p.highlightMuted)

	s.Tree.Normal.FileSize = s.Tree.Normal.FileName.
		Foreground(p.textMuted)

	s.Tree.Active.FileSize = s.Tree.Normal.FileName.
		Foreground(p.highlightMuted)

	s.Tree.FileContent = lipgloss.NewStyle()

	s.Tree.Paginator = s.Log.Paginator

	s.Tree.Blame.Hash = lipgloss.NewStyle().
		Foreground(p.hash).
		Bold(true)

	s.Tree.Blame.Message = lipgloss.NewStyle()

	s.Tree.Blame.Who = lipgloss.NewStyle().
		Faint(true)

	s.Spinner = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2).
		Foreground(p.accentMuted)

	s.SpinnerContainer = lipgloss.NewStyle()

	s.NoContent = lipgloss.NewStyle().
		MarginTop(1).
		MarginLeft(2).
		Foreground(p.textDim)

	s.StatusBar = lipgloss.NewStyle().
		Height(1)

	s.StatusBarKey = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(p.accentWarm).
		Foreground(p.bgMain)

	s.StatusBarValue = lipgloss.NewStyle().
		Padding(0, 1).
		Background(p.bgInactive).
		Foreground(p.textMuted)

	s.StatusBarInfo = lipgloss.NewStyle().
		Padding(0, 1).
		Background(p.accent).
		Foreground(p.bgMain)

	s.StatusBarBranch = lipgloss.NewStyle().
		Padding(0, 1).
		Background(p.accentSurface).
		Foreground(p.textStrong)

	s.StatusBarHelp = lipgloss.NewStyle().
		Padding(0, 1).
		Background(p.bgDim).
		Foreground(p.textMuted)

	s.Tabs = lipgloss.NewStyle().
		Height(1)

	s.TabInactive = lipgloss.NewStyle()

	s.TabActive = lipgloss.NewStyle().
		Underline(true).
		Foreground(p.accent)

	s.TabSeparator = lipgloss.NewStyle().
		SetString("│").
		Padding(0, 1).
		Foreground(p.borderStrong)

	s.Code.LineDigit = lipgloss.NewStyle().Foreground(p.textFaint)

	s.Code.LineBar = lipgloss.NewStyle().Foreground(p.border)

	s.Stash.Normal.Message = lipgloss.NewStyle().MarginLeft(1)

	s.Stash.Active.Message = s.Stash.Normal.Message.Foreground(p.accent)

	s.Stash.Title = lipgloss.NewStyle().
		Foreground(p.hash).
		Bold(true)

	s.Stash.Selector = lipgloss.NewStyle().
		Width(1).
		Foreground(p.accent)

	return s
}
