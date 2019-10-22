package gio

import (
	"image"
	"log"

	"gioui.org/io/system"

	gioapp "gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr/gio/helper"
)

type Desktop struct {
	window      *gioapp.Window
	pages       []page
	currentPage *page
	pageChanged bool

	theme *helper.Theme
}

const (
	appName      = "GoDcr"
	windowWidth  = 450
	windowHeight = 350

	navSectionWidth = 120
)

func LaunchApp() {
	desktop := &Desktop{
		theme: helper.NewTheme(),
	}
	desktop.prepareHandlers()

	go func() {
		desktop.window = gioapp.NewWindow(
			gioapp.Size(unit.Sp(windowWidth), unit.Sp(windowHeight)),
			gioapp.Title("GoDcr"),
		)

		if err := desktop.renderLoop(); err != nil {
			log.Fatal(err)
		}
	}()

	gioapp.Main()
}

func (d *Desktop) prepareHandlers() {
	list := &layout.List{
		Axis: layout.Vertical,
	}
	pages := getPages()
	d.pages = make([]page, len(pages))

	for index, page := range pages {
		d.pages[index] = page

		if index == 0 {
			d.changePage(page.name)
		}
	}
}

func (d *Desktop) changePage(pageName string) {
	if d.currentPage != nil && d.currentPage.name == pageName {
		return
	}

	if d.currentPage != nil && d.currentPage.name == pageName {
		return
	}

	for _, page := range d.pages {
		if page.name == pageName {
			d.currentPage = &page
			d.pageChanged = true
			break
		}
	}

}

func (d *Desktop) renderLoop() error {
	ctx := &layout.Context{
		Queue: d.window.Queue(),
	}

	for {
		e := <-d.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ctx.Reset(e.Config, e.Size)
			d.render(ctx)
			e.Frame(ctx.Ops)
		}
	}
}

func (d *Desktop) render(ctx *layout.Context) {
	inset := layout.Inset{
		Top:  unit.Dp(0),
		Left: unit.Dp(0),
	}

	inset.Layout(ctx, func() {
		flex := layout.Flex{
			Axis: layout.Horizontal,
		}

		navChild := flex.Rigid(ctx, func() {
			d.renderNavSection(ctx)
		})

		contentChild := flex.Rigid(ctx, func() {
			inset := layout.Inset{
				Left: unit.Sp(navSectionWidth - 103),
				Top:  unit.Dp(0),
			}

			inset.Layout(ctx, func() {
				d.renderContentSection(ctx)
			})
		})
		flex.Layout(ctx, navChild, contentChild)
	})
}

func (d *Desktop) renderNavSection(ctx *layout.Context) {
	var stack layout.Stack

	navAreaBounds := image.Point{
		X: navSectionWidth,
		Y: windowHeight * 2,
	}

	helper.PaintArea(ctx, helper.DecredDarkBlueColor, navAreaBounds)

	navArea := stack.Rigid(ctx, func() {
		inset := layout.Flex{Axis: layout.Vertical}
		inset.Layout(ctx, func() {
			//currentPositionTop := float32(0)
			for _, page := range d.pages {
				inset := layout.Flex{Axis: layout.Vertical}
				inset.Rigid(ctx, func() {
					for page.button.Clicked(ctx) {
						d.changePage(page.name)
					}
					d.theme.Button(page.label).Layout(ctx, page.button)
				})
				// inset.Layout(ctx, func() {

				// })
				// currentPositionTop += 32
			}
		})
	})
	stack.Layout(ctx, navArea)
}

func (d *Desktop) renderContentSection(ctx *layout.Context) {
	if d.pageChanged {
		d.pageChanged = false
		d.currentPage.handler.BeforeRender()
	}

	var stack layout.Stack

	contentAreaBounds := image.Point{
		X: windowWidth * 2,
		Y: windowHeight * 2,
	}

	helper.PaintArea(ctx, helper.BackgroundColor, contentAreaBounds)

	stack.Layout(ctx)

	/**
	stack := (&layout.Stack{})

	header := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top:  unit.Dp(0),
			Left: unit.Dp(0),
		}
		inset.Layout(ctx, func() {
			//widget.HeadingText(d.currentPage.label, widget.TextAlignLeft, ctx)
		})
	})

	headerLine := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(28),
			Left: unit.Dp(0),
		}

		inset.Layout(ctx, func(){
			/**bounds := image.Point{
				X: windowWidth - 30,
				Y: 1,
			}
			helper.PaintArea(helper.Theme.Brand, bounds, ctx.Ops)*
		})
	})

	stack.Layout(ctx, header, headerLine)**/
}
