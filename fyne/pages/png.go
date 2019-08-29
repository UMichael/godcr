package pages

//go:generate rice embed-go
import (
	"errors"

	"github.com/raedahgroup/godcr/fyne/log"

	rice "github.com/GeertJohan/go.rice"
)

var errorNoPNG = errors.New("No PNG files loaded")

var boxPng *rice.Box

// create
func init() {
	boxPng = rice.MustFindBox("Png")
}

func Png(name string) []byte {
	if boxPng == nil {
		log.Error("PNG files not loaded", name)
		return nil

	}
	b, err := boxPng.Bytes(name)
	if err != nil {
		log.Error("PNG:", name, err.Error())
	}
	return b
}
