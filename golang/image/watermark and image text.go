package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image/color"
	"path/filepath"
)

func main() {
	m, err := imaging.Open("bg.png")
	if err != nil {
		fmt.Println("open file failed", err)
		return
	}

	//dc := gg.NewContext(m.Bounds().Dx(), m.Bounds().Dy())
	backgroundImage, err := gg.LoadImage("bg.png")
	if err != nil {
		fmt.Println("load background image error: ", err)
		return
	}
	//dc.DrawImage(backgroundImage, 0, 0)

	logoImage, err := gg.LoadImage("logo.png")
	if err != nil {
		fmt.Println("load background image error: ", err)
		return
	}

	dc := gg.NewContext(m.Bounds().Dx(), m.Bounds().Dy())
	dc.DrawImage(backgroundImage, 0, 0)
	logoX := m.Bounds().Dx() - logoImage.Bounds().Dx() - 20
	logoY := m.Bounds().Dy() - logoImage.Bounds().Dy() - 20
	dc.DrawImage(logoImage, logoX, logoY)

	//margin := 20.0
	//x := margin
	//y := margin
	//w := float64(dc.Width()) - (2.0 * margin)
	//h := float64(dc.Height()) - (2.0 * margin)
	//dc.SetColor(color.RGBA{0, 0, 0, 204})
	//dc.DrawRectangle(x, y, w, h)
	//dc.Fill()

	fontPath := filepath.Join("", "Arial-Unicode-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 80); err != nil {
		fmt.Println("load font face error: ", err)
		return
	}
	dc.SetColor(color.White)
	s := "PACE."
	marginX := 50.0
	marginY := -10.0
	textWidth, textHeight := dc.MeasureString(s)
	x := float64(dc.Width()) - textWidth - marginX
	y := float64(dc.Height()) - textHeight - marginY
	dc.DrawString(s, x, y)

	textColor := color.White
	fontPath = filepath.Join("Arial-Unicode-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 60); err != nil {
		fmt.Println("load font face error: ", err)
		return
	}
	r, g, b, _ := textColor.RGBA()
	mutedColor := color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(100),
	}
	dc.SetColor(mutedColor)
	marginY = 30
	s = "https://pace.dev/"
	_, textHeight = dc.MeasureString(s)
	x = 70
	y = float64(dc.Height()) - textHeight - marginY
	dc.DrawString(s, x, y)

	title := "Programatically generate these gorgeous social media images in Go. Programatically generate these gorgeous social media images in Go. Programatically generate these gorgeous social media images in Go."
	//title := "Programatically generate these gorgeous social media images in Go, 卧槽这个教程有点牛逼啊，我要先看看怎么使用的，牛逼啊，我的天"
	textShadowColor := color.Black

	// 文本右边距
	textRightMargin := 60.0
	// 文本上边距
	textTopMargin := 90.0
	x = textRightMargin
	y = textTopMargin

	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin
	dc.SetColor(textShadowColor)
	dc.DrawStringWrapped(title, x+1.5, y+1.5, 0, 0, maxWidth, 1.5, gg.AlignLeft)
	dc.SetColor(textColor)
	dc.DrawStringWrapped(title, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	if err := dc.SavePNG("bgtest.png"); err != nil {
		fmt.Println("save png error: ", err)
		return
	}
}
