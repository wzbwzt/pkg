package image

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"

	"github.com/chai2010/webp"
)

type WatermarkPosition int

const (
	UpperLeft    WatermarkPosition = 0
	UpperCenter  WatermarkPosition = 1
	UpperRight   WatermarkPosition = 2
	CenterLeft   WatermarkPosition = 3
	Center       WatermarkPosition = 4
	CenterRight  WatermarkPosition = 5
	BottomLeft   WatermarkPosition = 6
	BottomCenter WatermarkPosition = 7
	BottomRight  WatermarkPosition = 8
)

var waterImage *image.Image

func InitWaterMarker(base64Image string) {
	if base64Image != "" {
		data, err := base64.StdEncoding.DecodeString(base64Image)
		if err == nil {
			waterImage0, err := png.Decode(bytes.NewReader(data))
			if err == nil {
				waterImage = &waterImage0
			}
		}
	}
}

func loadFont(path string) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	return f, err
}

func setFontFace(ctx *gg.Context, f *truetype.Font, points float64) {
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
	})
	ctx.SetFontFace(face)
}

var font *truetype.Font

func MarkerText(bg image.Image, text string, fontsize int, fontColor color.Color,
	position WatermarkPosition, hpadding int, vpadding int) image.Image {

	width := bg.Bounds().Size().X
	height := bg.Bounds().Size().Y

	fwidth := float64(width)
	fheight := float64(height)
	fhpadding := float64(hpadding)
	fvpadding := float64(vpadding)

	dc := gg.NewContext(width, height)

	if font == nil {
		font, _ = loadFont("xxxx")
	}

	setFontFace(dc, font, float64(fontsize))

	dc.DrawImage(bg, 0, 0)
	dc.SetColor(fontColor)

	tw, th := dc.MeasureString(text)

	// 计算位置
	var px float64 = 0
	var py float64 = 0

	switch position {
	case UpperLeft:
		px = fhpadding
		py = fvpadding + th
	case UpperCenter:
		py = fvpadding + th
		px = (fwidth - tw) / 2
	case UpperRight:
		py = fvpadding + th
		px = (fwidth - tw - fhpadding)
	case CenterLeft:
		px = fhpadding
		py = (fheight + th) / 2
	case Center:
		px = (fwidth - tw) / 2
		py = (fheight + th) / 2
	case CenterRight:
		py = (fheight + th) / 2
		px = (fwidth - tw - fhpadding)
	case BottomLeft:
		px = fhpadding
		py = fheight - fhpadding
	case BottomCenter:
		px = (fwidth - tw) / 2
		py = fheight - fhpadding
	case BottomRight:
		px = (fwidth - tw - fhpadding)
		py = fheight - fhpadding
	}

	dc.DrawString(text, px, py)
	dc.Clip()

	return dc.Image()
}

func ThumbImage(img image.Image, destWidth int) (image.Image, error) {
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y
	newWidth := destWidth
	newHeight := int(float64(height) * (float64(destWidth) / float64(width)))

	thumb := imaging.Fill(img, newWidth, newHeight, imaging.Center, imaging.Lanczos)
	return thumb, nil
}

func ThumbCropImage(img image.Image, destWidth, newHeight int) (image.Image, error) {
	thumb := imaging.Fill(img, destWidth, newHeight, imaging.Center, imaging.Lanczos)
	return thumb, nil
}

func MarkerImageWarp(bg *[]byte, result *[]byte, format string) error {

	if format != "webp" && format != "png" && format != "jpg" && format != "jpeg" {
		return errors.New("unsport img format")
	}

	if waterImage == nil {
		return nil
	}

	bgimg, _, err := image.Decode(bytes.NewReader(*bg))
	if err != nil {
		return err
	}

	resultImg := MarkerImage(bgimg, *waterImage, 1, UpperLeft, 0, 0)

	buf := new(bytes.Buffer)
	if format == "webp" {
		err = webp.Encode(buf, resultImg, &webp.Options{
			Quality: 70,
		})
		if err != nil {
			return err
		}
	}

	if format == "png" {
		err = png.Encode(buf, resultImg)
		if err != nil {
			return err
		}
	}

	if format == "jpg" || format == "jpeg" {
		jpeg.Encode(buf, resultImg, &jpeg.Options{
			Quality: 70,
		})
	}

	*result = buf.Bytes()
	return nil
}

func MarkerImage(bg image.Image, marker image.Image, scale float64, position WatermarkPosition, hpadding int, vpadding int) image.Image {

	width := marker.Bounds().Size().X
	height := marker.Bounds().Size().Y
	markFit := imaging.Fill(marker, int(float64(width)*scale), int(float64(height)*scale), imaging.Center, imaging.Lanczos)
	markerWidth := markFit.Rect.Size().X
	markerHeight := markFit.Rect.Size().Y

	bgWidth := bg.Bounds().Size().X
	bgHeight := bg.Bounds().Size().Y

	bgFit := imaging.Fit(bg, bgWidth, bgHeight, imaging.Lanczos)

	px := 0
	py := 0
	switch position {
	case UpperLeft:
		px = hpadding
		py = vpadding
	case UpperCenter:
		py = vpadding
		px = (bgWidth - markerWidth) / 2
	case UpperRight:
		py = vpadding
		px = (bgWidth - markerWidth - hpadding)
	case CenterLeft:
		px = hpadding
		py = (bgHeight - markerHeight) / 2
	case Center:
		px = (bgWidth - markerWidth) / 2
		py = (bgHeight - markerHeight) / 2
	case CenterRight:
		py = (bgHeight - markerHeight) / 2
		px = (bgWidth - markerWidth - hpadding)
	case BottomLeft:
		px = hpadding
		py = (bgHeight - markerHeight - vpadding)
	case BottomCenter:
		px = (bgWidth - markerWidth) / 2
		py = (bgHeight - markerHeight - vpadding)
	case BottomRight:
		px = (bgWidth - markerWidth - hpadding)
		py = (bgHeight - markerHeight - vpadding)
	}

	dst := imaging.Overlay(bgFit, markFit, image.Pt(px, py), 1.0)
	return dst
}
