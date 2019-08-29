package pages

import (
	"context"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

type menuPageData struct {
	peerConn  *widget.Label
	blkHeight *widget.Label
	// there might be situations we would want to get the particular opened tab
	tabs *widget.TabContainer
}

var menu menuPageData

func pageNotImplemented() fyne.CanvasObject {
	label := widget.NewLabelWithStyle("This page has not been implemented yet", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	return label
}

func menuPage(ctx context.Context, wallet godcrApp.WalletMiddleware, fyneApp fyne.App, window fyne.Window) fyne.CanvasObject {
	menu.peerConn = widget.NewLabel("")
	menu.blkHeight = widget.NewLabel("")

	menu.tabs = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", fyne.NewStaticResource("Overview", Png("overview.png")), overviewPage(wallet, fyneApp)),
		widget.NewTabItemWithIcon("History", fyne.NewStaticResource("History", Png("history.png")), pageNotImplemented()),
		widget.NewTabItemWithIcon("Send", fyne.NewStaticResource("Send", Png("send.png")), pageNotImplemented()),
		widget.NewTabItemWithIcon("Receive", fyne.NewStaticResource("Receive", Png("receive.png")), receivePage(wallet, window)),
		widget.NewTabItemWithIcon("Accounts", fyne.NewStaticResource("Accounts", Png("account.png")), pageNotImplemented()),
		widget.NewTabItemWithIcon("Staking", fyne.NewStaticResource("Staking", Png("stake.png")), stakingPage(wallet)),
		widget.NewTabItemWithIcon("More", fyne.NewStaticResource("More", Png("more.png")), morePage(wallet, fyneApp)),
		widget.NewTabItemWithIcon("Exit", fyne.NewStaticResource("Exit", Png("exit.png")), exit(ctx, fyneApp, window)))
	menu.tabs.SetTabLocation(widget.TabLocationLeading)

	// would update all labels for all pages every seconds, all objects to be updated should be placed here
	go func() {
		for {
			// update only when the user is on the page
			if menu.tabs.CurrentTabIndex() == 0 {
				overviewPageUpdates(wallet)
			} else if menu.tabs.CurrentTabIndex() == 3 {
				receivePageUpdates(wallet)
			} else if menu.tabs.CurrentTabIndex() == 4 {
				stakingPageReloadData(wallet)
			}
			statusUpdates(wallet)
			time.Sleep(time.Second * 1)
		}
	}()

	// where peerConn and blkHeight are the realtime status texts
	status := widget.NewVBox(menu.peerConn, menu.blkHeight)
	data := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, status, menu.tabs, nil), menu.tabs, status)

	return data
}
