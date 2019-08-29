package assets

//go:generate rice embed-go
import (
	"errors"

	rice "github.com/GeertJohan/go.rice"
)

var errorNoFont = errors.New("No Font files loaded")

var boxFont *rice.Box

// create
func init() {
	boxFont = rice.MustFindBox("font")
}

func Font(name string) ([]byte, error) {
	if boxFont == nil {
		return nil, errorNoFont
	}
	return boxFont.Bytes(name)
}
