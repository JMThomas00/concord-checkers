// Package game plugs the checkers engine into the Concord table kit: the
// rules adapter, the computer player, and the board people play on.
package game

import (
	"fmt"

	"github.com/JMThomas00/Concord/sdk/table"
	"github.com/JMThomas00/concord-checkers/engine"
	tea "github.com/charmbracelet/bubbletea"
)

// Rules is checkers for the table kit.
var Rules = table.Rules{
	Name:      "Checkers",
	SeatNames: []string{"Red", "White"},
	New:       func(map[string]string) table.Game { return &Game{Board: engine.New()} },
	NewBoard:  func(s *table.Seat) tea.Model { return newBoard(s) },
	AI: func(g table.Game, level int) string {
		return engine.Best(g.(*Game).Board, level).String()
	},
}

// Game adapts an engine.Board to table.Game.
type Game struct {
	Board *engine.Board
	Last  engine.Move // the move just played (for highlighting)
}

func (g *Game) Turn() int {
	if g.Outcome().Over {
		return -1
	}
	return int(g.Board.ToMove)
}

func (g *Game) Play(move string) error {
	if g.Outcome().Over {
		return fmt.Errorf("the game is over")
	}
	m, err := g.Board.Parse(move)
	if err != nil {
		return err
	}
	g.Board, g.Last = g.Board.Apply(m), m
	return nil
}

func (g *Game) Outcome() table.Outcome {
	if g.Board.Quiet >= 80 {
		return table.Outcome{Over: true, Winner: -1, Reason: "40 moves each without a capture or a man moving"}
	}
	if len(g.Board.Moves()) > 0 {
		return table.Outcome{}
	}
	loser := g.Board.ToMove
	red, white := g.Board.Count()
	reason := "no moves left"
	if (loser == engine.Red && red == 0) || (loser == engine.White && white == 0) {
		reason = "all pieces captured"
	}
	return table.Outcome{Over: true, Winner: int(1 - loser), Reason: reason}
}
