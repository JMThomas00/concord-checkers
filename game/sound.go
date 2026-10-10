package game

import "github.com/JMThomas00/Concord/sdk/table"

// sound is what everyone watching hears after a move.
func sound(tg table.Game, move string) string {
	g := tg.(*Game)
	switch o := g.Outcome(); {
	case o.Over && o.Winner >= 0:
		return "sounds/finish.wav"
	case o.Over:
		return "sounds/draw.wav"
	case len(g.Last.Captured) >= 2:
		return "sounds/combo.wav"
	case g.Crowned:
		return "sounds/king.wav"
	case len(g.Last.Captured) == 1:
		return "sounds/hop.wav"
	}
	return "sounds/click.wav"
}
