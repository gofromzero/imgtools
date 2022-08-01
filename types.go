package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"github.com/skip2/go-qrcode"
)

type ImgFile struct {
	Baseimg string  `json:"baseimg"`
	Items   []Items `json:"items"`
}

func (f *ImgFile) CheckParams() (err error) {
	if !IsFileExist(f.Baseimg) {
		return errors.New("baseimg path err")
	}

	for i, item := range f.Items {
		err = item.CheckParams()
		if err != nil {
			fmt.Errorf(" No. %d err: %v", i, err.Error())
			return
		}
	}
	return nil
}

func (f *ImgFile) Draw(no int, filedir string) (err error) {
	baseimage, err := gg.LoadImage(f.Baseimg)
	if err != nil {
		log.Fatalln(err)
	}

	dc := gg.NewContextForImage(baseimage)

	for _, item := range f.Items {
		err = item.Draw(dc)
		if err != nil {
			return
		}
	}

	ext := path.Ext(f.Baseimg)
	baseName := path.Base(f.Baseimg)
	name := strings.Trim(baseName, ext)

	file1name := fmt.Sprintf("%s_%d.png", name, no)
	fpath := filepath.Join(filedir, file1name)
	err = dc.SavePNG(fpath)
	if err != nil {
		return
	}
	log.Printf("%s 生成...\n", file1name)
	return
}

type Items struct {
	Type  int     `json:"type"`
	Value string  `json:"value"`
	Size  float64 `json:"size"`
	Point struct {
		X float64
		Y float64
	} `json:"point"`
	APoint struct {
		X float64
		Y float64
	} `json:"apoint"`
	RGB   []int  `json:"rgb"`
	Font  string `json:"font"`
	Scale struct {
		X float64
		Y float64
	}
}

const (
	word = iota
	itemQrcode
	img
)

func (f *Items) CheckParams() (err error) {

	if f.Value == "" {
		return errors.New("item value is ''")
	}
	switch f.Type {
	case word: // 文字
		if f.Font == "" {
			return errors.New("word type: font path is ''")
		}
		if len(f.RGB) != 3 {
			return errors.New("RGB params num not equal 3")
		}
	case itemQrcode: // 二维码
	case img: // 图片
	default:
		return errors.New("items type is unknown")
	}

	return nil
}

func (f *Items) Draw(dc *gg.Context) (err error) {
	dc.SetRGB255(f.RGB[0], f.RGB[1], f.RGB[2])

	switch f.Type {
	case word: // 文字
		if err = dc.LoadFontFace(f.Font, f.Size); err != nil {
			return err
		}

		_, h := dc.MeasureString(f.Value)

		dc.DrawStringAnchored(f.Value, f.Point.X, f.Point.Y+h, f.APoint.X, f.APoint.Y)
	case itemQrcode: // 二维码
		var qrc *qrcode.QRCode
		qrc, err = qrcode.New(f.Value, qrcode.Medium)
		if err != nil {
			log.Printf("could not generate QRCode: %v", err)
			return
		}
		qrc.DisableBorder = true
		dc.DrawImageAnchored(qrc.Image(int(f.Size)), int(f.Point.X), int(f.Point.Y), f.APoint.X, f.APoint.Y) // 二维码
	case img: // 图片
		var itemImg image.Image
		itemImg, err = gg.LoadImage(f.Value)
		if err != nil {
			return
		}
		itemdc := gg.NewContextForImage(itemImg)

		dc.ScaleAbout(f.Scale.X, f.Scale.Y, f.Point.X, f.Point.Y)
		dc.DrawImageAnchored(itemdc.Image(), int(f.Point.X), int(f.Point.Y), f.APoint.X, f.APoint.Y) // 图像

	default:
		return errors.New("items type is unknown")
	}

	dc.Fill()

	return nil

}
