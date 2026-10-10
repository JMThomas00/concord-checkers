package game

import (
	"github.com/JMThomas00/Concord/sdk/arcade"
	"github.com/JMThomas00/Concord/sdk/table"
)

// The arcade front door's personality: the corner diner and game parlour.
// Tickets buy piece sets and boards at the PRIZE COUNTER; multi-jumps get
// combo calls; a draw means you split the check.
var arcadeLook = &table.Arcade{
	Title:   "CHECKERS",
	Tagline: "KING ME!",
	HowTo: []string{
		"Red moves first. Men move diagonally forward onto dark squares.",
		"Jump an opposing piece to capture it. Captures are compulsory,",
		"and a piece that can keep jumping must: a combo.",
		"A man reaching the far row is crowned king and moves both ways.",
		"No pieces or no legal move: you lose. 40 moves each with no",
		"capture and no man moving: a draw.",
		"",
		"Win games to earn Tickets (one for each new achievement, every",
		"three wins in a row and every ten games); trade them at the",
		"PRIZE COUNTER for new pieces and boards.",
	},
	Keys: []arcade.Key{
		{Key: "←↑↓→", Does: "aim"},
		{Key: "Enter", Does: "pick up, then land"},
		{Key: "Esc", Does: "put it down"},
	},
	Sounds:      true,
	Reward:      "TICKET",
	Collection:  "PRIZE COUNTER",
	Unlockables: unlockables(),
	Kinds:       []table.Kind{{ID: kindPieces, Label: "PIECES"}, {ID: kindBoard, Label: "BOARD"}},
	Preview:     preview,
	Attract:     attract,
	Result:      result,
	ResultArt:   resultArt,
}

func result(g table.Game, o table.Outcome, seats []string) string {
	if o.Winner < 0 {
		return "DRAW!"
	}
	return seats[o.Winner] + " WINS!"
}

// theCheck is a diner check, torn down the middle (the gap column).
var theCheck = []string{
	"WWWWWWWWW.WWWWWWW",
	"WkkkWkkWW.WWkkWWW",
	"WWWWWWWWW.WWWWWWW",
	"WkkWWWkkkW.WkkkWW",
	"WWWWWWWWW.WWWWWWW",
	"WkkkkkWWW.WWWRRRW",
	"WWWWWWWWW.WWWWWWW",
	"W.W.W.W.W.W.W.W.W",
}

// resultArt splits the check on a draw; wins show the final board.
func resultArt(c *arcade.Canvas, g table.Game, o table.Outcome, x, y, w, h, frame int) bool {
	if o.Winner >= 0 {
		return false
	}
	pal := map[rune]string{'W': "fg", 'k': "dim", 'R': "red"}
	var left, right []string
	for _, r := range theCheck {
		left, right = append(left, r[:9]), append(right, r[10:])
	}
	cx := x + w/2
	bob := frame / 3 % 2
	c.Sprite(cx-11, 2*(y+1)+bob, left, pal, 1)
	c.Sprite(cx+2, 2*(y+1)+1-bob, right, pal, 1)
	c.CenterIn(x, w, y+7, "SPLIT THE CHECK.", "yellow", "", true)
	if o.Reason != "" {
		c.CenterIn(x, w, y+8, o.Reason+".", "dim", "", false)
	}
	return true
}
