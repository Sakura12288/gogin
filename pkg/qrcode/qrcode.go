package qrcode

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"gogin/pkg/file"
	"gogin/pkg/setting"
	"gogin/pkg/util"
	"image/jpeg"
)

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

func NewQrCode(url string, width int, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Ext:    EXT_JPG,
		Level:  level,
		Mode:   mode,
	}
}

func GetQrCodeSavePath() string {
	return setting.AppSetting.QrCodeSavePath
}
func GetQrCodeFullSavePath() string {
	return setting.AppSetting.RuntimeRootPath + GetQrCodeSavePath()
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodeSavePath() + name
}
func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}
func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.Ext
	if file.CheckNotExist(src) {
		return false
	}
	return true
}

func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.Ext
	src := path + name
	if file.CheckNotExist(src) {
		qrcode, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}
		qrcode, err = barcode.Scale(qrcode, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}
		f, err := file.MustOpen(name, GetQrCodeFullSavePath())
		if err != nil {
			return "", "", err
		}
		if err := jpeg.Encode(f, qrcode, nil); err != nil {
			return "", "", err
		}
	}
	return name, path, nil
}
