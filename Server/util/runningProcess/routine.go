package runningProcess

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func HlsVideoConversion(resultChan <-chan string) {

	for result := range resultChan {
		dir, file := filepath.Split(result)

		// For single(original resolution)
		// cmd := exec.Command("ffmpeg", "-i", paths[4], "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", "index.m3u8")

		// For Multiple(360p,480p & 720p resolutions)
		cmd := exec.Command("ffmpeg", "-i", result, "-map", "0:v:0", "-map", "0:a:0", "-map",
			"0:v:0", "-map", "0:a:0", "-map", "0:v:0", "-map", "0:a:0", "-c:v", "libx264", "-crf",
			"22", "-c:a", "aac", "-ar", "48000", "-filter:v:0", "scale=w=480:h=360", "-maxrate:v:0",
			"600k", "-b:a:0", "64k", "-filter:v:1", "scale=w=640:h=480", "-maxrate:v:1", "900k",
			"-b:a:1", "128k", "-filter:v:2", "scale=w=1280:h=720", "-maxrate:v:2", "1500k", "-b:a:2",
			"128k", "-var_stream_map", "v:0,a:0,name:360p v:1,a:1,name:480p v:2,a:2,name:720p",
			"-preset", "slow", "-hls_list_size", "0", "-threads", "0", "-f", "hls", "-hls_playlist_type",
			"event", "-hls_time", "10", "-hls_flags", "independent_segments", "-master_pl_name",
			"index.m3u8", "-y", dir+"%v/index.m3u8")

		err := cmd.Run()

		if err != nil {
			log.Println(err)
			continue
		}

		imageFileName := strings.Split(file, ".")[0] + ".png"

		cmd = exec.Command("ffmpeg", "-i", result, "-ss", "00:00:01.000", "-vframes", "1", filepath.Join(dir, imageFileName))

		err = cmd.Run()

		if err != nil {
			log.Println(err)
			continue
		}

		err = os.Remove(result)

		if err != nil {
			log.Println(err)
		}
	}

}

func UnzipFile(resultPaths <-chan string) {

	dir := "../media/examProofs/"

	for result := range resultPaths {
		reader, err := zip.OpenReader(result)

		if err != nil {
			log.Println(err)
			continue
		}

		for _, file := range reader.File {

			fpath := dir + file.Name

			// Make Folder
			// if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			// 	log.Println(err)
			// 	continue
			// }

			// Create/Open dst File
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, file.Mode())
			if err != nil {
				log.Println(err)
				continue
			}

			// Open src File
			inFile, err := file.Open()
			if err != nil {
				log.Println(err)
				continue
			}

			// Copy src to dst
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				log.Println(err)
				continue
			}

			// Close the file without defer to close before next iteration of loop
			outFile.Close()
			inFile.Close()
		}
		reader.Close()

		err = os.Remove(result)
		if err != nil {
			log.Println(err)
		}
	}
}
