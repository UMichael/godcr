package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/skip2/go-qrcode"
)

type receivePageData struct {
	accountSelect *widget.Select
	errorLabel    *widget.Label
}

var receive receivePageData

// todo: remove this when account page is implemented
func receivePageUpdates(wallet godcrApp.WalletMiddleware) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		receive.errorLabel.Text = err.Error()
		widget.Refresh(receive.errorLabel)
		receive.errorLabel.Show()
		return
	}

	var options []string
	for _, account := range accounts {
		options = append(options, account.Name)
	}
	receive.accountSelect.Options = options
	widget.Refresh(receive.accountSelect)
}

func receivePage(wallet godcrApp.WalletMiddleware, window fyne.Window) fyne.CanvasObject {
	// if there were to be situations, wallet fails and new address cant be generated, then simply show fyne logo
	qrImage := canvas.NewImageFromResource(theme.FyneLogo())
	qrImage.SetMinSize(fyne.NewSize(300, 300))

	label := widget.NewLabelWithStyle("Receiving Funds", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	info := widget.NewLabelWithStyle("Each time you request a payment, a new address is created to protect your privacy.", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true, Italic: true})
	accountLabel := widget.NewLabelWithStyle("Account:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	generatedAddress := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	receive.errorLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	widget.Refresh(receive.errorLabel)
	receive.errorLabel.Hide()

	var addr string
	copy := widget.NewToolbar(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		clipboard := window.Clipboard()
		clipboard.SetContent(addr)
	}))

	generateAddressFunc := func() {
		name, err := wallet.AccountNumber(receive.accountSelect.Selected)
		if err != nil {
			receive.errorLabel.SetText("error getting account name, " + err.Error())
			receive.errorLabel.Show()
			return
		}

		addr, err = wallet.GenerateNewAddress(name)
		if err != nil {
			receive.errorLabel.SetText("could not generate new address, " + err.Error())
			receive.errorLabel.Show()
			return
		}
		// if there was a rectified error and user clicks the generate again, this hides the error text
		if receive.errorLabel.Hidden == false {
			receive.errorLabel.Text = ""
			receive.errorLabel.Hide()
		}

		generatedAddress.Text = (addr)
		widget.Refresh(generatedAddress)

		png, _ := qrcode.Encode(addr, qrcode.High, 256)
		qrImage.Resource = fyne.NewStaticResource("Address", png)
		qrImage.Show()
		canvas.Refresh(qrImage)
	}

	button := widget.NewButton("Generate Address", func() {
		generateAddressFunc()
	})
	receive.accountSelect = widget.NewSelect(nil, nil)
	receivePageUpdates(wallet)
	if receive.errorLabel.Text == "" {
		receive.accountSelect.SetSelected(receive.accountSelect.Options[0])
	}
	generateAddressFunc()

	output := widget.NewVBox(
		label,
		info,
		widget.NewHBox(accountLabel, receive.accountSelect),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(button.MinSize()), button),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), qrImage, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), generatedAddress, copy, layout.NewSpacer()),
		receive.errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}
