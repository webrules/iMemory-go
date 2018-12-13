package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var of1 = "./samples/2018/3/12/111.jpg"
var of2 = "./samples/2018/4/5/222.jpg"
var tf1 = "./thumbnails/2018/3/12/17623_031.217669_121.469039.jpg"
var tf2 = "./thumbnails/2018/4/5/10857_031.176947_121.445183.jpg"

func TestIsAcceptableFormat(t *testing.T) {
	assert.True(t, isAcceptableFormat(".jpg"))
	assert.False(t, isAcceptableFormat(".gif"))
	// assert.True(t, isAcceptableFormat(".png"))
}

func TestGetTimeDetails(t *testing.T) {
	date := time.Date(2018, 12, 9, 1, 2, 3, 4, time.UTC)
	year, month, day, hour, minute, second, ns := getTimeDetails(date)
	assert.Equal(t, 2018, year)
	assert.Equal(t, 12, month)
	assert.Equal(t, 9, day)
	assert.Equal(t, 1, hour)
	assert.Equal(t, 2, minute)
	assert.Equal(t, 3, second)
	assert.Equal(t, 4, ns)
}
func TestTimeToPath(t *testing.T) {
	tp := timeToPath([]int{2006, 1, 26}, "/")
	assert.Equal(t, "2006/1/26", tp)

	tf := timeToPath([]int{11, 36, 28, 1234}, "")
	assert.Equal(t, "1136281234", tf)
}

func TestTargetFileName(t *testing.T) {
	assert.Equal(t, tf1, targetFileName(of1))
}

func TestGenThumbnail(t *testing.T) {
	path := of1
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	genThumbnail(path)

	fileInfo, err = os.Stat(tf1)
	assert.True(t, fileInfo.Size() > 0)
}

func TestTraverse(t *testing.T) {
	assert.Equal(t, 6, traverse("./samples"))
}
func TestPhotoExif(t *testing.T) {
	creationDate, lat, long := photoExif(of2)
	assert.Equal(t, 31.176947222222225, lat)
	assert.Equal(t, 121.44518333333333, long)

	year, month, day, hour, minute, second, ns := getTimeDetails(creationDate)
	assert.Equal(t, 2018, year)
	assert.Equal(t, 4, month)
	assert.Equal(t, 5, day)
	assert.Equal(t, 10, hour)
	assert.Equal(t, 8, minute)
	assert.Equal(t, 57, second)
	assert.Equal(t, 0, ns)

}
