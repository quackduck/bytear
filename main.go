package main //schlobbin, shlobbers

import (
	"bufio"
	"io"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/youpy/go-wav"
)

// curl "https://releases.ubuntu.com/22.04.3/ubuntu-22.04.3-desktop-amd64.iso" | ./bytear 8

// const sampleRate = 640000 //44100 / 2
//const sampleRate = 1000

const sampleRate = 44100 / 2

func main() {
	file, err := os.Create("out.wav")
	if err != nil {
		panic(err)
	}
	chunkLenS := os.Args[1]

	chunkLen, err := strconv.Atoi(chunkLenS)
	if err != nil {
		panic(err)
	}
	writeWAV(bufio.NewReaderSize(os.Stdin, 1024*1024), int64(chunkLen), file, 60*3, sampleRate)
	return

	writer := wav.NewWriter(file, sampleRate*60, 1, sampleRate, uint16(chunkLen*8))
	//reader, err := aces.NewBitReader(uint8(chunklen), os.Stdin)
	//if err != nil {
	//	panic(err)
	//}
	sr := beep.SampleRate(sampleRate) // 44100
	speaker.Init(sr, sr.N(time.Second/10))
	println("now playing")
	doneChan := make(chan bool, 1)
	speaker.Play(audio(os.Stdin, int64(chunkLen), doneChan, writer))
	<-doneChan
}

func writeWAV(reader io.Reader, chunkLen int64, writer io.Writer, seconds int, sampleRate int) {
	buf := make([]byte, chunkLen)
	num := big.NewInt(0)
	//r := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(chunkLen*8), nil)
	w := wav.NewWriter(writer, uint32(sampleRate*seconds), 1, uint32(sampleRate), uint16(chunkLen*8))
	s := make([]wav.Sample, seconds*sampleRate)
	for i := range s {
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			if err.Error() == "EOF" {
				break
				//return
			}
			panic(err)
		}
		num.SetBytes(buf)
		s[i] = wav.Sample{Values: [2]int{int(num.Int64()), int(num.Int64())}}
		//if i%(sampleRate*30) == 0 {
		//	fmt.Println(float64(i) / float64(seconds*sampleRate))
		//}
	}
	err := w.WriteSamples(s)
	if err != nil {
		panic(err)
	}
}

func audio(reader io.Reader, chunkLen int64, done chan bool, writer *wav.Writer) beep.Streamer {
	buf := make([]byte, chunkLen)
	num := big.NewInt(0)
	r := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(chunkLen*8), nil)
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		s := make([]wav.Sample, len(samples))
		for i := range samples {
			_, err := io.ReadFull(reader, buf)
			if err != nil {
				if err.Error() == "EOF" {
					err = writer.WriteSamples(s)
					if err != nil {
						panic(err)
					}
					done <- true
					return len(samples), false
				}
				panic(err)
			}

			num.SetBytes(buf)
			s[i].Values[0] = int(num.Int64())
			s[i].Values[1] = int(num.Int64())

			samples[i][0], _ = new(big.Rat).SetFrac(num, r).Float64()
			samples[i][0] = samples[i][0]*2 - 1
			//fmt.Println(samples[i][0])
			samples[i][1] = samples[i][0]
		}
		//fmt.Println(s)
		err := writer.WriteSamples(s)
		if err != nil {
			panic(err)
		}
		return len(samples), true
	})
}
