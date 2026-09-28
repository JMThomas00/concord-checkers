# Checkers

American checkers (English draughts) for the terminal and for
[Concord](https://github.com/JMThomas00/Concord) channels, from the same
program.

- **In a terminal:** `go run .` offers two players on one keyboard, the computer at
  three levels, or a network game. For a network game, one player hosts and the other
  joins with the host's address and a 6-character code.
- **In Concord:** install it in **Server Settings → Plugins** (press **I** and
  type `JMThomas00/concord-checkers`), then create a **Checkers** channel.
  Choose how people play in the channel's settings:
  - **seats**: one board; sit down with Tab, and everyone else watches.
  - **challenge**: a lobby where members challenge each other.
  - **private**: your own games with opponents you pick.

## Rules

8×8 board. Red moves first. Pieces move diagonally forward onto dark squares.
Captures are compulsory, and a piece that can keep jumping must. A man reaching
the far row becomes a king, which moves both ways, and crowning ends the move.
You lose when you have no pieces or no legal move. 40 moves each without a
capture or a man moving is a draw.

## Playing

- **Arrow keys** move the cursor.
- **Enter** picks up a piece, then Enter on each square it lands on. Multi-jumps are chosen one hop at a time.
- **Esc** cancels.
- **Tab** opens the table menu: resign, rematch, and so on.

Moves are recorded in standard notation: `11-15` for a step, `15x24x31` for jumps.

## Layout

- `engine/`: the rules and the computer player. Move generation is checked
  against the published move counts (perft) for English draughts.
- `game/`: connects the engine to the Concord SDK's table kit, and the board you play on.
- `release.go`: `go run release.go` builds the release zips Concord installs.

Tag a version (`git tag v0.1.0 && git push --tags`) and the workflow publishes them.
