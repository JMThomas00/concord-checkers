package game

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/JMThomas00/Concord/sdk/arcade"
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

// Two players sit down at the plain board (a pane under 64 x 24) and play
// the opening moves with the keyboard.
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
	for _, v := range []*plugintest.Viewer{red, white} { // too small for the arcade: the plain door
		srv.FrameContaining(v, "Enter to play")
		srv.Key(v, "enter")
		srv.FrameContaining(v, "M: sit down")
	}
	for _, v := range []*plugintest.Viewer{red, white} { // Sit (first open seat)
		srv.Key(v, "m")
		srv.Key(v, "enter")
	}
	srv.FrameContaining(red, "your move")

	// Red's cursor starts on square 11; 15 is up and to the right on Red's
	// (flipped) board.
	// Esc is Concord's until a piece is picked up; then it puts it back.
	if srv.Key(red, "esc") {
		t.Fatal("Esc was claimed with no move started")
	}
	srv.Key(red, "enter")
	srv.FrameContaining(red, "choose where it goes")
	if !srv.Key(red, "esc") {
		t.Fatal("Esc wasn't claimed with a piece picked up")
	}
	srv.FrameContaining(red, "your move")
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

func startServer(t *testing.T) (*plugintest.Server, uuid.UUID) {
	t.Helper()
	srv := plugintest.NewServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { plugin.Run(ctx, srv.Config(), table.New(Rules).Handler()); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	srv.WaitReady()
	ch := uuid.New()
	srv.Channel(wire.Channel{ID: ch, Name: "checkers"})
	return srv, ch
}

// Through the arcade's front door to the table, and the opening moves.
func TestArcadeGame(t *testing.T) {
	srv, ch := startServer(t)
	red, white := srv.Enter(ch, "alice", 80, 24), srv.Enter(ch, "bob", 80, 24)
	for _, v := range []*plugintest.Viewer{red, white} {
		srv.FrameContaining(v, "PRESS ENTER")
		srv.Key(v, "enter")
		srv.FrameContaining(v, "PRIZE COUNTER")
		srv.Key(v, "down") // TAKE A SEAT
		srv.Key(v, "enter")
	}
	frame := srv.FrameContaining(red, "YOUR MOVE")
	if !strings.ContainsAny(frame, "▀▄█") {
		t.Fatalf("no pixel pieces:\n%s", frame)
	}
	srv.Key(red, "enter") // pick up the piece on 11
	srv.FrameContaining(red, "ENTER WHERE IT LANDS")
	srv.Key(red, "esc") // the board has it, so Esc puts the piece down
	srv.FrameContaining(red, "YOUR MOVE")
	srv.Key(red, "enter")
	srv.Key(red, "up")
	srv.Key(red, "right")
	srv.Key(red, "enter") // 11-15
	srv.FrameContaining(white, "YOUR MOVE")
}

// Crowning is noticed, and moves earn the right callouts.
func TestCallouts(t *testing.T) {
	for _, c := range []struct {
		caps    int
		crowned bool
		want    string
	}{{0, false, ""}, {1, false, ""}, {2, false, "DOUBLE!"}, {3, true, "TRIPLE!"}, {4, false, "QUAD!"}, {1, true, "KING ME!"}} {
		if got := calloutFor(c.caps, c.crowned); got != c.want {
			t.Errorf("%d captures, crowned %v: %q, want %q", c.caps, c.crowned, got, c.want)
		}
	}
	// Red man on 26 steps to 30 (White's back row) and is crowned.
	g := &Game{Board: &engine.Board{ToMove: engine.Red}}
	g.Board.Sq[25] = engine.RedMan
	g.Board.Sq[20] = engine.WhiteMan
	if err := g.Play("26-30"); err != nil {
		t.Fatal(err)
	}
	if !g.Crowned || sound(g, "26-30") != "sounds/king.wav" {
		t.Fatalf("crowned %v, sound %s", g.Crowned, sound(g, "26-30"))
	}
	if err := g.Play("21-17"); err != nil {
		t.Fatal(err)
	}
	if g.Crowned {
		t.Fatal("an ordinary move counted as a crowning")
	}
}

// Every set, board and size draws inside its space, and attract mode is
// the same picture for the same frame.
func TestDrawing(t *testing.T) {
	sgr := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	for _, s := range pieceSets {
		for side := 0; side < 2; side++ {
			rows, _ := s.sprite(side, true, false)
			if len(rows) != 4 || arcade.SpriteWidth(rows, 1) != 5 {
				t.Errorf("%s: a piece isn't 5 x 4", s.ID)
			}
			if _, pal := s.puck(side, false, false); pal['m'] == "fg" && s.ID != "classic" && s.ID != "coins" {
				t.Errorf("%s: its puck fell back to plain fg", s.ID)
			}
		}
	}
	const ax, ay = 10, 3
	for i := range boardStyles {
		for _, size := range [][2]int{{42, 18}, {40, 17}, {40, 12}, {42, 10}, {21, 6}, {11, 6}} {
			for _, s := range []*pieceSet{&pieceSets[0], &pieceSets[7]} {
				c := arcade.New(80, 24, arcade.NewPalette(nil))
				sel := map[string]string{kindPieces: s.ID, kindBoard: boardPrefix + boardStyles[i].ID}
				preview(c, sel[kindBoard], sel, ax, ay, size[0], size[1], false)
				preview(c, s.ID, sel, ax, ay, size[0], size[1], true)
				for y, raw := range strings.Split(c.String(), "\n") {
					inside := y >= ay && y < ay+size[1]
					for x, r := range []rune(sgr.ReplaceAllString(raw, "")) {
						if r != ' ' && (!inside || x < ax || x >= ax+size[0]) {
							t.Fatalf("%s at %v: %q at %d,%d, outside its space", boardStyles[i].ID, size, r, x, y)
						}
					}
					if !inside && strings.Contains(raw, "\x1b[") {
						t.Fatalf("%s at %v: colour on row %d, outside its space", boardStyles[i].ID, size, y)
					}
				}
			}
		}
	}
	draw := func(frame int) string {
		c := arcade.New(80, 24, arcade.NewPalette(nil))
		attract(c, 2, 6, 76, 14, frame)
		return c.String()
	}
	for f := 0; f < 400; f += 9 {
		if draw(f) != draw(f) {
			t.Fatalf("attract mode changes within frame %d", f)
		}
	}
	if draw(0) == draw(60) {
		t.Fatal("attract mode doesn't move")
	}
}
