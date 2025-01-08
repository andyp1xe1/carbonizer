package utils

import (
	"image"
	"image/png"
	"log"
	"os"
	"strings"
	"time"
)

func SaveImage(startName string, img image.Image) (string, error) {

	resName := "./results/" + startName + "_" + time.Now().Format("060102_150405") + ".png"
	resName = strings.ReplaceAll(resName, " ", "_")
	rf, err := os.Create(resName)
	if err != nil {
		log.Fatal(err)
	}
	defer rf.Close()

	png.Encode(rf, img)

	return resName, nil
}
