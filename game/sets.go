package game

import "github.com/JMThomas00/Concord/sdk/arcade"

// The unlockables: piece sets (a man for each side, crowned for a king) and
// boards, traded for Tickets at the PRIZE COUNTER. Each player sees the
// game in their own set on their own board.

const (
	kindPieces = "pieces"
	kindBoard  = "board"
	// Board ids are prefixed in the unlockables list so they can't clash
	// with a piece set's ("classic" is both).
	boardPrefix = "board:"
)

// pieceSet is a man, 5 x 4 pixels (one rune a pixel, "." see-through):
// R is the side's colour, r its shade, h its shine; other runes are fixed
// roles. White can have a different man (burgers and fries). A king is the
// man with its top row swapped for a crown.
type pieceSet struct {
	arcade.Unlockable
	man, whiteMan []string
	red, white    map[rune]string
}

func set(id, name string, tier arcade.Tier, blurb string) arcade.Unlockable {
	return arcade.Unlockable{ID: id, Name: name, Tier: tier, Kind: kindPieces, Blurb: blurb}
}

var pieceSets = []pieceSet{
	{Unlockable: set("classic", "CLASSIC", arcade.Starter, "Red and white, since 1952."),
		man: []string{".hhh.", "RRRRR", "RRRRR", ".rrr."},
		red: map[rune]string{'R': "red", 'r': "redD", 'h': "pink"}, white: map[rune]string{'R': "fg", 'r': "dim", 'h': "fg"}},
	{Unlockable: set("chips", "POKER CHIPS", arcade.Starter, "House money."),
		man: []string{".hWh.", "RRRRR", "WRRRW", ".rWr."},
		red: map[rune]string{'R': "red", 'r': "redD", 'h': "pink", 'W': "fg"}, white: map[rune]string{'R': "cyan", 'r': "cyanD", 'h': "cyan", 'W': "fg"}},
	{Unlockable: set("caps", "BOTTLE CAPS", arcade.Common, "From the soda fountain."),
		man: []string{"R.R.R", "RhhhR", "RRRRR", ".R.R."},
		red: map[rune]string{'R': "red", 'h': "pink"}, white: map[rune]string{'R': "green", 'h': "fg"}},
	{Unlockable: set("cookies", "COOKIES", arcade.Common, "Chocolate chip against sugar."),
		man: []string{".RRR.", "RkRRk", "RRkRR", ".RRR."},
		red: map[rune]string{'R': "orange", 'k': "tire"}, white: map[rune]string{'R': "yellow", 'k': "pink"}},
	{Unlockable: set("coins", "COINS", arcade.Common, "Gold against silver."),
		man: []string{".hRR.", "hRRRR", "RRRRr", ".Rrr."},
		red: map[rune]string{'R': "yellow", 'r': "orangeD", 'h': "fg"}, white: map[rune]string{'R': "fg", 'r': "dim", 'h': "fg"}},
	{Unlockable: set("buttons", "BUTTONS", arcade.Common, "From the jar by the till."),
		man: []string{".RRR.", "RkRkR", "RRRRR", ".RRR."},
		red: map[rune]string{'R': "pink", 'k': "tire"}, white: map[rune]string{'R': "cyan", 'k': "tire"}},
	{Unlockable: set("records", "JUKEBOX RECORDS", arcade.Rare, "B-sides only."),
		man: []string{".ddd.", "dkRkd", "dkRkd", ".ddd."},
		red: map[rune]string{'k': "tire", 'd': "comment", 'R': "red"}, white: map[rune]string{'k': "tire", 'd': "comment", 'R': "cyan"}},
	{Unlockable: set("burgers", "BURGERS & FRIES", arcade.Rare, "Order up."),
		man: []string{".ooo.", "GGGGG", "rrrrr", ".ooo."}, whiteMan: []string{"Y.Y.Y", "YYYYY", "RRRRR", ".RRR."},
		red: map[rune]string{'o': "orange", 'G': "green", 'r': "redD"}, white: map[rune]string{'Y': "yellow", 'R': "red"}},
	{Unlockable: set("grapes", "GRAPES", arcade.Rare, "Purple against green."),
		man: []string{"..g..", ".RhR.", "RRRhR", ".RRR."},
		red: map[rune]string{'R': "purple", 'h': "hi", 'g': "green"}, white: map[rune]string{'R': "green", 'h': "fg", 'g': "orange"}},
	{Unlockable: set("eyes", "EYEBALLS", arcade.Legendary, "They're watching your next move."),
		man: []string{".WWW.", "WWRWW", "WWkWW", ".WWW."},
		red: map[rune]string{'W': "fg", 'R': "red", 'k': "tire"}, white: map[rune]string{'W': "fg", 'R': "cyan", 'k': "tire"}},
	{Unlockable: set("planets", "PLANETS", arcade.Legendary, "Ringed, every one."),
		man: []string{".RRR.", "yyyyy", "RRRRR", ".rrr."},
		red: map[rune]string{'R': "orange", 'r': "orangeD", 'y': "yellow"}, white: map[rune]string{'R': "cyan", 'r': "cyanD", 'y': "fg"}},
	{Unlockable: set("frogs", "FROGS & TOADS", arcade.Legendary, "Jumping is what they do."),
		man: []string{"W.W.W", "RRRRR", "RkRkR", ".R.R."},
		red: map[rune]string{'W': "fg", 'R': "green", 'k': "tire"}, white: map[rune]string{'W': "fg", 'R': "orange", 'k': "tire"}},
}

// boardStyle is a board's colours: light and dark squares, and an optional
// frame.
type boardStyle struct {
	arcade.Unlockable
	light, dark, edge string
}

func board(id, name string, tier arcade.Tier, blurb, light, dark, edge string) boardStyle {
	return boardStyle{arcade.Unlockable{ID: id, Name: name, Tier: tier, Kind: kindBoard, Blurb: blurb}, light, dark, edge}
}

var boardStyles = []boardStyle{
	board("classic", "CLASSIC", arcade.Starter, "Red and black.", "redB", "tire", ""),
	board("tile", "DINER TILE", arcade.Starter, "Mind the mop.", "fgD", "tire", ""),
	board("gingham", "GINGHAM", arcade.Common, "The tablecloth.", "redB", "fgB", ""),
	board("walnut", "WALNUT", arcade.Common, "Grandpa's board.", "orangeB", "orangeD", ""),
	board("mint", "MINT & BUBBLEGUM", arcade.Rare, "A milkshake menu.", "pinkB", "greenB", ""),
	board("neon", "NEON SIGN", arcade.Rare, "Open 24 hours.", "purpleB", "bg", "pink"),
	board("vineyard", "VINEYARD", arcade.Legendary, "Squares of vines.", "purpleB", "greenB", "green"),
}

// unlockables is every piece set and board, for the arcade.
func unlockables() []arcade.Unlockable {
	var out []arcade.Unlockable
	for _, s := range pieceSets {
		out = append(out, s.Unlockable)
	}
	for _, b := range boardStyles {
		u := b.Unlockable
		u.ID = boardPrefix + u.ID
		out = append(out, u)
	}
	return out
}

// pieces looks up a piece set (the first starter if id is unknown).
func pieces(id string) *pieceSet {
	for i := range pieceSets {
		if pieceSets[i].ID == id {
			return &pieceSets[i]
		}
	}
	return &pieceSets[0]
}

// style looks up a board (with or without its prefix).
func style(id string) *boardStyle {
	if len(id) > len(boardPrefix) && id[:len(boardPrefix)] == boardPrefix {
		id = id[len(boardPrefix):]
	}
	for i := range boardStyles {
		if boardStyles[i].ID == id {
			return &boardStyles[i]
		}
	}
	return &boardStyles[0]
}

var crown = "Y.Y.Y"

// sprite is a piece's picture: side 0 Red, 1 White.
func (s *pieceSet) sprite(side int, king, ghost bool) ([]string, map[rune]string) {
	rows, src := s.man, s.red
	if side == 1 {
		src = s.white
		if s.whiteMan != nil {
			rows = s.whiteMan
		}
	}
	pal := map[rune]string{'Y': "yellow"}
	for k, v := range src {
		pal[k] = v
	}
	if king {
		rows = append([]string{crown}, rows[1:]...)
	}
	if ghost {
		for k := range pal { // "comment", not "ghost": it has to read on dark squares
			pal[k] = "comment"
		}
	}
	return rows, pal
}

// puck is a piece seen side on, 5 x 2 pixels, for one-row squares: its
// main colour, and a gold top for a king.
func (s *pieceSet) puck(side int, king, ghost bool) ([]string, map[rune]string) {
	src := s.red
	if side == 1 {
		src = s.white
	}
	main, top := src['R'], src['h']
	for _, k := range []rune{'W', 'o', 'Y'} {
		if main == "" {
			main = src[k]
		}
	}
	if main == "" {
		main = "fg"
	}
	if top == "" {
		top = main
	}
	if king {
		top = "yellow"
	}
	if ghost {
		main, top = "comment", "comment"
	}
	return []string{".ttt.", "mmmmm"}, map[rune]string{'t': top, 'm': main}
}
