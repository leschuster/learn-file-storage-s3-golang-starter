package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

const EPSILON = 1e-2

func getVideoAspectRatio(filePath string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var res ffprobeOutput
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		return "", err
	}
	if len(res.Streams) == 0 {
		return "", fmt.Errorf("not enough return values")
	}

	w, h := res.Streams[0].Width, res.Streams[0].Height
	ratio := float64(w) / float64(h)

	if math.Abs(ratio-float64(16)/float64(9)) < EPSILON {
		return "16:9", nil
	} else if math.Abs(ratio-float64(9)/float64(16)) < EPSILON {
		return "9:16", nil
	}
	return "other", nil
}
