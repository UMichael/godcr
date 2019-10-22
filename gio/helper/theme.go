package helper

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Theme struct {
	material.Theme
}

const (
	fontSize = 13
)

func NewTheme() *Theme {
	t := &Theme{
		Theme: material.Theme{Shaper: font.Default()},
	}
	t.Color.Primary = color.RGBA{63, 81, 181, 255}
	t.Color.Text = color.RGBA{0, 0, 0, 255}
	t.Color.Hint = color.RGBA{187, 187, 187, 255}
	t.TextSize = unit.Sp(16)
	return t
}

// func getFont() text.Font {
// 	return text.Font{
// 		Size: unit.Dp(fontSize),
// 	}
// }

// func getItalicFont() text.Font {
// 	return text.Font {
// 		Size: unit.Dp(fontSize),
// 		Style: text.Italic,
// 	}
// }
