package engine

import (
	"math/rand"
	"time"
)

// Best picks a move for the side to move at the given level:
//
//	1 (easy)   -- looks two moves ahead, and often plays a random decent move
//	2 (normal) -- looks ~6 moves ahead
//	3 (hard)   -- searches as deep as it can in about two seconds
func Best(b *Board, level int) Move {
	moves := b.Moves()
	if len(moves) == 1 {
		return moves[0]
	}
	switch level {
	case 1:
		scored := rootScores(b, moves, 2)
		// Pick at random among moves within a man's value of the best.
		best := scored[0].score
		var ok []Move
		for _, s := range scored {
			if s.score >= best-manValue {
				ok = append(ok, s.move)
			}
		}
		return ok[rand.Intn(len(ok))]
	case 2:
		return rootScores(b, moves, 6)[0].move
	default:
		deadline := time.Now().Add(2 * time.Second)
		best := rootScores(b, moves, 4)[0].move
		for depth := 5; depth <= 20 && time.Now().Before(deadline); depth++ {
			s := &search{deadline: deadline}
			scored := s.root(b, moves, depth)
			if s.timeout {
				break
			}
			best = scored[0].move
		}
		return best
	}
}

const (
	manValue  = 100
	kingValue = 160
	win       = 100000
)

type scoredMove struct {
	move  Move
	score int
}

func rootScores(b *Board, moves []Move, depth int) []scoredMove {
	s := &search{}
	return s.root(b, moves, depth)
}

type search struct {
	deadline time.Time
	nodes    int
	timeout  bool
}

// root scores every move, best first.
func (s *search) root(b *Board, moves []Move, depth int) []scoredMove {
	out := make([]scoredMove, len(moves))
	for i, m := range moves {
		out[i] = scoredMove{m, -s.negamax(b.Apply(m), depth-1, -win-1, win+1)}
	}
	// Stable sort, best first (insertion sort: move lists are short).
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].score > out[j-1].score; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// negamax scores b for the side to move with alpha-beta pruning. Captures
// extend the search, so it doesn't stop in the middle of an exchange.
func (s *search) negamax(b *Board, depth, alpha, beta int) int {
	s.nodes++
	if s.nodes&1023 == 0 && !s.deadline.IsZero() && time.Now().After(s.deadline) {
		s.timeout = true
	}
	if s.timeout {
		return 0
	}
	moves := b.Moves()
	if len(moves) == 0 {
		return -win + (20 - depth) // lose; prefer losing later
	}
	if b.Quiet >= 80 {
		return 0
	}
	if depth <= 0 && len(moves[0].Captured) == 0 {
		return evaluate(b)
	}
	for _, m := range moves {
		score := -s.negamax(b.Apply(m), depth-1, -beta, -alpha)
		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}
	return alpha
}

// evaluate scores a quiet position for the side to move: material, with a
// little for advancing men and keeping the back row (which stops enemy
// men crowning).
func evaluate(b *Board) int {
	score := 0
	for sq, p := range b.Sq {
		if p == Empty {
			continue
		}
		v := manValue
		if IsKing(p) {
			v = kingValue
		} else {
			row, col := Coords(sq)
			advance := row // Red advances down the board
			back := row == 0
			if Owner(p) == White {
				advance, back = 7-row, row == 7
			}
			v += advance * 3
			if back {
				v += 8
			}
			if col >= 2 && col <= 5 && row >= 2 && row <= 5 {
				v += 4 // center
			}
		}
		if Owner(p) == b.ToMove {
			score += v
		} else {
			score -= v
		}
	}
	return score
}
