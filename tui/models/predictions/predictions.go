package predictions

import (
	"fmt"
	"log"
	"os"
)

// Group data
type Group struct {
	Name       string
	Slidenames []string
}

// Predictions data
type Predictions struct {
	Datadir string
	Groups  map[string]Group
	Loaded  chan bool
}

func NewPredictions(datadir string) *Predictions {
	m := Predictions{Datadir: datadir, Loaded: make(chan bool)}
	return &m
}

func (m *Predictions) String() string {
	base := fmt.Sprintf(
		`
Output Dir: %s
    Groups: %d
		`,
		m.Datadir, len(m.Groups),
	)
	for k, g := range m.Groups {
		base += fmt.Sprintf("\n  %s: %d", k, len(g.Slidenames))
	}
	return base
}

func (m *Predictions) SetGroups() {
	// Allocate
	m.Groups = make(map[string]Group)

	// Iterate over directories
	entries, err := os.ReadDir(m.Datadir)
	if err != nil {
		log.Println("Could not read datadir: ", m.Datadir, err)
		return
	}
	for _, e := range entries {
		if e.IsDir() {
			// At least record the name
			group := Group{Name: e.Name()}

			// Get slides if readable
			entries, err := os.ReadDir(fmt.Sprintf("%s/%s", m.Datadir, e.Name()))
			if err != nil {
				log.Println("Could not read group: ", group.Name, err)
			} else {
				for _, s := range entries {
					if s.IsDir() {
						group.Slidenames = append(group.Slidenames, s.Name())
					}
				}
			}

			// Update global state
			m.Groups[group.Name] = group
		}
	}

	// Inform that load is finished
	m.Loaded <- true
}
