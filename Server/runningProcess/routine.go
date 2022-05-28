package runningProcess

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func HlsVideoConversion(resultChan <-chan string) {

	for result := range resultChan {
		paths := strings.Split(result, "/")

		path := paths[0] + "/" + paths[1] + "/" + paths[2] + "/" + paths[3]

		err := os.Chdir(path)

		if err != nil {
			fmt.Println(err.Error())
		}

		cmd := exec.Command("ffmpeg", "-i", paths[4], "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", "index.m3u8")

		err = cmd.Run()

		if err != nil {
			fmt.Println(err.Error())
		}

		imageFileName := strings.Split(paths[4], ".")[0] + ".png"

		cmd = exec.Command("ffmpeg", "-i", paths[4], "-ss", "00:00:01.000", "-vframes", "1", imageFileName)

		err = cmd.Run()

		if err != nil {
			fmt.Println(err.Error())
		}

		err = os.Remove(paths[4])

		if err != nil {
			fmt.Println(err.Error())
		}

		err = os.Chdir("../../../Server")
		if err != nil {
			fmt.Println(err.Error())
		}

	}

}
