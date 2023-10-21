package main //schlobbin, shlobbers

import (
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

// curl "https://releases.ubuntu.com/22.04.3/ubuntu-22.04.3-desktop-amd64.iso" | ./bytear 8

const sampleRate = 44100 / 2

func main() {
	chunklenS := os.Args[1]

	chunklen, err := strconv.Atoi(chunklenS)
	if err != nil {
		panic(err)
	}
	//reader, err := aces.NewBitReader(uint8(chunklen), os.Stdin)
	//if err != nil {
	//	panic(err)
	//}
	sr := beep.SampleRate(sampleRate) // 44100
	speaker.Init(sr, sr.N(time.Second/10))
	println("now playing")
	doneChan := make(chan bool, 1)
	speaker.Play(audio(os.Stdin, int64(chunklen), doneChan))
	<-doneChan
}

func audio(reader io.Reader, chunkLen int64, done chan bool) beep.Streamer {
	buf := make([]byte, chunkLen)
	num := big.NewInt(0)
	r := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(chunkLen*8), nil)
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			_, err := io.ReadFull(reader, buf)
			if err != nil {
				if err.Error() == "EOF" {
					done <- true
					return len(samples), false
				}
				panic(err)
			}

			num.SetBytes(buf)
			samples[i][0], _ = new(big.Rat).SetFrac(num, r).Float64()
			samples[i][0] = samples[i][0]*2 - 1
			//fmt.Println(samples[i][0])
			samples[i][1] = samples[i][0]
		}
		return len(samples), true
	})
}
