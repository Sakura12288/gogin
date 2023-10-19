package article_service

import (
	"github.com/golang/freetype"
	"gogin/pkg/file"
	"gogin/pkg/qrcode"
	"gogin/pkg/setting"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

type ArticlePoster struct {
	PosterName string
	*Article
	Qr *qrcode.QrCode
}

func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
		Qr:         qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}
func (a *ArticlePoster) CheckMergeImage(path string) bool {
	if file.CheckNotExist(path + a.PosterName) {
		return false
	}
	return true
}
func (a *ArticlePoster) OpenMergeImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//下面为背景图

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Point
}
type Rect struct {
	Name           string
	X0, Y0, X1, Y1 int
}
type Point struct {
	X, Y int
}

type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	Subtitle string
	X1       int
	Y1       int
	Size1    float64
}

const (
	FontName = "msyhbd.ttc"
)

func (a *ArticlePosterBg) DrawPoster(d *DrawText, fontName string) error {
	fontSource := setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath + fontName
	fontSourceBytes, err := os.ReadFile(fontSource)
	if err != nil {
		return err
	}
	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}
	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(trueTypeFont)
	fc.SetFontSize(d.Size0)
	fc.SetClip(d.JPG.Bounds())
	fc.SetDst(d.JPG)
	fc.SetSrc(image.Black)

	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}
	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.Subtitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}
	return nil
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Point) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Point:         pt,
	}
}
func (a *ArticlePosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullSavePath()
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}
	if !a.CheckMergeImage(path) {
		merdF, err := a.OpenMergeImage(path)
		if err != nil {
			return "", "", err
		}
		defer merdF.Close()

		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()

		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}
		jpg := image.NewRGBA(image.Rect(a.X0, a.Y0, a.X1, a.Y1))

		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Point{a.X, a.Y}), draw.Over)

		err = a.DrawPoster(&DrawText{
			JPG:    jpg,
			Merged: merdF,
			Title:  "私奔是傻逼 wf9e",
			X0:     80,
			Y0:     160,
			Size0:  42,

			Subtitle: "猴子也是",
			X1:       320,
			Y1:       220,
			Size1:    36,
		}, FontName)
		if err != nil {
			return "", "", err
		}
	}
	return fileName, path, nil
}
