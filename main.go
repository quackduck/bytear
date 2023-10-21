package main //schlobbin, shlobbers

import (
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/quackduck/aces"
)

// curl "https://releases.ubuntu.com/22.04.3/ubuntu-22.04.3-desktop-amd64.iso" | ./bytear 8

const sampleRate = 44100 / 2

func main() {
	chunklenS := os.Args[1]

	chunklen, err := strconv.Atoi(chunklenS)
	if err != nil {
		panic(err)
	}
	reader, err := aces.NewBitReader(uint8(chunklen), os.Stdin)
	if err != nil {
		panic(err)
	}
	sr := beep.SampleRate(sampleRate) // 44100
	speaker.Init(sr, sr.N(time.Second/10))
	println("now playing")
	doneChan := make(chan bool, 1)
	speaker.Play(audio(reader, uint8(chunklen), doneChan))
	<-doneChan
}

func audio(reader *aces.BitReader, chunkLen uint8, done chan bool) beep.Streamer {
	r := 1 << chunkLen
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			b, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					done <- true
					return len(samples), false
				}
				panic(err)
			}
			samples[i][0] = float64(b) / float64(r)
			samples[i][1] = samples[i][0]
		}
		return len(samples), true
	})
}
