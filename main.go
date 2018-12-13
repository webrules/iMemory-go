package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func main() {
	traverse("./samples")
	// photoExif("./samples/2018/4/18/222.jpg")
}

func traverse(root string) int {
	var photos int

	err := filepath.Walk(root,
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				if isAcceptableFormat(filepath.Ext(file)) {
					err := genThumbnail(file)
					if err != nil {
						return err
					}
					photos++
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return photos
}
func genThumbnail(originalFileName string) error {
	file, err := os.Open(originalFileName)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// // resize to width 320*240 using Lanczos resampling
	// // and preserve aspect ratio
	m := resize.Resize(320, 0, img, resize.Lanczos3)

	out, err := os.Create(targetFileName(originalFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	return nil
}

func targetFileName(originalFileName string) string {
	creationDate, lat, long := photoExif(originalFileName)
	year, month, day, hour, minute, second, ns := getTimeDetails(creationDate)

	tp := "./thumbnails/" + timeToPath([]int{year, month, day}, "/")
	err := os.MkdirAll(tp, 0700)
	if err != nil {
		log.Fatal(err)
	}

	tf := timeToPath([]int{hour, minute, second}, "")
	return tp + "/" + tf + "_" + strconv.Itoa(ns) + strconv.FormatFloat(lat, 'f', 6, 64) + "_" + strconv.FormatFloat(long, 'f', 6, 64) + ".jpg"
}

func getTimeDetails(modTime time.Time) (int, int, int, int, int, int, int) {
	year, month, day := modTime.Date()
	hour, minute, second := modTime.Clock()
	ns := modTime.Nanosecond()
	return year, int(month), day, hour, minute, second, ns

}

func timeToPath(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func isAcceptableFormat(suffix string) bool {
	suffix = strings.ToLower(suffix)
	var formats = map[string]bool{
		".jpeg": true,
		".jpg":  true,
	}

	_, ok := formats[suffix]
	return ok
}

func photoExif(file string) (time.Time, float64, float64) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	tm, _ := x.DateTime()
	// if (tm == time.Time{}) {
	// log.Fatal("No timestamp available")
	// }
	fmt.Println("Taken: ", tm)

	lat, long, _ := x.LatLong()
	// fmt.Println("lat, long: ", lat, ", ", long)
	return tm, lat, long
}
