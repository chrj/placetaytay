package taytay

import (
	"encoding/json"
	"image"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"
	_ "image/png"
)

type TaylorSwift struct {
	Filename string `json:"-"`
	Bounds   image.Rectangle
}

var taytays = []*TaylorSwift{}

func (ts *TaylorSwift) Reader() io.ReadCloser {
	f, err := os.Open(ts.Filename)
	if err != nil {
		log.Printf("couldn't open taytay: %v", err)
		return nil
	}
	return f
}

func (ts *TaylorSwift) Image() image.Image {

	rd := ts.Reader()

	if rd == nil {
		return nil
	}

	defer rd.Close()

	im, _, err := image.Decode(rd)
	if err != nil {
		log.Printf("couldn't decode taytay: %v", err)
		return nil
	}

	return im

}

func (ts *TaylorSwift) AspectRatio() float64 {
	return float64(ts.Bounds.Max.X) / float64(ts.Bounds.Max.Y)
}

func (ts *TaylorSwift) saveMeta() error {

	f, err := os.OpenFile(ts.Filename+".meta", os.O_CREATE|os.O_RDWR, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	return json.NewEncoder(f).Encode(ts)

}

func load(filename string) error {

	if strings.HasSuffix(filename, ".html") {
		return nil
	}

	if strings.HasSuffix(filename, ".meta") {
		return nil
	}

	m, err := os.Open(filename + ".meta")

	if err == nil {
		defer m.Close()
		ts := &TaylorSwift{
			Filename: filename,
		}
		if err := json.NewDecoder(m).Decode(ts); err != nil {
			return err
		}
		taytays = append(taytays, ts)
		return nil
	}

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	imf, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer imf.Close()

	im, _, err := image.Decode(imf)
	if err != nil {
		return err
	}

	ts := &TaylorSwift{
		Bounds:   im.Bounds(),
		Filename: filename,
	}

	if err := ts.saveMeta(); err != nil {
		return err
	}

	taytays = append(taytays, ts)

	return nil

}

func Index(path string) error {

	d, err := os.Open(path)
	if err != nil {
		return err
	}

	defer d.Close()

	files, err := d.Readdir(-1)

	for _, fi := range files {

		if err := load(filepath.Join(path, fi.Name())); err != nil {
			log.Printf("error opening image: %v:", fi.Name(), err)
			continue
		}

		log.Printf("opened image: %v:", fi.Name())

	}

	return nil
}

func Closest(width, height int) *TaylorSwift {

	ratio := float64(width) / float64(height)

	closestDiff := float64(1000)
	closestIdx := 0

	for idx, ts := range taytays {

		if ts.Bounds.Max.X < width || ts.Bounds.Max.Y < height {
			continue
		}

		aspectDiff := math.Abs(ts.AspectRatio() - ratio)

		sizeDiff := float64(ts.Bounds.Max.X) / float64(width) *
			float64(ts.Bounds.Max.Y) / float64(height)

		diff := aspectDiff * sizeDiff

		if diff < closestDiff {
			closestIdx = idx
			closestDiff = diff
		}

	}

	return taytays[closestIdx]

}
