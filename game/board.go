package game

import (
	"strings"

	"github.com/JMThomas00/Concord/sdk/table"
	"github.com/JMThomas00/concord-checkers/engine"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Board is what a player sees and moves on. The cursor lives in screen
// coordinates (row 0 at the top as drawn), so arrow keys always move the
// way they look, whichever side the board is drawn from.
type Board struct {
	seat          *table.Seat
	row, col      int   // cursor, screen coordinates
	path          []int // squares chosen so far for the move being made
	err           string
	width, height int
}

func newBoard(s *table.Seat) *Board {
	// Start the cursor on the player's own side (square 11 for Red, 22 for White).
	return &Board{seat: s, row: 5, col: 2}
}

func (b *Board) game() *Game { return b.seat.Game().(*Game) }

// flipped: the board is stored with Red's pieces at the top; flipping it
// puts Red's at the bottom for Red (and spectators), the way a player
// sits at a board. White sees it unflipped, their own pieces at the bottom.
func (b *Board) flipped() bool { return b.seat.Perspective() == int(engine.Red) }

// toBoard converts screen coordinates to board coordinates.
func (b *Board) toBoard(row, col int) (int, int) {
	if b.flipped() {
		return 7 - row, 7 - col
	}
	return row, col
}

func (b *Board) cursorSquare() int { return engine.SquareAt(b.toBoard(b.row, b.col)) }

func (b *Board) Init() tea.Cmd { return nil }

func (b *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width, b.height = msg.Width, msg.Height
	case table.ChangedMsg:
		b.path, b.err = nil, ""
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			b.row = max(0, b.row-1)
		case "down", "j":
			b.row = min(7, b.row+1)
		case "left", "h":
			b.col = max(0, b.col-1)
		case "right", "l":
			b.col = min(7, b.col+1)
		case "enter", " ":
			b.choose()
		case "esc":
			if len(b.path) > 0 {
				b.path, b.err = nil, ""
			} else {
				return b, tea.Quit // hand the keyboard back (in Concord)
			}
		case "q":
			return b, tea.Quit
		}
	}
	return b, nil
}

// candidates are the legal moves continuing the squares chosen so far.
func (b *Board) candidates() []engine.Move {
	var out []engine.Move
	for _, m := range b.game().Board.Moves() {
		if len(m.Path) > len(b.path) && prefix(b.path, m.Path) {
			out = append(out, m)
		}
	}
	return out
}

func prefix(p, of []int) bool {
	for i := range p {
		if p[i] != of[i] {
			return false
		}
	}
	return true
}

// choose acts on the square under the cursor: pick up a piece, step to a
// landing square, or finish the move.
func (b *Board) choose() {
	b.err = ""
	if !b.seat.MyTurn() {
		b.err = "not your turn"
		return
	}
	sq := b.cursorSquare()
	if sq < 0 {
		return
	}
	g := b.game()
	own := engine.Owner(g.Board.Sq[sq]) == g.Board.ToMove
	if own && len(b.path) <= 1 { // (re)select a piece
		b.path = []int{sq}
		if len(b.candidates()) == 0 {
			b.path = nil
			if len(g.Board.Moves()) > 0 && len(g.Board.Moves()[0].Captured) > 0 {
				b.err = "you must capture"
			} else {
				b.err = "that piece can't move"
			}
		}
		return
	}
	if len(b.path) == 0 {
		return
	}
	next := append(append([]int(nil), b.path...), sq)
	var extends, finishes bool
	var done engine.Move
	for _, m := range b.candidates() {
		if !prefix(next, m.Path) {
			continue
		}
		if len(m.Path) == len(next) {
			finishes, done = true, m
		} else {
			extends = true
		}
	}
	switch {
	case finishes && !extends:
		b.path = nil
		if err := b.seat.Play(done.String()); err != nil {
			b.err = err.Error()
		}
	case extends || finishes:
		b.path = next // keep jumping
	default:
		b.err = "can't move there"
	}
}

func (b *Board) View() string {
	g := b.game()
	cellW, cellH := 3, 1
	if b.width >= 8*5 && b.height >= 8*2+2 {
		cellW, cellH = 5, 2
	}
	dark := b.seat.Color("current_line", lipgloss.Color("238"))
	light := b.seat.Color("selection", lipgloss.Color("244"))
	redC := b.seat.Color("red", lipgloss.Color("9"))
	whiteC := b.seat.Color("foreground", lipgloss.Color("15"))
	hint := b.seat.Color("green", lipgloss.Color("10"))
	cursorC := b.seat.Color("cyan", lipgloss.Color("14"))
	lastC := b.seat.Color("yellow", lipgloss.Color("11"))

	targets := map[int]bool{}
	for _, m := range b.candidates() {
		if len(b.path) > 0 {
			targets[m.Path[len(b.path)]] = true
		}
	}
	last := map[int]bool{}
	for _, s := range g.Last.Path {
		last[s] = true
	}

	var lines []string
	for row := 0; row < 8; row++ {
		rowLines := make([]string, cellH)
		for col := 0; col < 8; col++ {
			br, bc := b.toBoard(row, col)
			sq := engine.SquareAt(br, bc)
			bg := light
			if sq >= 0 {
				bg = dark
			}
			style := lipgloss.NewStyle().Background(bg).Width(cellW).Align(lipgloss.Center)
			glyph := ""
			if sq >= 0 {
				switch p := g.Board.Sq[sq]; {
				case p != engine.Empty:
					glyph = "●"
					if engine.IsKing(p) {
						glyph = "♛"
					}
					fg := redC
					if engine.Owner(p) == engine.White {
						fg = whiteC
					}
					style = style.Foreground(fg).Bold(true)
				case targets[sq]:
					glyph = "·"
					style = style.Foreground(hint).Bold(true)
				}
				if len(b.path) > 0 && b.path[len(b.path)-1] == sq {
					style = style.Underline(true)
				}
				if last[sq] && glyph == "" {
					glyph = "∙"
					style = style.Foreground(lastC)
				}
			}
			if row == b.row && col == b.col && b.seat.MyTurn() {
				style = style.Background(cursorC)
			}
			for i := range rowLines {
				content := ""
				if i == (cellH-1)/2 {
					content = glyph
				}
				rowLines[i] += style.Render(content)
			}
		}
		lines = append(lines, rowLines...)
	}

	status := ""
	switch {
	case b.err != "":
		status = b.err
	case b.seat.MyTurn() && len(b.path) > 0:
		status = "choose where it goes (Esc to cancel)"
	case b.seat.MyTurn():
		status = "your move — arrows, Enter to pick a piece"
	case len(g.Last.Path) > 0:
		status = "last move " + g.Last.String()
	}
	return strings.Join(append(lines, status), "\n")
}
