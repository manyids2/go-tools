/*
* Minimum standard ui keymaps
*   0.            ? : contextual help
*   1.       q, esc : home / someplace predictable / back
*   2.          c-q : exit
*   3.   tab, s-tab : cycle focus / autocomplete
*   4.        space : select
*   5.        enter : open / enter focus
*   6.    backspace : back
*   7.      c-space : command pallette
*   8.   c-c,     y : copy / stop
*   9.   c-v,     p : paste
*  10.   c-z,     u : undo
*  11. c-s-z,   c-r : redo
*  12.   c-f,     / : find
*  13.     j,  down : down
*  14.     k,    up : up
*  15.     h,  left : left
*  16.     l, right : right
 */
package layout

import (
	"github.com/gdamore/tcell/v2"
	"github.com/manyids2/go-tools/tui/components/breadcrumbs"
	"github.com/manyids2/go-tools/tui/components/filebrowser"
	"github.com/rivo/tview"
)

type UI struct {
	// Subclass from grid
	*tview.Grid

	// Slots
	Status  *breadcrumbs.Breadcrumbs
	Sidebar *filebrowser.Filebrowser
	Content *tview.TextArea

	// Basic info
	Datadir string

	// Views
	State  string
	Layout *tview.Grid
	Views  map[string]*tview.Grid

	// Focused
	FocusedChild int
	Children     []*tview.Box
}

func (p *UI) Focus(delegate func(p tview.Primitive)) {
	if p.FocusedChild >= 0 {
		delegate(p.Children[p.FocusedChild])
	} else {
		p.Box.Focus(delegate)
	}
}

func (p *UI) HasFocus() bool {
	if p.FocusedChild >= 0 {
		return p.Children[p.FocusedChild].HasFocus()
	} else {
		return p.Box.HasFocus()
	}
}

func (p *UI) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return p.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if p.FocusedChild >= 0 {
			if handler := p.Children[p.FocusedChild].InputHandler(); handler != nil {
				handler(event, setFocus)
			}
			return
		} else {
			switch event.Key() {
			// ...handle key events not forwarded to the child primitive...
			case tcell.KeyTab:
				p.FocusedChild = (p.FocusedChild + 1) % len(p.Children)
			}
		}
	})
}

func (r *UI) Draw(screen tcell.Screen) {
	switch r.State {
	case "without-sidebar":
		r.DrawForSubclass(screen, r.Views["without-sidebar"])
	default:
		r.DrawForSubclass(screen, r.Views["with-sidebar"])
	}
}

func NewUI(datadir string) *UI {
	ui := UI{
		Datadir:      datadir,
		Status:       breadcrumbs.NewBreadcrumbs([]string{"hi", "hello"}),
		Sidebar:      filebrowser.NewFilebrowser(datadir),
		Content:      tview.NewTextArea(),
		FocusedChild: 0,
		State:        "with-sidebar",
	}

	// AddItem(p Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool)
	ui.Views = make(map[string]*tview.Grid, 2)

	// Without sidebar
	LayoutWithoutSidebar := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0).
		SetBorders(false)
	LayoutWithoutSidebar.AddItem(ui.Status, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.Content, 1, 0, 1, 1, 0, 0, false)
	ui.Views["without-sidebar"] = LayoutWithoutSidebar

	// With sidebar
	LayoutWithSidebar := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(-1, -3).
		SetBorders(false)
	LayoutWithSidebar.AddItem(ui.Status, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.Sidebar.Tree, 0, 0, 2, 1, 0, 0, false).
		AddItem(ui.Content, 1, 1, 1, 1, 0, 0, false)
	ui.Views["with-sidebar"] = LayoutWithSidebar

	ui.Children = []*tview.Box{
		ui.Sidebar.SetBorder(false),
		ui.Status.SetBorder(false),
		ui.Content.SetBorder(false),
	}

	return &ui
}
