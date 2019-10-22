package widget

func StrokeLine(gtx)
square := f32.Rectangle{Max: f32.Point{X: 5, Y: float32(e.Size.Y)}}

			// Position
			op.TransformOp{}.Offset(f32.Point{ // HLdraw
				X: float32(100), // HLdraw // this should be a max width of the buttons
				Y: 0,            // HLdraw
			}).Add(gtx.Ops) // HLdraw
			// Color
			paint.ColorOp{Color: color.RGBA{0xd9, 0xd9, 0xd9, 0xff}}.Add(gtx.Ops) // HLdraw
			// Clip corners
			// Draw
			paint.PaintOp{Rect: square}.Add(gtx.Ops) // HLdraw
			// Animate
			op.InvalidateOp{}.Add(gtx.Ops) // HLdraw