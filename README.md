# Checkers

American checkers (English draughts) for the terminal and for
[Concord](https://github.com/JMThomas00/Concord) channels, from the same
program.

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
   - **Seating**: *seats* (one board; sit down with Tab, everyone else
     watches), *challenge* (a lobby where members challenge each other), or
     *private* (your own games with opponents you pick).
   - **Allow spectators**, **Computer opponent**, and **Computer strength**
     (easy, normal or hard).
4. Select the channel and press **Tab** (or click the board) so your keys go
   to the game. **Tab** again opens the table menu: sit down, play the
   computer, resign, rematch. **Ctrl+]** gives the keyboard back to Concord.

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
- **Esc** cancels.
- **Tab** opens the table menu: resign, rematch, and so on.

Moves are recorded in standard notation: `11-15` for a step, `15x24x31` for jumps.

## Layout

- `engine/`: the rules and the computer player. Move generation is checked
  against the published move counts (perft) for English draughts.
- `game/`: connects the engine to the Concord SDK's table kit, and the board you play on.
- `release.go`: `go run release.go` builds the release zips Concord installs.

Tag a version (`git tag v0.1.0 && git push --tags`) and the workflow publishes them.
