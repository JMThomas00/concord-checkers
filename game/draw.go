package game

import (
	"math/rand/v2"
	"time"

	"github.com/JMThomas00/Concord/sdk/arcade"
	"github.com/JMThomas00/concord-checkers/engine"
)

// Drawing on the arcade canvas. A square is 5 columns by 2 rows, holding a
// 5 x 4 pixel piece (the board is 40 x 16), or 5 x 1 holding a flat puck
// when there's less room (40 x 8).

// look is how to draw one board.
type look struct {
	set     *pieceSet
	style   *boardStyle
	flip    bool         // Red at the bottom (draw rows and columns reversed)
	cursor  [2]int       // screen row, col of the cursor; -1 for none
	picked  map[int]bool // squares of the move being made
	targets map[int]bool // where it can land next
	last    map[int]bool // the last move's squares
	ghost   bool         // the pieces are a locked silhouette
	grey    bool         // the board is locked: grey squares
}

const sqW = 5

// boardFits picks the square height for w x h cells: 2 (pictures) or 1
// (pucks), and whether there's room for the style's frame too.
func boardFits(w, h int, s *boardStyle) (sqH int, frame bool) {
	sqH = 1
	if w >= 8*sqW && h >= 16 {
		sqH = 2
	}
	frame = s.edge != "" && w >= 8*sqW+2 && h >= 8*sqH+2
	return sqH, frame
}

// drawBoard draws the position centred in w x h cells at (x, y).
func drawBoard(c *arcade.Canvas, x, y, w, h int, sq *[32]int8, l look) {
	sqH, frame := boardFits(w, h, l.style)
	bw, bh := 8*sqW, 8*sqH
	x += max(0, (w-bw)/2)
	y += max(0, (h-bh)/2)
	if frame {
		role := l.style.edge
		if l.grey {
			role = "ghost"
		}
		c.Box(x-1, y-1, bw+2, bh+2, role, "", "")
	}
	for r := 0; r < 8; r++ {
		for col := 0; col < 8; col++ {
			br, bc := r, col
			if l.flip {
				br, bc = 7-r, 7-col
			}
			sx, sy := x+col*sqW, y+r*sqH
			dark := (br+bc)%2 == 1
			bg := l.style.light
			if dark {
				bg = l.style.dark
			}
			if l.grey {
				bg = "line"
				if dark {
					bg = "bg"
				}
			}
			n := engine.SquareAt(br, bc)
			switch {
			case l.cursor[0] == r && l.cursor[1] == col:
				bg = "yellowB"
			case n >= 0 && l.picked[n]:
				bg = "greenB"
			case n >= 0 && l.targets[n]:
				bg = "cyanB"
			case n >= 0 && l.last[n]:
				bg = "line"
			}
			c.Shade(sx, sy, sqW, sqH, bg)
			if n < 0 {
				continue
			}
			p := sq[n]
			if p == engine.Empty {
				if l.targets[n] {
					c.Text(sx+2, sy+sqH-1, "•", "cyan", bg, true)
				}
				continue
			}
			side, king := int(engine.Owner(p)), engine.IsKing(p)
			if sqH == 2 {
				rows, pal := l.set.sprite(side, king, l.ghost)
				c.Sprite(sx, 2*sy, rows, pal, 1)
			} else {
				rows, pal := l.set.puck(side, king, l.ghost)
				c.Sprite(sx, 2*sy, rows, pal, 1)
			}
		}
	}
}

// drawPieces draws a set's men and kings, Red above White, centred in w x h.
func drawPieces(c *arcade.Canvas, s *pieceSet, x, y, w, h int, ghost bool) {
	x += max(0, (w-11)/2)
	rows := 2
	if h >= 5 {
		y += (h - 5) / 2
	} else {
		rows = 1
	}
	for side := 0; side < rows; side++ {
		for k, king := range []bool{false, true} {
			sp, pal := s.sprite(side, king, ghost)
			c.Sprite(x+k*6, 2*(y+side*3), sp, pal, 1)
		}
	}
}

// samplePosition is a game in progress, for previews.
var samplePosition = func() [32]int8 {
	var sq [32]int8
	for n, p := range map[int]int8{0: engine.RedKing, 5: engine.RedMan, 9: engine.RedMan, 10: engine.RedMan, 14: engine.RedMan,
		17: engine.WhiteMan, 21: engine.WhiteMan, 22: engine.WhiteMan, 26: engine.WhiteMan, 30: engine.WhiteKing} {
		sq[n] = p
	}
	return sq
}()

// preview draws an unlockable for the arcade: a piece set as its men and
// kings on cards and the menu, or in a sample game on the player's board
// when there's room; a board as a sample game in the player's pieces, or a
// tiny pixel board on a card.
func preview(c *arcade.Canvas, id string, sel map[string]string, x, y, w, h int, ghost bool) {
	set, st := pieces(sel[kindPieces]), style(sel[kindBoard])
	isBoard := len(id) > len(boardPrefix) && id[:len(boardPrefix)] == boardPrefix
	if isBoard {
		st = style(id)
	} else {
		set = pieces(id)
	}
	if h >= 8 { // grey out only what's locked
		drawBoard(c, x, y, w, h, &samplePosition, look{set: set, style: st, flip: true, cursor: [2]int{-1, -1}, ghost: ghost && !isBoard, grey: ghost && isBoard})
		return
	}
	if !isBoard {
		drawPieces(c, set, x, y, w, h, ghost)
		return
	}
	// A tiny board: each square 2 x 1 pixels.
	if w < 16 {
		return
	}
	x += max(0, (w-16)/2)
	y += max(0, (h-4)/2)
	for r := 0; r < 8; r++ {
		for col := 0; col < 8; col++ {
			role := st.light
			if (r+col)%2 == 1 {
				role = st.dark
			}
			if ghost {
				role = "bg"
				if (r+col)%2 == 1 {
					role = "line"
				}
			}
			c.Px(x+2*col, 2*y+r, role)
			c.Px(x+2*col+1, 2*y+r, role)
		}
	}
}

// ── Callouts ───────────────────────────────────────────────────────────────

// calloutFor is what a move earns: a combo for a multi-jump, KING ME! for a
// crowning, or "".
func calloutFor(captured int, crowned bool) string {
	switch {
	case captured >= 4:
		return "QUAD!"
	case captured == 3:
		return "TRIPLE!"
	case captured == 2:
		return "DOUBLE!"
	case crowned:
		return "KING ME!"
	}
	return ""
}

// calloutLasts is how long a callout shows.
const calloutLasts = 1400 * time.Millisecond

// drawCallout writes text in the logo font, centred on (cx, row), on a
// plain backing so it reads over the board.
func drawCallout(c *arcade.Canvas, text string, cx, row int) {
	w := arcade.LogoWidth(text, 1)
	x := cx - w/2
	c.Shade(x-1, row-1, w+3, arcade.LogoHeight(1)+2, "bg")
	c.Logo(text, x, 2*row, 1, []string{"yellow", "yellow", "yellow", "orange", "orange", "orange", "orange"})
}

// ── Attract mode ───────────────────────────────────────────────────────────

// The computer plays itself on a small board: an opening skirmish of
// demoPlies moves, a pause, then another game in the next set on the next
// board. Which game comes first changes every hour.
const (
	demoStep  = 3 // ticks a move
	demoPlies = 30
	demoPause = 6 // moves' worth of ticks to hold the end
)

type demoFrame struct {
	board   *engine.Board
	last    engine.Move
	callout string
}

var demoCache struct {
	n      int
	frames []demoFrame
}

// demoGame is game n's positions, the same every time it's asked for.
func demoGame(n int) []demoFrame {
	if demoCache.frames != nil && demoCache.n == n {
		return demoCache.frames
	}
	rnd := rand.New(rand.NewPCG(uint64(n)+1, 0xc4ec))
	b := engine.New()
	frames := []demoFrame{{board: b}}
	for len(frames) <= demoPlies {
		moves := b.Moves()
		if len(moves) == 0 {
			break
		}
		// Most captures first (they're compulsory anyway), else crowning,
		// else anything.
		best := moves[rnd.IntN(len(moves))]
		for _, m := range moves {
			if len(m.Captured) > len(best.Captured) {
				best = m
			}
		}
		nb := b.Apply(best)
		crowned := !engine.IsKing(b.Sq[best.Path[0]]) && engine.IsKing(nb.Sq[best.Path[len(best.Path)-1]])
		frames = append(frames, demoFrame{board: nb, last: best, callout: calloutFor(len(best.Captured), crowned)})
		b = nb
	}
	demoCache.n, demoCache.frames = n, frames
	return frames
}

func attract(c *arcade.Canvas, x, y, w, h, frame int) {
	step := frame / demoStep
	n := int(time.Now().Unix() / 3600)
	for i := 0; i < 1000; i++ {
		length := len(demoGame(n)) + demoPause
		if step < length {
			break
		}
		step -= length
		n++
	}
	frames := demoGame(n)
	f := frames[min(step, len(frames)-1)]
	set := &pieceSets[n%len(pieceSets)]
	st := &boardStyles[n%len(boardStyles)]
	bw := 42
	last := map[int]bool{}
	for _, s := range f.last.Path {
		last[s] = true
	}
	drawBoard(c, x, y, bw, h-4, &f.board.Sq, look{set: set, style: st, flip: true, cursor: [2]int{-1, -1}, last: last})
	if f.callout != "" {
		c.CenterIn(x, bw, y+h-3, "★ "+f.callout+" ★", "yellow", "", true)
	}
	px := x + bw + 4
	c.Text(px, y+1, "NOW SHOWING", "pink", "", true)
	c.Text(px, y+3, set.Name, set.Tier.Role(), "", true)
	c.Text(px, y+4, "on "+st.Name, "dim", "", false)
	if set.Tier != arcade.Starter {
		c.Text(px, y+6, set.Tier.Stars()+" "+set.Tier.Name(), set.Tier.Role(), "", false)
	}
}
