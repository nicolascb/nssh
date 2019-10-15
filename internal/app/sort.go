package app

import (
	"reflect"
	"sort"
)

func Prop(field string, asc bool) func(p1, p2 SSHHost) bool {
	return func(p1, p2 SSHHost) bool {

		v1 := reflect.Indirect(reflect.ValueOf(p1)).FieldByName(field)
		v2 := reflect.Indirect(reflect.ValueOf(p2)).FieldByName(field)

		ret := false

		switch v1.Kind() {
		case reflect.Int64:
			ret = int64(v1.Int()) < int64(v2.Int())
		case reflect.Float64:
			ret = float64(v1.Float()) < float64(v2.Float())
		case reflect.String:
			ret = string(v1.String()) < string(v2.String())
		}

		if asc {
			return ret
		}
		return !ret
	}
}

type By func(p1, p2 SSHHost) bool

func (by By) Sort(entries []SSHHost) {
	ps := &entriesSort{
		entries: entries,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

type entriesSort struct {
	entries []SSHHost
	by      func(p1, p2 SSHHost) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *entriesSort) Len() int { return len(s.entries) }

// Swap is part of sort.Interface.
func (s *entriesSort) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *entriesSort) Less(i, j int) bool {
	return s.by(s.entries[i], s.entries[j])
}
