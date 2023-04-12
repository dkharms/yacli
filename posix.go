package yacli

import (
	"fmt"
	"strings"
)

func (c *command) Usage() string {
	var s strings.Builder

	s.WriteString(c.Name())

	if len(c.cs) > 0 {
		s.WriteString(" [")
		var i int
		for name := range c.cs {
			s.WriteString(fmt.Sprintf(" %s", name))
			if i < len(c.cs)-1 {
				s.WriteString(" |")
			}
			i++
		}
		s.WriteString(" ]")
	}

	var (
		defaultGroup  []*flag
		mutexGroup    []*flag
		togetherGroup []*flag
	)

	for _, g := range c.fg {
		switch g.ttype {
		case groupDefault:
			defaultGroup = append(defaultGroup, g.flags...)
		case groupMutex:
			mutexGroup = append(mutexGroup, g.flags...)
		case groupTogether:
			togetherGroup = append(togetherGroup, g.flags...)
		}
	}

	if len(defaultGroup) > 0 {
		s.WriteString(formatDefaultGroup(defaultGroup...))
	}

	if len(mutexGroup) > 0 {
		s.WriteString(formatMutexGroup(mutexGroup...))
	}

	if len(togetherGroup) > 0 {
		s.WriteString(formatTogetherGroup(togetherGroup...))
	}

	for _, arg := range c.as {
		s.WriteString(fmt.Sprintf(" %s", arg.Name()))
	}

	return s.String()
}

func formatDefaultGroup(flags ...*flag) string {
	var s strings.Builder
	for _, f := range flags {
		s.WriteString(" [")
		s.WriteString(fmt.Sprintf(" -%s", f.Short()))
		s.WriteString(" ]")
	}
	return s.String()
}

func formatMutexGroup(flags ...*flag) string {
	var s strings.Builder
	s.WriteString(" [")
	for i, f := range flags {
		s.WriteString(fmt.Sprintf(" -%s", f.Short()))
		if i < len(flags)-1 {
			s.WriteString(" |")
		}
	}
	s.WriteString(" ]")
	return s.String()
}

func formatTogetherGroup(flags ...*flag) string {
	var s strings.Builder
	s.WriteString(" [")
	for _, f := range flags {
		s.WriteString(fmt.Sprintf(" -%s", f.Short()))
	}
	s.WriteString(" ]")
	return s.String()
}
