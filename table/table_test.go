package table

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// https://github.com/calyptia/go-bubble-table/blob/main/table_test.go

func TestFromValues(t *testing.T) {
	input := "foo1,bar1\nfoo2,bar2\nfoo3,bar3"
	table := New(WithColumns([]Column{{Title: "Foo"}, {Title: "Bar"}}))
	table.FromValues(input, ",")

	if len(table.rows) != 3 {
		t.Fatalf("expect table to have 3 rows but it has %d", len(table.rows))
	}

	expect := []Row{
		{"foo1", "bar1"},
		{"foo2", "bar2"},
		{"foo3", "bar3"},
	}
	if !deepEqual(table.rows, expect) {
		t.Fatal("table rows is not equals to the input")
	}
}

func TestFromValuesWithTabSeparator(t *testing.T) {
	input := "foo1.\tbar1\nfoo,bar,baz\tbar,2"
	table := New(WithColumns([]Column{{Title: "Foo"}, {Title: "Bar"}}))
	table.FromValues(input, "\t")

	if len(table.rows) != 2 {
		t.Fatalf("expect table to have 2 rows but it has %d", len(table.rows))
	}

	expect := []Row{
		{"foo1.", "bar1"},
		{"foo,bar,baz", "bar,2"},
	}
	if !deepEqual(table.rows, expect) {
		t.Fatal("table rows is not equals to the input")
	}
}

func deepEqual(a, b []Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i, r := range a {
		for j, f := range r {
			if f != b[i][j] {
				return false
			}
		}
	}
	return true
}

func newTable() Model {
	columns := []Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}
	rows := []Row{
		{"0", "London", "United Kingdom", "9,540,576"},
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
	}
	t := New(
		WithColumns(columns),
		WithRows(rows),
		WithFocused(true),
		WithHeight(3),
	)
	return t
}

func newTableLong(length int) Model {
	columns := []Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}
	rows := make([]Row, 0, length)
	for i := 0; i < length; i++ {
		rows = append(rows, Row{strconv.Itoa(i), "", "", ""})
	}
	t := New(
		WithColumns(columns),
		WithRows(rows),
		WithFocused(true),
		WithHeight(3),
	)
	return t
}

func TestNav(t *testing.T) {

	testcases := []struct {
		name         string
		setCursor    int
		keyMsgs      []tea.KeyMsg
		wantSelected []int
	}{
		{
			name: "down one",
			keyMsgs: []tea.KeyMsg{tea.KeyMsg(
				tea.Key{
					Type: tea.KeyDown,
				},
			)},
			wantSelected: []int{1},
		},
		{
			name: "down three, multi-select up two",
			keyMsgs: []tea.KeyMsg{
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftUp,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftUp,
					}),
			},
			wantSelected: []int{1, 2, 3},
		},
		{
			name: "down three, multi-select to bottom",
			keyMsgs: []tea.KeyMsg{
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftEnd,
					}),
			},
			wantSelected: []int{3, 4, 5, 6},
		},
		{
			name: "multi-select back",
			keyMsgs: []tea.KeyMsg{
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftUp,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftUp,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftUp,
					}),
			},
			wantSelected: []int{2, 3},
		},
		{
			name: "multi-select to top",
			keyMsgs: []tea.KeyMsg{
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyEnd,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyUp,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftHome,
					}),
			},
			wantSelected: []int{0, 1, 2, 3, 4, 5},
		},
		{
			name: "go to top",
			keyMsgs: []tea.KeyMsg{
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyHome,
					}),
			},
			wantSelected: []int{0},
		},
		{
			name: "multi-select to top, then try multi-selecting up",
			keyMsgs: []tea.KeyMsg{
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyDown,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftUp,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftUp,
					}),
				tea.KeyMsg(
					tea.Key{
						Type: tea.KeyShiftUp,
					}),
			},
			wantSelected: []int{0, 1, 2},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			table := newTable()
			for _, msg := range tc.keyMsgs {
				table, _ = table.Update(msg)
			}
			selected := table.SelectedRows()
			if !reflect.DeepEqual(selected, tc.wantSelected) {
				t.Errorf("selected rows %v, want %v", selected, tc.wantSelected)
			}
		})
	}
}

// Pseudo fuzzing
func TestWindowingSetCursor(t *testing.T) {
	t.Parallel()
	table := newTable()
	for i := 0; i < 1000; i++ {
		table.SetCursor(rand.Intn(10))
		if table.end-table.start != table.Height()-1 {
			t.Error("height", table.end, table.start, table.Height())
		}
	}
}

func TestWindowingHeight(t *testing.T) {
	t.Parallel()
	table := newTable()
	for i := 0; i < 1000; i++ {
		// alter the table
		n := rand.Intn(10)
		table.SetHeight(n)

		// check the table
		if table.Cursor() != clamp(table.Cursor(), table.start, table.end) {
			t.Errorf("cursor %d, start %d, end %d", table.Cursor(), table.start, table.end)
		}
		if table.Height() == 0 {
			if table.end-table.start != 0 {
				t.Error("height:", table.Height(), "end:", table.end, "start:", table.start, "n:", n)
			}
			continue
		}
		if table.end-table.start+1 != table.Height() {
			t.Error("height:", table.Height(), "end:", table.end, "start:", table.start, "n:", n)
		}
	}
}

func TestWindowingCursorAndHeight(t *testing.T) {
	t.Parallel()
	tableLength := 101
	table := newTableLong(tableLength)
	for i := 0; i < 1000; i++ {

		preHeight := table.Height()
		preStart := table.start
		preEnd := table.end
		preCursor := table.Cursor()

		r := rand.Intn(tableLength * 2)
		if rand.Intn(2)%2 == 0 {
			r = -r
		}

		if rand.Intn(2)%2 == 0 {
			table.SetHeight(r)
		}
		if rand.Intn(2)%2 == 0 {
			table.SetCursor(rand.Intn(tableLength * 2))
		}

		if i == 17 {
			table.SetCursor(tableLength)
		}

		// check the table
		if table.Cursor() != clamp(table.Cursor(), table.start, table.end) {
			t.Errorf("cursor %d, start %d, end %d", table.Cursor(), table.start, table.end)
		}
		if table.Cursor() >= tableLength {
			t.Error("Cursor table length Error", "height:", table.Height(), "start:", table.start, "end:", table.end, "cursor:", table.Cursor())
		}
		if table.end >= tableLength {
			t.Error("table length error", "height:", table.Height(), "end:", table.end, "start:", table.start, "cursor:", table.Cursor())
		}
		if table.start < 0 {
			t.Error("Start is less than 0", "height:", table.Height(), "end:", table.end, "start:", table.start, "cursor:", table.Cursor())
		}
		if table.Height() == 0 {
			if table.end-table.start != 0 {
				t.Error("Zero height error", "height:", table.Height(), "end:", table.end, "start:", table.start)
			}
			continue
		}
		if table.end-table.start != table.Height()-1 {
			t.Error("Length Error PRE ", "height:", preHeight, "start:", preStart, "end:", preEnd, "cursor:", preCursor)
			t.Error("Length Error POST", "height:", table.Height(), "start:", table.start, "end:", table.end, "cursor:", table.Cursor())
			t.Fatal()
		}
		table.View()
	}
}

func FuzzHeightCursor(f *testing.F) {
	testcases := []int{-2, 3, 0, 4}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	tableLength := 10
	table := newTableLong(tableLength)
	f.Fuzz(func(t *testing.T, n int) {

		preHeight := table.Height()
		preStart := table.start
		preEnd := table.end
		preCursor := table.Cursor()

		if rand.Intn(2)%2 == 0 {
			table.SetHeight(n)
		} else {
			table.SetCursor(n)
		}

		if table.Cursor() != clamp(table.Cursor(), table.start, table.end) {
			t.Errorf("cursor %d, start %d, end %d", table.Cursor(), table.start, table.end)
		}
		if table.Cursor() >= tableLength {
			t.Error("Cursor table length Error", "height:", table.Height(), "start:", table.start, "end:", table.end, "cursor:", table.Cursor())
		}
		if table.end >= tableLength {
			t.Error("table length error", "height:", table.Height(), "end:", table.end, "start:", table.start, "cursor:", table.Cursor())
		}
		if table.start < 0 {
			t.Error("Start is less than 0", "height:", table.Height(), "end:", table.end, "start:", table.start, "cursor:", table.Cursor())
		}
		if table.Height() == 0 {
			if table.end-table.start != 0 {
				t.Error("Zero height error", "height:", table.Height(), "end:", table.end, "start:", table.start)
			}
			return
		}
		if table.end-table.start != table.Height()-1 {
			t.Error("Length Error PRE ", "height:", preHeight, "start:", preStart, "end:", preEnd, "cursor:", preCursor)
			t.Error("Length Error POST", "height:", table.Height(), "start:", table.start, "end:", table.end, "cursor:", table.Cursor())
		}

	})
}
