package stmt

import (
	"reflect"
	"strings"
	"testing"
)

func TestGrab(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		exp rune
	}{
		{"", 0, 0},
		{"a", 0, 'a'},
		{" a", 0, ' '},
		{"a ", 1, ' '},
		{"a", 1, 0},
	}

	for i, test := range tests {
		z := []rune(test.s)
		r := grab(z, test.i, len(z))
		if r != test.exp {
			t.Errorf("test %d expected %c, got: %c", i, test.exp, r)
		}
	}
}

func TestFindSpace(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		exp int
		b   bool
	}{
		{"", 0, 0, false},
		{" ", 0, 0, true},
		{"a", 0, 1, false},
		{"a ", 0, 1, true},
		{" a ", 0, 0, true},
		{"aaa", 0, 3, false},

		{" a ", 1, 2, true},
		{"aaa", 1, 3, false},
		{" aaa", 1, 4, false},
	}

	for i, test := range tests {
		z := []rune(test.s)
		n, b := findSpace(z, test.i, len(z))
		if n != test.exp {
			t.Errorf("test %d expected %d, got: %d", i, test.exp, n)
		}
		if b != test.b {
			t.Errorf("test %d expected %t, got: %t", i, test.b, b)
		}
	}
}

func TestFindNonSpace(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		exp int
		b   bool
	}{
		{"", 0, 0, false},
		{" ", 0, 1, false},
		{"a", 0, 0, true},
		{"a ", 0, 0, true},
		{" a ", 0, 1, true},
		{"    ", 0, 4, false},

		{" a ", 1, 1, true},
		{"aaa", 1, 1, true},
		{" aaa", 1, 1, true},
		{"  aa", 1, 2, true},
		{"    ", 1, 4, false},
	}

	for i, test := range tests {
		z := []rune(test.s)
		n, b := findNonSpace(z, test.i, len(z))
		if n != test.exp {
			t.Errorf("test %d expected %d, got: %d", i, test.exp, n)
		}
		if b != test.b {
			t.Errorf("test %d expected %t, got: %t", i, test.b, b)
		}
	}
}

func TestIsEmptyLine(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		exp bool
	}{
		{"", 0, true},
		{"a", 0, false},
		{" a", 0, false},
		{" a ", 0, false},
		{" \na", 0, false},
		{" \n\ta", 0, false},

		{"a ", 1, true},
		{" a", 1, false},
		{" a ", 1, false},
		{" \na", 1, false},
		{" \n\t ", 1, true},
	}

	for i, test := range tests {
		z := []rune(test.s)
		b := isEmptyLine(z, test.i, len(z))
		if b != test.exp {
			t.Errorf("test %d expected %t, got: %t", i, test.exp, b)
		}
	}
}

func TestStartsWith(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		exp bool
	}{
		{"", 0, false},
		{" ", 0, false},
		{" help", 0, true},
		{"     helpfoo", 0, true},
		{"     help foo", 1, true},
		{"     foo help", 1, false},
	}

	for i, test := range tests {
		z := []rune(test.s)
		b := StartsWith(z, test.i, len(z), "help")
		if b != test.exp {
			t.Errorf("test %d expected %t, got: %t", i, test.exp, b)
		}
	}
}

func TestTrimSplit(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		exp string
	}{
		{"", 0, ""},
		{"   ", 0, ""},
		{" \t\n  ", 0, ""},

		{"a", 0, "a"},
		{"a ", 0, "a"},
		{" a", 0, "a"},
		{" a ", 0, "a"},

		{"a b", 0, "a b"},
		{"a b ", 0, "a b"},
		{" a b", 0, "a b"},
		{" a b ", 0, "a b"},

		{"foo bar", 0, "foo bar"},
		{"foo bar ", 0, "foo bar"},
		{" foo bar", 0, "foo bar"},
		{" foo bar ", 0, "foo bar"},

		{"\\c foo bar z", 1, "c foo bar z"},
		{"\\c foo bar z ", 1, "c foo bar z"},
		{"\\c foo bar z  ", 1, "c foo bar z"},
		{"\\c    foo    bar    z  ", 1, "c foo bar z"},
		{"\\c    pg://blah    bar    z  ", 1, "c pg://blah bar z"},
		{"\\foo    pg://blah    bar    z  ", 1, "foo pg://blah bar z"},
	}

	for i, test := range tests {
		z := []rune(test.s)
		y := trimSplit(z, test.i, len(z))
		exp := strings.Split(test.exp, " ")
		if test.exp == "" {
			if len(y) != 0 {
				t.Errorf("test %d expected result to have length 0, has length: %d", i, len(y))
			}
		} else if !reflect.DeepEqual(y, exp) {
			t.Errorf("test %d expected %v, got: %v", i, exp, y)
		}
	}
}

func TestReadCommand(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		exp string
		r   string
	}{
		{"\\c foo bar z", 0, "c foo bar z", ""},
		{"\\c foo bar z ", 0, "c foo bar z", ""},
		{"\\c foo bar z  ", 0, "c foo bar z", ""},
		{"\\c    foo    bar    z  ", 0, "c foo bar z", ""},
		{"\\c    pg://blah    bar    z  ", 0, "c pg://blah bar z", ""},
		{"\\foo    pg://blah    bar    z  ", 0, "foo pg://blah bar z", ""},
		{"\\p \\p", 0, "p", "\\p"},
		{"\\p foo \\p", 0, "p foo", "\\p"},
		{"\\p foo   \\p bar", 0, "p foo", "\\p bar"},
		{"\\p\\p", 0, "p\\p", ""},
		{"\\p \\r foo", 0, "p", "\\r foo"},
		{"\\print   \\reset    foo", 0, "print", "\\reset    foo"},
		{"\\print   \\reset    foo", 9, "reset foo", ""},
		{"\\print   \\reset    foo  ", 9, "reset foo", ""},
		{"\\print   \\reset    foo  bar  ", 9, "reset foo bar", ""},
	}

	for i, test := range tests {
		z := []rune(test.s)
		a := strings.Split(test.exp, " ")

		cmd, params, pos := readCommand(z, test.i, len(z))
		if cmd != a[0] {
			t.Errorf("test %d expected command to be `%s`, got: `%s`", i, a[0], cmd)
		}
		if !reflect.DeepEqual(params, a[1:]) {
			t.Errorf("test %d expected %v, got: %v", i, a[1:], params)
		}

		m := string(z[pos:])
		if m != test.r {
			t.Errorf("test %d expected remaining to be `%s`, got: `%s`", i, test.r, m)
		}
	}
}

func TestFindWords(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		w   int
		exp string
		c   int
	}{
		{"", 0, 4, "", 0},
		{"  ", 0, 4, "", 0},
		{"  ", 0, 4, "", 0},
		{" select ", 0, 4, " select", 1},
		{" select to ", 0, 4, " select to", 2},
		{" select to ", 1, 4, "select to", 2}, // 5
		{" select   to   ", 0, 4, " select   to", 2},
		{"select into from", 0, 2, "select into", 2},
		{"select into * from", 0, 4, "select into * from", 4},
		{" select  into  *   from  ", 0, 4, " select  into  *   from", 4},
		{"  select\n\n\tb\t\tzfrom j\n\n  ", 1, 2, " select\n\n\tb", 2}, // 10
	}
	for i, test := range tests {
		z := []rune(test.s)

		end, c := findEndOfWords(z, test.i, len(z), test.w)
		s := string(z[test.i:end])
		if s != test.exp {
			t.Errorf("test %d expected `%s`, got: `%s`", i, test.exp, s)
		}

		if c != test.c {
			t.Errorf("test %d expected word count %d, got: %d", i, test.c, c)
		}
	}
}

func TestFindPrefix(t *testing.T) {
	tests := []struct {
		s   string
		i   int
		w   int
		exp string
	}{
		{"", 0, 4, ""},
		{"  ", 0, 4, ""},
		{"  ", 0, 4, ""},
		{" select ", 0, 4, "SELECT"},
		{" select to ", 0, 4, "SELECT TO"},
		{" select to ", 1, 4, "SELECT TO"}, // 5
		{" select   to   ", 0, 4, "SELECT TO"},
		{"select into from", 0, 2, "SELECT INTO"},
		{"select into * from", 0, 4, "SELECT INTO"},
		{" select  into  *   from  ", 0, 4, "SELECT INTO"},
		{"  select\n\n\tb\t\tzfrom j\n\n  ", 1, 2, "SELECT B"}, // 10
	}
	for i, test := range tests {
		z := []rune(test.s)
		p := findPrefix(z, test.i, len(z), test.w)
		if p != test.exp {
			t.Errorf("test %d expected `%s`, got: `%s`", i, test.exp, p)
		}
	}
}
