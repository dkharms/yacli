package yacli

type gtype int

const (
	groupDefault gtype = iota
	groupMutex
	groupTogether
)

type group struct {
	id    string
	ttype gtype
	met   int
	flags []*flag
}

func (g *group) add(f *flag) {
	g.flags = append(g.flags, f)
}

type flaggroup map[string]*group

func (fg flaggroup) new(t gtype) *group {
	g := &group{id: uniqId(), ttype: t}
	fg[g.id] = g

	return g
}

func (fg flaggroup) add(id string, f *flag) {
	fg[id].flags = append(fg[id].flags, f)
}

func (fg flaggroup) get(id string) *group {
	return fg[id]
}
