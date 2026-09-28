// Package engine is American checkers (English draughts): an 8x8 board,
// men move diagonally forward, captures are compulsory and a jumping piece
// must keep jumping, men reaching the far row become kings (which ends the
// move), and a player with no legal move loses.
//
// Squares use standard notation: the 32 dark squares are numbered 1-32, row
// by row from the top of the board as drawn from Red's side. Red starts on
// 1-12 and moves first; White starts on 21-32. Moves are written
// "11-15" (a step) or "15x24" / "15x24x31" (jumps, every landing square).
package engine

import (
	"fmt"
	"strconv"
	"strings"
)

// Piece values on a square.
const (
	Empty     = 0
	RedMan    = 1
	RedKing   = 2
	WhiteMan  = 3
	WhiteKing = 4
)

// Side is 0 (Red, moves first) or 1 (White).
type Side int

const (
	Red   Side = 0
	White Side = 1
)

// Board is a position: 32 squares (index 0 = square 1) and who's to move.
type Board struct {
	Sq     [32]int8
	ToMove Side
	// Quiet counts half-moves since the last capture or man move; at 80
	// (40 moves each) the game is drawn.
	Quiet int
}

// Move is a legal move: the squares visited (from, then each landing
// square) and the squares of pieces it captures.
type Move struct {
	Path     []int // square indexes (0-31)
	Captured []int
}

// String is the move in standard notation, e.g. "11-15" or "15x24x31".
func (m Move) String() string {
	sep := "-"
	if len(m.Captured) > 0 {
		sep = "x"
	}
	parts := make([]string, len(m.Path))
	for i, s := range m.Path {
		parts[i] = strconv.Itoa(s + 1)
	}
	return strings.Join(parts, sep)
}

// New is the starting position.
func New() *Board {
	b := &Board{}
	for i := 0; i < 12; i++ {
		b.Sq[i] = RedMan
		b.Sq[31-i] = WhiteMan
	}
	return b
}

// Coords returns a square's row (0 = top, Red's back row) and column (0-7).
func Coords(sq int) (row, col int) {
	row = sq / 4
	col = (sq % 4) * 2
	if row%2 == 0 {
		col++
	}
	return row, col
}

// SquareAt returns the square index at row, col, or -1 for a light square
// or a position off the board.
func SquareAt(row, col int) int {
	if row < 0 || row > 7 || col < 0 || col > 7 || (row+col)%2 == 0 {
		return -1
	}
	return row*4 + col/2
}

// Owner reports which side a piece belongs to (-1 for empty).
func Owner(p int8) Side {
	switch p {
	case RedMan, RedKing:
		return Red
	case WhiteMan, WhiteKing:
		return White
	}
	return -1
}

// IsKing reports whether p is a king.
func IsKing(p int8) bool { return p == RedKing || p == WhiteKing }

// directions a piece on sq may move: Red men go down (+1 row), White men up.
func dirs(p int8) [][2]int {
	switch p {
	case RedMan:
		return [][2]int{{1, -1}, {1, 1}}
	case WhiteMan:
		return [][2]int{{-1, -1}, {-1, 1}}
	case RedKing, WhiteKing:
		return [][2]int{{1, -1}, {1, 1}, {-1, -1}, {-1, 1}}
	}
	return nil
}

func crownRow(s Side) int {
	if s == Red {
		return 7
	}
	return 0
}

// Moves lists every legal move for the side to move. If any capture is
// possible, only captures are legal (and each is played to the end).
func (b *Board) Moves() []Move {
	var jumps []Move
	for sq := 0; sq < 32; sq++ {
		if Owner(b.Sq[sq]) == b.ToMove {
			b.jumpsFrom(sq, b.Sq[sq], []int{sq}, nil, &jumps)
		}
	}
	if len(jumps) > 0 {
		return jumps
	}
	var steps []Move
	for sq := 0; sq < 32; sq++ {
		p := b.Sq[sq]
		if Owner(p) != b.ToMove {
			continue
		}
		r, c := Coords(sq)
		for _, d := range dirs(p) {
			if to := SquareAt(r+d[0], c+d[1]); to >= 0 && b.Sq[to] == Empty {
				steps = append(steps, Move{Path: []int{sq, to}})
			}
		}
	}
	return steps
}

// jumpsFrom extends a capture sequence from the piece's current square
// (the end of path) as far as it goes, collecting finished sequences.
func (b *Board) jumpsFrom(start int, p int8, path, captured []int, out *[]Move) {
	at := path[len(path)-1]
	r, c := Coords(at)
	extended := false
	for _, d := range dirs(p) {
		over := SquareAt(r+d[0], c+d[1])
		to := SquareAt(r+2*d[0], c+2*d[1])
		if over < 0 || to < 0 {
			continue
		}
		if Owner(b.Sq[over]) != 1-b.ToMove || contains(captured, over) {
			continue
		}
		if b.Sq[to] != Empty && to != start { // the moving piece has left start
			continue
		}
		extended = true
		np := append(append([]int(nil), path...), to)
		nc := append(append([]int(nil), captured...), over)
		if !IsKing(p) {
			if tr, _ := Coords(to); tr == crownRow(b.ToMove) {
				*out = append(*out, Move{Path: np, Captured: nc}) // crowning ends the move
				continue
			}
		}
		b.jumpsFrom(start, p, np, nc, out)
	}
	if !extended && len(captured) > 0 {
		*out = append(*out, Move{Path: path, Captured: captured})
	}
}

func contains(xs []int, x int) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// Apply plays m (which must come from Moves) and returns the new position.
func (b *Board) Apply(m Move) *Board {
	nb := *b
	from, to := m.Path[0], m.Path[len(m.Path)-1]
	p := nb.Sq[from]
	nb.Sq[from] = Empty
	for _, c := range m.Captured {
		nb.Sq[c] = Empty
	}
	if !IsKing(p) {
		if r, _ := Coords(to); r == crownRow(b.ToMove) {
			p++ // man → king of the same color
		}
	}
	nb.Sq[to] = p
	if len(m.Captured) > 0 || !IsKing(b.Sq[from]) {
		nb.Quiet = 0
	} else {
		nb.Quiet++
	}
	nb.ToMove = 1 - b.ToMove
	return &nb
}

// Parse finds the legal move written as notation ("11-15", "15x24x31");
// a jump may also be given by just its start and end ("15x31") when that's
// unambiguous.
func (b *Board) Parse(notation string) (Move, error) {
	notation = strings.TrimSpace(strings.ToLower(notation))
	fields := strings.FieldsFunc(notation, func(r rune) bool { return r == '-' || r == 'x' })
	if len(fields) < 2 {
		return Move{}, fmt.Errorf("write moves like 11-15 or 15x24")
	}
	var squares []int
	for _, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil || n < 1 || n > 32 {
			return Move{}, fmt.Errorf("%q isn't a square (1-32)", f)
		}
		squares = append(squares, n-1)
	}
	var match []Move
	for _, m := range b.Moves() {
		if len(squares) == len(m.Path) && equal(squares, m.Path) {
			return m, nil
		}
		if len(squares) == 2 && m.Path[0] == squares[0] && m.Path[len(m.Path)-1] == squares[1] {
			match = append(match, m)
		}
	}
	switch len(match) {
	case 1:
		return match[0], nil
	case 0:
		if len(b.Moves()) > 0 && len(b.Moves()[0].Captured) > 0 {
			return Move{}, fmt.Errorf("you must capture")
		}
		return Move{}, fmt.Errorf("%s isn't a legal move", notation)
	}
	return Move{}, fmt.Errorf("more than one jump goes there; give every square")
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Count returns how many pieces each side has.
func (b *Board) Count() (red, white int) {
	for _, p := range b.Sq {
		switch Owner(p) {
		case Red:
			red++
		case White:
			white++
		}
	}
	return
}
