//go:build ignore

// gen writes Checkers' sounds into client/sounds: its own (a move, a hop,
// a combo, crowning, a draw) and the Concord Arcade kit. Run from the repo
// root: go run tools/gen.go
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/JMThomas00/Concord/sdk/arcade"
)

const rate = 22050

func main() {
	check(arcade.WriteSoundKit("client"))
	sound("click.wav", .08, func(t float64) float64 { // a piece set down on wood
		return .6*math.Sin(2*math.Pi*220*t)*math.Exp(-t*60) + .2*math.Sin(2*math.Pi*880*t)*math.Exp(-t*140)
	})
	sound("hop.wav", .16, func(t float64) float64 { // a jump: up a fifth
		f := 523.25
		if t > .07 {
			f = 783.99
		}
		return .35 * sq(f, t) * math.Exp(-math.Mod(t, .07)*18)
	})
	sound("combo.wav", .6, func(t float64) float64 { // a rising run, one note a hop
		notes := []float64{523.25, 659.25, 783.99, 1046.5, 1318.5}
		i := min(int(t/.09), len(notes)-1)
		return .3 * sq(notes[i], t) * math.Exp(-(t-float64(i)*.09)*10)
	})
	sound("king.wav", .7, func(t float64) float64 { // a little fanfare
		notes := []float64{783.99, 987.77, 1174.66, 1567.98}
		i := min(int(t/.08), len(notes)-1)
		decay := 14.0
		if i == len(notes)-1 {
			decay = 5
		}
		return .28 * sq(notes[i], t) * math.Exp(-(t-float64(i)*.08)*decay)
	})
	sound("draw.wav", .6, func(t float64) float64 { // two gentle notes, the second lower
		v := 0.0
		for i, f := range []float64{587.33, 493.88} {
			if s := t - float64(i)*.16; s >= 0 {
				v += .25 * math.Sin(2*math.Pi*f*s) * math.Exp(-s*6)
			}
		}
		return v
	})
	fmt.Println("checkers sounds and the arcade kit written to client/sounds")
}

// sq is a square wave, softened.
func sq(f, t float64) float64 { return math.Tanh(4 * math.Sin(2*math.Pi*f*t)) }

func sound(name string, dur float64, f func(t float64) float64) {
	n := int(dur * rate)
	data := make([]byte, 2*n)
	for i := 0; i < n; i++ {
		t := float64(i) / rate
		fade := math.Min(1, (dur-t)/.01)
		v := max(-1, min(1, f(t)*fade))
		binary.LittleEndian.PutUint16(data[2*i:], uint16(int16(v*32000)))
	}
	var h bytes.Buffer
	h.WriteString("RIFF")
	binary.Write(&h, binary.LittleEndian, uint32(36+len(data)))
	h.WriteString("WAVEfmt ")
	for _, v := range []any{uint32(16), uint16(1), uint16(1), uint32(rate), uint32(rate * 2), uint16(2), uint16(16)} {
		binary.Write(&h, binary.LittleEndian, v)
	}
	h.WriteString("data")
	binary.Write(&h, binary.LittleEndian, uint32(len(data)))
	h.Write(data)
	path := filepath.Join("client", "sounds", name)
	check(os.MkdirAll(filepath.Dir(path), 0o755))
	check(os.WriteFile(path, h.Bytes(), 0o644))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
