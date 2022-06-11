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

		// For single(original resolution)
		// cmd := exec.Command("ffmpeg", "-i", paths[4], "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", "index.m3u8")

		// For Multiple(360p,480p & 720p resolutions)
		cmd := exec.Command("ffmpeg", "-i", paths[4], "-map", "0:v:0", "-map", "0:a:0", "-map",
			"0:v:0", "-map", "0:a:0", "-map", "0:v:0", "-map", "0:a:0", "-c:v", "libx264", "-crf",
			"22", "-c:a", "aac", "-ar", "48000", "-filter:v:0", "scale=w=480:h=360", "-maxrate:v:0",
			"600k", "-b:a:0", "64k", "-filter:v:1", "scale=w=640:h=480", "-maxrate:v:1", "900k",
			"-b:a:1", "128k", "-filter:v:2", "scale=w=1280:h=720", "-maxrate:v:2", "1500k", "-b:a:2",
			"128k", "-var_stream_map", "v:0,a:0,name:360p v:1,a:1,name:480p v:2,a:2,name:720p",
			"-preset", "slow", "-hls_list_size", "0", "-threads", "0", "-f", "hls", "-hls_playlist_type",
			"event", "-hls_time", "10", "-hls_flags", "independent_segments", "-master_pl_name",
			"index.m3u8", "-y", "%v/index.m3u8")

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
