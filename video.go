package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
)

const EPSILON = 1e-2

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %v", err)
	}

	var output struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		return "", err
	}
	if len(output.Streams) == 0 {
		return "", errors.New("no video stream found")
	}

	w, h := output.Streams[0].Width, output.Streams[0].Height
	ratio := float64(w) / float64(h)

	if math.Abs(ratio-float64(16)/float64(9)) < EPSILON {
		return "16:9", nil
	} else if math.Abs(ratio-float64(9)/float64(16)) < EPSILON {
		return "9:16", nil
	}
	return "other", nil
}

func processVideoForFastStart(filePath string) (string, error) {
	processedFilePath := filePath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", processedFilePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error processing video: %s, %v", stderr.String(), err)
	}

	fileInfo, err := os.Stat(processedFilePath)
	if err != nil {
		return "", fmt.Errorf("could not stat processed file: %v", err)
	}
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("processed file is empty")
	}

	return processedFilePath, nil
}
