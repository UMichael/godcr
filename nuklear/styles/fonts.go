package styles

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/raedahgroup/godcr/nuklear/assets"
	"golang.org/x/image/font"
)

var (
	NavFont             font.Face
	PageHeaderFont      font.Face
	PageContentFont     font.Face
	BoldPageContentFont font.Face
)

const (
	pageHeaderFontSize  = 18
	pageContentFontSize = 16
	navFontSize         = 16
)

func InitFonts() error {
	boldItalicsFontBytes, err := assets.Font("SourceSansPro-SemiboldIt.ttf")
	if err != nil {
		return err
	}

	semiBoldFontBytes, err := assets.Font("SourceSansPro-Semibold.ttf")
	if err != nil {
		return err
	}

	regularFontBytes, err := assets.Font("SourceSansPro-Regular.ttf")
	if err != nil {
		return err
	}

	NavFont, err = getFont(navFontSize, regularFontBytes)
	if err != nil {
		return err
	}

	PageHeaderFont, err = getFont(pageHeaderFontSize, boldItalicsFontBytes)
	if err != nil {
		return err
	}

	PageContentFont, err = getFont(pageContentFontSize, regularFontBytes)
	if err != nil {
		return err
	}

	BoldPageContentFont, err = getFont(pageContentFontSize, semiBoldFontBytes)
	if err != nil {
		return err
	}

	return nil
}

func getFont(fontSize int, fontBytes []byte) (font.Face, error) {
	ttfont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	options := &truetype.Options{
		Size: float64(fontSize),
	}

	return truetype.NewFace(ttfont, options), nil
}
