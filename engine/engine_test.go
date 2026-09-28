package engine

import "testing"

func perft(b *Board, depth int) int {
	if depth == 0 {
		return 1
	}
	n := 0
	for _, m := range b.Moves() {
		n += perft(b.Apply(m), depth-1)
	}
	return n
}

// Move counts from the opening position, as published for English
// draughts; any mistake in move generation shows up here.
func TestPerft(t *testing.T) {
	want := []int{1, 7, 49, 302, 1469, 7361, 36768, 179740}
	for depth, n := range want {
		if testing.Short() && depth > 6 {
			break
		}
		if got := perft(New(), depth); got != n {
			t.Errorf("perft(%d) = %d, want %d", depth, got, n)
		}
	}
}

// place builds a position from square numbers (1-32).
func place(toMove Side, pieces map[int]int8) *Board {
	b := &Board{ToMove: toMove}
	for sq, p := range pieces {
		b.Sq[sq-1] = p
	}
	return b
}

func notations(ms []Move) map[string]bool {
	out := map[string]bool{}
	for _, m := range ms {
		out[m.String()] = true
	}
	return out
}

func TestCaptureIsCompulsoryAndMultiJumpsAreCompleted(t *testing.T) {
	// Red man on 9 can jump 14 (landing 18) and then 23 (landing 27).
	b := place(Red, map[int]int8{9: RedMan, 14: WhiteMan, 23: WhiteMan, 1: RedMan, 32: WhiteMan})
	moves := notations(b.Moves())
	if len(moves) != 1 || !moves["9x18x27"] {
		t.Fatalf("moves = %v, want only the full double jump 9x18x27", moves)
	}
	after := b.Apply(b.Moves()[0])
	if red, white := after.Count(); red != 2 || white != 1 {
		t.Fatalf("after the double jump: red %d, white %d", red, white)
	}
	if _, err := b.Parse("1-5"); err == nil {
		t.Fatal("a quiet move was accepted while a capture was available")
	}
	if m, err := b.Parse("9x27"); err != nil || m.String() != "9x18x27" {
		t.Fatalf("shorthand 9x27 = %v, %v", m, err)
	}
}

func TestCrowningEndsTheMove(t *testing.T) {
	// A Red man on 23 jumps 26 and lands on 30, the crowning row. A king on 30
	// could jump on over 25, but crowning ends the move.
	b := place(Red, map[int]int8{23: RedMan, 26: WhiteMan, 25: WhiteMan, 1: WhiteMan})
	var found bool
	for _, m := range b.Moves() {
		if m.Path[len(m.Path)-1] == 29 { // square 30
			found = true
			if len(m.Captured) != 1 {
				t.Fatalf("move %s continued after crowning", m)
			}
			after := b.Apply(m)
			if after.Sq[29] != RedKing {
				t.Fatalf("man wasn't crowned: %d", after.Sq[29])
			}
		}
	}
	if !found {
		t.Fatalf("expected a jump to 30; moves = %v", notations(b.Moves()))
	}
}

func TestKingsMoveBothWays(t *testing.T) {
	b := place(Red, map[int]int8{14: RedKing, 32: WhiteMan})
	moves := notations(b.Moves())
	for _, want := range []string{"14-9", "14-10", "14-17", "14-18"} {
		if !moves[want] {
			t.Errorf("king on 14 can't play %s: %v", want, moves)
		}
	}
}

func TestCoordsRoundTrip(t *testing.T) {
	for sq := 0; sq < 32; sq++ {
		r, c := Coords(sq)
		if SquareAt(r, c) != sq {
			t.Fatalf("square %d -> (%d,%d) -> %d", sq+1, r, c, SquareAt(r, c))
		}
	}
}
