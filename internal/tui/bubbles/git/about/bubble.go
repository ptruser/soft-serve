package about

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/soft-serve/internal/tui/bubbles/git/refs"
	"github.com/charmbracelet/soft-serve/internal/tui/bubbles/git/types"
	vp "github.com/charmbracelet/soft-serve/internal/tui/bubbles/git/viewport"
	"github.com/charmbracelet/soft-serve/internal/tui/style"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/muesli/reflow/wrap"
)

type Bubble struct {
	readmeViewport *vp.ViewportBubble
	repo           types.Repo
	styles         *style.Styles
	height         int
	heightMargin   int
	width          int
	widthMargin    int
	ref            *plumbing.Reference
}

func NewBubble(repo types.Repo, styles *style.Styles, width, wm, height, hm int) *Bubble {
	b := &Bubble{
		readmeViewport: &vp.ViewportBubble{
			Viewport: &viewport.Model{},
		},
		repo:         repo,
		styles:       styles,
		widthMargin:  wm,
		heightMargin: hm,
		ref:          repo.GetHEAD(),
	}
	b.SetSize(width, height)
	return b
}
func (b *Bubble) Init() tea.Cmd {
	return b.setupCmd
}

func (b *Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.SetSize(msg.Width, msg.Height)
		// XXX: if we find that longer readmes take more than a few
		// milliseconds to render we may need to move Glamour rendering into a
		// command.
		md, err := b.glamourize()
		if err != nil {
			return b, nil
		}
		b.readmeViewport.Viewport.SetContent(md)
	case tea.KeyMsg:
		switch msg.String() {
		case "R":
			b.GotoTop()
		}
	case refs.RefMsg:
		b.ref = msg
		return b, b.setupCmd
	}
	rv, cmd := b.readmeViewport.Update(msg)
	b.readmeViewport = rv.(*vp.ViewportBubble)
	cmds = append(cmds, cmd)
	return b, tea.Batch(cmds...)
}

func (b *Bubble) SetSize(w, h int) {
	b.width = w
	b.height = h
	b.readmeViewport.Viewport.Width = w - b.widthMargin
	b.readmeViewport.Viewport.Height = h - b.heightMargin
}

func (b *Bubble) GotoTop() {
	b.readmeViewport.Viewport.GotoTop()
}

func (b *Bubble) View() string {
	return b.readmeViewport.View()
}

func (b *Bubble) Help() []types.HelpEntry {
	return nil
}

func (b *Bubble) glamourize() (string, error) {
	w := b.width - b.widthMargin - b.styles.RepoBody.GetHorizontalFrameSize()
	rm := b.repo.GetReadme()
	if rm == "" {
		return b.styles.AboutNoReadme.Render("No readme found."), nil
	}
	f, err := types.RenderFile(b.repo.GetReadmePath(), rm, w)
	if err != nil {
		return "", err
	}
	// For now, hard-wrap long lines in Glamour that would otherwise break the
	// layout when wrapping. This may be due to #43 in Reflow, which has to do
	// with a bug in the way lines longer than the given width are wrapped.
	//
	//     https://github.com/muesli/reflow/issues/43
	//
	// TODO: solve this upstream in Glamour/Reflow.
	return wrap.String(f, w), nil
}

func (b *Bubble) setupCmd() tea.Msg {
	md, err := b.glamourize()
	if err != nil {
		return types.ErrMsg{Err: err}
	}
	b.readmeViewport.Viewport.SetContent(md)
	b.GotoTop()
	return nil
}
