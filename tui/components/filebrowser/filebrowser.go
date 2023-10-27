package filebrowser

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Filebrowser struct {
	*tview.Box
	Datadir string
	Tree    *tview.TreeView
}

func (r *Filebrowser) Draw(screen tcell.Screen) {
	r.Tree.DrawForSubclass(screen, r)
}

func (r *Filebrowser) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.Tree.InputHandler()
}

func NewFilebrowser(datadir string) *Filebrowser {
	root := tview.NewTreeNode(datadir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir())
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, datadir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	tree.SetBorder(true)

	return &Filebrowser{Datadir: datadir, Tree: tree}
}
