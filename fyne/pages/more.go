package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

func morePage(wallet godcrApp.WalletMiddleware, fyneApp fyne.App) fyne.CanvasObject {
	container := widget.NewTabContainer(
		widget.NewTabItemWithIcon("Settings", fyne.NewStaticResource("More", Png("settings.png")), settingsPage(fyneApp)),
		widget.NewTabItemWithIcon("Security Tools", fyne.NewStaticResource("More", Png("security.png")), pageNotImplemented()),
		widget.NewTabItemWithIcon("Help", fyne.NewStaticResource("More", Png("help.png")), pageNotImplemented()),
		widget.NewTabItemWithIcon("Debug", fyne.NewStaticResource("More", Png("settings.png")), pageNotImplemented()), // TODO - needs its own PNG ?
		widget.NewTabItemWithIcon("About", fyne.NewStaticResource("More", Png("about.png")), pageNotImplemented()))
	return container
}
