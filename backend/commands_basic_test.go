package backend

import (
	. "lime/backend/primitives"
	"reflect"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	data := `Hello world
Test
Goodbye world
`
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, data)
	v.EndEdit(e)

	v.Sel().Clear()
	v.Sel().Add(Region{11, 11})
	v.Sel().Add(Region{16, 16})
	v.Sel().Add(Region{30, 30})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if v.Buffer().String() != `Hello worl
Tes
Goodbye worl
` {
		t.Error(v.Buffer().String())
	}
	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "a"})
	if d := v.Buffer().String(); d != "Hello worla\nTesa\nGoodbye worla\n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	v.Settings().Set("translate_tabs_to_spaces", true)
	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "\t"})
	if v.Buffer().String() != "Hello worla \nTesa    \nGoodbye worla   \n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.Buffer().String(); d != "Hello worla\nTesa\nGoodbye worla\n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.Buffer().String(); d != "Hello worl\nTes\nGoodbye worl\n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": "\t"})
	if d := v.Buffer().String(); d != "Hello worl  \nTes \nGoodbye worl    \n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if v.Buffer().String() != "Hello worl\nTes\nGoodbye worl\n" {
		lines := strings.Split(v.Buffer().String(), "\n")
		for _, l := range lines {
			t.Errorf("%d: '%s'", len(l), l)
		}
	}

	v.Buffer().Erase(0, len(v.Buffer().String()))
	v.Buffer().Insert(0, "€þıœəßðĸʒ×ŋµåäö𝄞")
	orig := "€þıœəßðĸʒ×ŋµåäö𝄞"
	if d := v.Buffer().String(); d != orig {
		t.Errorf("%s\n\t%v\n\t%v", d, []byte(d), []byte(orig))
	} else {
		t.Logf("ref %s\n\t%v\n\t%v", d, []byte(d), []byte(orig))
	}
	v.Sel().Clear()
	v.Sel().Add(Region{3, 3})
	v.Sel().Add(Region{6, 6})
	v.Sel().Add(Region{9, 9})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	exp := "€þœəðĸ×ŋµåäö𝄞"
	if d := v.Buffer().String(); d != exp {
		t.Errorf("%s\n\t%v\n\t%v", d, []byte(d), []byte(exp))
	}
}

func TestLeftDelete(t *testing.T) {
	ed := GetEditor()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, "12345678")
	v.EndEdit(e)
	v.Sel().Clear()
	v.Sel().Add(Region{1, 1})
	v.Sel().Add(Region{2, 2})
	v.Sel().Add(Region{3, 3})
	v.Sel().Add(Region{4, 4})
	ed.CommandHandler().RunTextCommand(v, "left_delete", nil)
	if d := v.buffer.String(); d != "5678" {
		t.Error(d)
	}
}

func TestMove(t *testing.T) {
	ed := GetEditor()

	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in      []Region
		by      string
		extend  bool
		forward bool
		exp     []Region
		args    Args
	}

	tests := []Test{
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"characters",
			false,
			true,
			[]Region{{2, 2}, {4, 4}, {7, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"characters",
			false,
			false,
			[]Region{{0, 0}, {2, 2}, {5, 5}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			false,
			true,
			[]Region{{2, 2}, {4, 4}, {7, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			false,
			false,
			[]Region{{0, 0}, {2, 2}, {5, 5}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			true,
			true,
			[]Region{{1, 2}, {3, 4}, {10, 7}},
			nil,
		},
		{
			[]Region{{1, 1}, {3, 3}, {10, 6}},
			"characters",
			true,
			false,
			[]Region{{1, 0}, {3, 2}, {10, 5}},
			nil,
		},
		{
			[]Region{{1, 3}, {3, 5}, {10, 7}},
			"characters",
			true,
			true,
			[]Region{{1, 6}, {10, 8}},
			nil,
		},
		{
			[]Region{{1, 1}},
			"stops",
			true,
			true,
			[]Region{{1, 5}},
			Args{"word_end": true},
		},
		{
			[]Region{{1, 1}},
			"stops",
			false,
			true,
			[]Region{{6, 6}},
			Args{"word_begin": true},
		},
		{
			[]Region{{6, 6}},
			"stops",
			false,
			false,
			[]Region{{0, 0}},
			Args{"word_begin": true},
		},
	}
	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		args := Args{"by": test.by, "extend": test.extend, "forward": test.forward}
		if test.args != nil {
			for k, v := range test.args {
				args[k] = v
			}
		}
		ed.CommandHandler().RunTextCommand(v, "move", args)
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.exp) {
			t.Errorf("Move test %d failed: %v", i, sr)
		}
	}
}

func TestGlueCmds(t *testing.T) {
	ed := GetEditor()
	ch := ed.CommandHandler()
	w := ed.NewWindow()
	v := w.NewFile()
	v.SetScratch(true)
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)
	v.SetScratch(false)
	ch.RunTextCommand(v, "mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "a"})
	ch.RunTextCommand(v, "insert", Args{"characters": "b"})
	ch.RunTextCommand(v, "insert", Args{"characters": "c"})
	ch.RunTextCommand(v, "glue_marked_undo_groups", nil)
	if v.undoStack.position != 1 {
		t.Error(v.undoStack.position)
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "redo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	if v.undoStack.position != 1 {
		t.Error(v.undoStack.position)
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}

	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "a"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "b"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "insert", Args{"characters": "c"})
	ch.RunTextCommand(v, "maybe_mark_undo_groups_for_gluing", nil)
	ch.RunTextCommand(v, "glue_marked_undo_groups", nil)
	if v.undoStack.position != 1 {
		t.Error(v.undoStack.position)
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "undo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\n" {
		t.Error(d)
	}
	ch.RunTextCommand(v, "redo", nil)
	if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
	if v.undoStack.position != 1 {
		t.Error(v.undoStack.position)
	} else if d := v.Buffer().String(); d != "Hello World!\nTest123123\nAbrakadabra\nabc" {
		t.Error(d)
	}
}

func TestInsert(t *testing.T) {
	ed := GetEditor()
	ch := ed.CommandHandler()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	v.EndEdit(e)

	type Test struct {
		in   []Region
		data string
		expd string
		expr []Region
	}

	tests := []Test{
		{
			[]Region{{1, 1}, {3, 3}, {6, 6}},
			"a",
			"Haelalo aWorld!\nTest123123\nAbrakadabra\n",
			[]Region{{2, 2}, {5, 5}, {9, 9}},
		},
		{
			[]Region{{1, 1}, {3, 3}, {6, 9}},
			"a",
			"Haelalo ald!\nTest123123\nAbrakadabra\n",
			[]Region{{2, 2}, {5, 5}, {9, 9}},
		},
		{
			[]Region{{1, 1}, {3, 3}, {6, 9}},
			"€þıœəßðĸʒ×ŋµåäö𝄞",
			"H€þıœəßðĸʒ×ŋµåäö𝄞el€þıœəßðĸʒ×ŋµåäö𝄞lo €þıœəßðĸʒ×ŋµåäö𝄞ld!\nTest123123\nAbrakadabra\n",
			[]Region{{17, 17}, {35, 35}, {54, 54}},
		},
	}
	for i, test := range tests {
		v.Sel().Clear()
		for _, r := range test.in {
			v.Sel().Add(r)
		}
		ed.CommandHandler().RunTextCommand(v, "insert", Args{"characters": test.data})
		if d := v.buffer.String(); d != test.expd {
			t.Errorf("Insert test %d failed: %s", i, d)
		}
		if sr := v.Sel().Regions(); !reflect.DeepEqual(sr, test.expr) {
			t.Errorf("Insert test %d failed: %v", i, sr)
		}
		ch.RunTextCommand(v, "undo", nil)
	}
}

type scfe struct {
	DummyFrontend
	show Region
}

func (f *scfe) VisibleRegion(v *View) Region {
	s := v.Buffer().Line(v.Buffer().TextPoint(3*3, 1))
	e := v.Buffer().Line(v.Buffer().TextPoint(6*3, 1))
	return Region{s.Begin(), e.End()}
}

func (f *scfe) Show(v *View, r Region) {
	f.show = r
}

func TestScrollLines(t *testing.T) {
	var fe scfe
	ed := GetEditor()
	ed.SetFrontend(&fe)
	ch := ed.CommandHandler()
	w := ed.NewWindow()
	v := w.NewFile()
	e := v.BeginEdit()
	for i := 0; i < 10; i++ {
		v.Insert(e, 0, "Hello World!\nTest123123\nAbrakadabra\n")
	}
	v.EndEdit(e)
	ch.RunTextCommand(v, "scroll_lines", Args{"amount": 0})

	if c := v.Buffer().Line(v.Buffer().TextPoint(3*3, 1)); fe.show.Begin() != c.Begin() {
		t.Errorf("Expected %v, but got %v", c, fe.show)
	}

	ch.RunTextCommand(v, "scroll_lines", Args{"amount": 1})
	if c := v.Buffer().Line(v.Buffer().TextPoint(3*3-1, 1)); fe.show.Begin() != c.Begin() {
		t.Errorf("Expected %v, but got %v", c, fe.show)
	}
	t.Log(fe.VisibleRegion(v), v.Buffer().Line(v.Buffer().TextPoint(6*3+1, 1)))
	ch.RunTextCommand(v, "scroll_lines", Args{"amount": -1})
	if c := v.Buffer().Line(v.Buffer().TextPoint(6*3+1, 1)); fe.show.Begin() != c.Begin() {
		t.Errorf("Expected %v, but got %v", c, fe.show)
	}
}
