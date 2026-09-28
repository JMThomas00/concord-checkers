package game

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/JMThomas00/Concord/sdk/plugin"
	"github.com/JMThomas00/Concord/sdk/plugintest"
	"github.com/JMThomas00/Concord/sdk/table"
	"github.com/JMThomas00/Concord/sdk/wire"
	"github.com/JMThomas00/concord-checkers/engine"
	"github.com/google/uuid"
)

func TestGameAdapter(t *testing.T) {
	g := Rules.New(nil).(*Game)
	if g.Turn() != 0 {
		t.Fatalf("Red moves first; Turn() = %d", g.Turn())
	}
	if err := g.Play("11-15"); err != nil {
		t.Fatal(err)
	}
	if g.Turn() != 1 || g.Last.String() != "11-15" {
		t.Fatalf("after 11-15: turn %d, last %s", g.Turn(), g.Last)
	}
	if err := g.Play("11-15"); err == nil {
		t.Fatal("White played Red's move")
	}

	// White has no pieces left: Red wins.
	g.Board = &engine.Board{ToMove: engine.White}
	g.Board.Sq[5] = engine.RedMan
	if o := g.Outcome(); !o.Over || o.Winner != 0 || o.Reason != "all pieces captured" {
		t.Fatalf("outcome = %+v", o)
	}
}

func TestComputerPlaysLegalMovesAndTakesFreePieces(t *testing.T) {
	for level := 1; level <= 3; level++ {
		g := Rules.New(nil)
		began := time.Now()
		for i := 0; i < 6 && !g.Outcome().Over; i++ {
			mv := Rules.AI(g, level)
			if err := g.Play(mv); err != nil {
				t.Fatalf("level %d suggested an illegal move %q: %v", level, mv, err)
			}
		}
		if took := time.Since(began); took > 20*time.Second {
			t.Errorf("level %d took %v for 6 moves", level, took)
		}
	}

	// Red on 10 can take an undefended White man on 15 (landing on 19).
	g := &Game{Board: &engine.Board{ToMove: engine.Red}}
	g.Board.Sq[9] = engine.RedMan    // 10
	g.Board.Sq[14] = engine.WhiteMan // 15
	g.Board.Sq[31] = engine.WhiteMan // 32
	g.Board.Sq[0] = engine.RedMan    // 1
	if mv := Rules.AI(g, 2); mv != "10x19" {
		t.Fatalf("normal computer played %s instead of the free capture 10x19", mv)
	}
}

// Two players sit down and play the opening moves with the keyboard.
func TestTwoPlayersInAChannel(t *testing.T) {
	srv := plugintest.NewServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go plugin.Run(ctx, srv.Config(), table.New(Rules).Handler())
	srv.WaitReady()
	ch := uuid.New()
	srv.Channel(wire.Channel{ID: ch, Name: "checkers"})

	red := srv.Enter(ch, "alice", 60, 24)
	white := srv.Enter(ch, "bob", 60, 24)
	srv.FrameContaining(red, "Tab: sit down")
	srv.FrameContaining(white, "Tab: sit down")
	for _, v := range []*plugintest.Viewer{red, white} { // Sit (first open seat)
		srv.Key(v, "tab")
		srv.Key(v, "enter")
	}
	srv.FrameContaining(red, "your move")

	// Red's cursor starts on square 11; 15 is up and to the right on Red's
	// (flipped) board.
	srv.Key(red, "enter")
	srv.FrameContaining(red, "choose where it goes")
	srv.Key(red, "up")
	srv.Key(red, "right")
	srv.Key(red, "enter")
	frame := srv.FrameContaining(white, "your move")
	if strings.Count(frame, "●") != 24 {
		t.Fatalf("expected 24 pieces on White's board:\n%s", frame)
	}

	// White answers 22-18; its cursor starts on 22, and 18 is up and right.
	srv.Key(white, "enter")
	srv.Key(white, "up")
	srv.Key(white, "right")
	srv.Key(white, "enter")
	srv.FrameContaining(red, "your move")
}
