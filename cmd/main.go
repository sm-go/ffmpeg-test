package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	ff "github.com/u2takey/ffmpeg-go"
)

const (
	InputFile1  = "./videos/atom.mp4"
	OutputFile1 = "./videos/out1.mp4"
	OverlayFile = "./videos/overlay.png"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("DD")

	ExampleStream(InputFile1, OutputFile1, false)
	log.Println("EDDS")

}

func ExampleStream(InputFile1, OutputFile1 string, dream bool) {
	if dream {
		panic("Use DeepDream with Tensflow haven't been implemented")
	}
	w, h := getVideoSize(InputFile1)
	log.Println(w, h)
	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()
	done1 := startFFmpegProcess1(InputFile1, pw1)
	process(pr1, pw2, w, h)
	don2 := startFFmpegProcess2(OutputFile1, pr2, w, h)
	err := <-done1
	if err != nil {
		log.Println("DD")
		panic(err)
	}
	err = <-don2
	if err != nil {
		log.Println("DED")

		panic(err)
	}
	log.Println("Done")
}

func getVideoSize(fileName string) (int, int) {
	log.Println("Getting video size for ", fileName)
	data, err := ff.Probe(fileName)
	if err != nil {
		log.Println("HAHA")
		panic(err)
	}
	log.Println("got video info", data)
	type VideoInfo struct {
		Streams []struct {
			CodeType string `json:"codec_type"`
			Width    int
			Height   int
		} `json:"streams"`
	}
	vInfo := &VideoInfo{}
	err = json.Unmarshal([]byte(data), vInfo)
	if err != nil {
		panic(err)
	}
	for _, s := range vInfo.Streams {
		if s.CodeType == "video" {
			return s.Width, s.Height
		}
	}
	return 0, 0
}

func startFFmpegProcess1(InputFile string, writer io.WriteCloser) <-chan error {
	log.Println("Starting ffmpeg process1")
	done := make(chan error)
	go func() {
		err := ff.Input(InputFile).Output("pipe:", ff.KwArgs{"format": "rawvideo", "pix_fmt": "rgb24"}).WithOutput(writer).Run()
		log.Println("ffmpeg process1 done")
		_ = writer.Close()
		done <- err
		close(done)
	}()
	return done
}

func startFFmpegProcess2(outfileName string, buf io.Reader, width, height int) <-chan error {
	log.Println("Starting ffmpeg process2")
	done := make(chan error)
	go func() {
		err := ff.Input("pipe:",
			ff.KwArgs{"format": "rawvideo",
				"pix_fmt": "rgb24", "s": fmt.Sprintf("%dx%d", width, height),
			}).
			Output(outfileName, ff.KwArgs{"pix_fmt": "yuv420p"}).
			OverWriteOutput().
			WithInput(buf).
			Run()
		log.Println("ffmpeg process2 done")
		done <- err
		close(done)
	}()
	return done
}

func process(reader io.ReadCloser, writer io.WriteCloser, w, h int) {
	go func() {
		frameSize := w * h * 3
		buf := make([]byte, frameSize, frameSize)
		for {
			n, err := io.ReadFull(reader, buf)
			if n == 0 || err == io.EOF {
				_ = writer.Close()
				return
			} else if n != frameSize || err != nil {
				panic(fmt.Sprintf("read error: %d, %s", n, err))
			}
			for i := range buf {
				buf[i] = buf[i] / 3
			}
			n, err = writer.Write(buf)
			if n != frameSize || err != nil {
				panic(fmt.Sprintf("write error: %d, %s", n, err))
			}
		}
	}()
	return
}
