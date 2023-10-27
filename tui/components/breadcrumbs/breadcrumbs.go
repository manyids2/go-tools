package breadcrumbs

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Breadcrumbs struct {
	*tview.Box
	Crumbs        []string
	currentOption int
}

func (r *Breadcrumbs) Draw(screen tcell.Screen) {
	r.Box.DrawForSubclass(screen, r)
	x, y, width, _ := r.GetInnerRect()

	line := ""
	separator := "\uf054" // Checked.
	for idx, crumb := range r.Crumbs {
		color := "white"
		if r.currentOption == idx {
			color = "red"
		}
		line += fmt.Sprintf(`%s[%s]  %s[orange]  `, separator, color, crumb)
	}
	tview.Print(screen, line, x, y, width, tview.AlignLeft, tcell.ColorRed)
}

func NewBreadcrumbs(crumbs []string) *Breadcrumbs {
	return &Breadcrumbs{
		Box:    tview.NewBox(),
		Crumbs: crumbs,
	}
}

func (r *Breadcrumbs) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		// Sane keys
		case tcell.KeyUp, tcell.KeyLeft:
			r.currentOption--
			if r.currentOption < 0 {
				r.currentOption = 0
			}
		case tcell.KeyDown, tcell.KeyRight:
			r.currentOption++
			if r.currentOption >= len(r.Crumbs) {
				r.currentOption = len(r.Crumbs) - 1
			}

		// Vim keys
		case tcell.KeyRune:
			switch event.Rune() {
			case 'k', 'h':
				r.currentOption--
				if r.currentOption < 0 {
					r.currentOption = 0
				}
			case 'j', 'l':
				r.currentOption++
				if r.currentOption >= len(r.Crumbs) {
					r.currentOption = len(r.Crumbs) - 1
				}
			}
		}
	},
	)
}
