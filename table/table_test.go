package table

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

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
		WithHeight(7),
	)
	// t.SetStyles(s)
	t.SetCursor(0)
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
