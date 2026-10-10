# Checkers

American checkers (English draughts) for the terminal and for
[Concord](https://github.com/JMThomas00/Concord) channels, from the same
program. In Concord it's a little arcade cabinet at the corner diner: a title
screen where the computer plays itself, pixel-art pieces in your theme's
colours, a combo call for every multi-jump (**DOUBLE!**, **TRIPLE!**), **KING
ME!**, chiptune sounds, a Hall of Fame, and **Tickets** to trade at the
**PRIZE COUNTER** for new pieces and boards. A draw means you split the check.

## Play it on your own computer

**Download:** from the [Releases](https://github.com/JMThomas00/concord-checkers/releases)
page, get the zip for your system (`concord-checkers_windows_amd64.zip`,
`concord-checkers_darwin_arm64.zip` for Apple silicon, `concord-checkers_linux_amd64.zip`, ...),
unzip it, and run the program inside from a terminal:

```sh
./concord-checkers          # Windows: .\concord-checkers.exe
```

On macOS, if it's blocked as being from an unidentified developer, run
`xattr -d com.apple.quarantine concord-checkers` once. **Or with Go installed:**
`go install github.com/JMThomas00/concord-checkers@latest`, then run `concord-checkers`.

It starts with a menu: two players on one keyboard, against the computer at
three levels, or over the network. For a network game, one player hosts and is
shown their address and a 6-character code; the other chooses join and types
both.

## Play it on a Concord server

You need to be the server owner, or have the **Manage Plugins** permission.

1. In Concord, open **Server Settings → Plugins** and press **I** (install).
2. Type `JMThomas00/concord-checkers` and press Enter. Concord downloads the latest
   release for the server's own system, verifies it, and starts it: no
   restart, no files to edit.
3. Open **Server Settings → Channels**, create a channel, and choose
   **Checkers** as its type. Its options:
   - **Seating**: *seats* (one board; sit down with M, everyone else
     watches), *challenge* (a lobby where members challenge each other), or
     *private* (your own games with opponents you pick).
   - **Allow spectators**, **Computer opponent**, and **Computer strength**
     (easy, normal or hard).
4. Select the channel and press **Tab** (or click it) so your keys go to the
   game. Press **Enter** on the title screen, then pick from the menu.

## The arcade

Everyone who opens the channel starts on the title screen. The menu:

- **1 PLAYER VS CPU** `◂ NORMAL ▸`: a game of your own against the computer
  (←/→ picks easy, normal or hard); leave and come back to carry on.
- **TAKE A SEAT** (seats channels) sits you at the channel's table; in a
  challenge channel it's **2 PLAYERS** and **WATCH**, in a private one
  **NEW GAME** and **YOUR GAMES**.
- **PRIZE COUNTER**: your pieces and board. 12 piece sets, each with a crowned
  king (classic, poker chips, bottle caps, cookies, coins, buttons, jukebox
  records, burgers and fries, grapes ... and eyeballs, planets, frogs and
  toads), and 7 boards (diner tile, gingham, walnut, mint and bubblegum, neon
  sign, vineyard). Everyone sees the game in their own.
- **HALL OF FAME**, **HOW TO PLAY** and **OPTIONS** (your sound and effects).

You earn a **Ticket** for each new achievement, every three wins in a row and
every ten games. Trade one at the PRIZE COUNTER for a locked set or board:
you're offered three and pick one.

After a game, **Enter** shows the results; Enter again asks for a rematch,
which starts once both players have. **Esc** puts down a piece you've picked
up, or goes back a screen (on the title screen it gives the keyboard back to
Concord). A pane smaller than 64 x 24 gets the plain board instead.

To update later: select it in **Server Settings → Plugins**, press **U**, then
Enter. A failed update rolls back by itself.

## Rules

8×8 board. Red moves first. Pieces move diagonally forward onto dark squares.
Captures are compulsory, and a piece that can keep jumping must. A man reaching
the far row becomes a king, which moves both ways, and crowning ends the move.
You lose when you have no pieces or no legal move. 40 moves each without a
capture or a man moving is a draw.

## Playing

- **Arrow keys** move the cursor.
- **Enter** picks up a piece, then Enter on each square it lands on. Multi-jumps are chosen one hop at a time.
- **Esc** cancels a move you've started; with nothing to cancel it hands the
  keyboard back to Concord.
- **M** opens the table menu: resign, rematch, and so on. (playing standalone, Tab does too)

Moves are recorded in standard notation: `11-15` for a step, `15x24x31` for jumps.

## Layout

- `engine/`: the rules and the computer player. Move generation is checked
  against the published move counts (perft) for English draughts.
- `game/`: connects the engine to the Concord SDK's table kit: the board you
  play on (`board.go`), the piece sets and boards (`sets.go`, `draw.go`), the
  arcade's personality (`arcade.go`: attract mode, the split check) and the
  move sounds (`sound.go`).
- `client/`: the sounds Concord sends to members' clients; `go run tools/gen.go`
  regenerates them (and the arcade sound kit).
- `release.go`: `go run release.go` builds the release zips Concord installs.

Tag a version (`git tag v0.1.0 && git push --tags`) and the workflow publishes them.

## License

MIT License — see LICENSE file for details.
