package pages

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/skip2/go-qrcode"
)

const receivingDecredHint = "Each time you request a payment, a new \naddress is created to protect your privacy."
const pageTitleText = "Receive DCR"

var receiveHandler struct {
	generatedReceiveAddress string
	recieveAddressError     error

	generatedQrCode []byte
	qrcodeError     error

	wallet             *dcrlibwallet.LibWallet
	accountNumberError error

	accountNumber       uint32
	selectedAccountName string
}

func ReceivePageContent(dcrlw *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) fyne.CanvasObject {
	receiveHandler.wallet = dcrlw

	// error handler
	errorLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide()
	errorHandler := func(err string) {
		errorLabel.SetText(err)
		errorLabel.Show()
	}

	qrImage := widget.NewIcon(theme.FyneLogo())
	generatedAddressLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	generateAddress := func(generateNewAddress bool) {
		receiveHandler.accountNumber, receiveHandler.accountNumberError = receiveHandler.wallet.AccountNumber(receiveHandler.selectedAccountName)
		if receiveHandler.accountNumberError != nil {
			errorHandler(fmt.Sprintf("Error: %s", receiveHandler.accountNumberError.Error()))
			return
		}

		err := generateAddressAndQrCode(generateNewAddress)
		if err != nil {
			errorHandler(fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		widget.Refresh(generatedAddressLabel)
		generatedAddressLabel.SetText(receiveHandler.generatedReceiveAddress)

		// If there was a rectified error and user clicks the generate address button again, this hides the error text.
		if !errorLabel.Hidden {
			errorLabel.Hide()
		}

		qrImage.SetResource(fyne.NewStaticResource("", receiveHandler.generatedQrCode))
	}

	// get icons used on the page.
	icons, err := assets.GetIcons(assets.CollapseIcon, assets.ReceiveAccountIcon, assets.MoreIcon, assets.InfoIcon)

	// clickableInfoIcon holds receiving decred hint-text pop-up
	var clickableInfoIcon *widgets.ClickableIcon
	clickableInfoIcon = widgets.NewClickableIcon(icons[assets.InfoIcon], nil, func() {
		receivingDecredHintLabel := widget.NewLabelWithStyle(receivingDecredHint, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		gotItLabel := canvas.NewText("Got it", color.RGBA{41, 112, 255, 255})
		gotItLabel.TextStyle = fyne.TextStyle{Bold: true}
		gotItLabel.TextSize = 16

		var popup *widget.PopUp
		popup = widget.NewPopUp(widget.NewVBox(
			widgets.NewVSpacer(14),
			widget.NewHBox(widgets.NewHSpacer(14), widget.NewLabelWithStyle(pageTitleText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			widgets.NewVSpacer(5),
			widget.NewHBox(widgets.NewHSpacer(14), receivingDecredHintLabel, widgets.NewHSpacer(14)),
			widgets.NewVSpacer(10),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(gotItLabel), func() { popup.Hide() }), widgets.NewHSpacer(14)),
			widgets.NewVSpacer(14),
		), window.Canvas())

		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableInfoIcon).Add(fyne.NewPos(0, clickableInfoIcon.Size().Height)))
	})

	// generate new address pop-up
	var clickableMoreIcon *widgets.ClickableIcon
	clickableMoreIcon = widgets.NewClickableIcon(icons[assets.MoreIcon], nil, func() {
		var popup *widget.PopUp
		popup = widget.NewPopUp(widgets.NewClickableBox(widget.NewHBox(widget.NewLabel("Generate new address")), func() {
			generateAddress(true)
			popup.Hide()
		}), window.Canvas())

		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(clickableMoreIcon).Add(fyne.NewPos(0, clickableMoreIcon.Size().Height)))
		popup.Show()
	})

	// get user accounts
	accounts, err := receiveHandler.wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel(fmt.Sprintf("Error: %s", err.Error()))
	}

	// automatically generate address for first account
	receiveHandler.selectedAccountName = accounts.Acc[0].Name
	selectedAccountLabel := widget.NewLabel(receiveHandler.selectedAccountName)
	selectedAccountBalanceLabel := widget.NewLabel(dcrutil.Amount(accounts.Acc[0].TotalBalance).String())
	generateAddress(false)

	var accountSelectionPopup *widget.PopUp
	accountsBox := widget.NewVBox()

	for index, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		accountProperties := widget.NewHBox()
		accountProperties.Append(widgets.NewHSpacer(15))
		accountProperties.Append(widget.NewIcon(icons[assets.ReceiveAccountIcon]))
		accountProperties.Append(widgets.NewHSpacer(15))

		spendableLabel := canvas.NewText("Spendable", color.Black)
		spendableLabel.TextSize = 12
		spendableLabel.Alignment = fyne.TextAlignLeading

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), color.Black)
		spendableAmountLabel.TextSize = 12
		spendableAmountLabel.Alignment = fyne.TextAlignCenter

		accountName := account.Name
		accountProperties.Append(widget.NewVBox(
			widget.NewLabel(accountName),
			spendableLabel,
		))

		accountbalance := dcrutil.Amount(account.Balance.Total).String()
		accountProperties.Append(widgets.NewHSpacer(84))
		accountProperties.Append(widget.NewVBox(
			widget.NewLabel(accountbalance),
			spendableAmountLabel,
		))

		accountProperties.Append(widgets.NewHSpacer(8))

		checkmarkIcon := widget.NewIcon(theme.ConfirmIcon())
		if index != 0 {
			checkmarkIcon.Hide()
		}

		accountProperties.Append(checkmarkIcon)

		accountsBox.Append(widgets.NewClickableBox(accountProperties, func() {
			// hide checkmark icon of other accounts
			for _, children := range accountsBox.Children {
				if box, ok := children.(*widgets.ClickableBox); !ok {
					continue
				} else {
					if len(box.Children) != 8 {
						continue
					}

					if icon, ok := box.Children[7].(*widget.Icon); !ok {
						continue
					} else {
						icon.Hide()
					}
				}
			}

			selectedAccountLabel.SetText(accountName)
			selectedAccountBalanceLabel.SetText(accountbalance)
			receiveHandler.selectedAccountName = accountName
			generateAddress(false)
			accountSelectionPopup.Hide()
		}))
	}

	hidePopUp := widget.NewHBox(widgets.NewHSpacer(16),
		widgets.NewClickableIcon(theme.CancelIcon(), nil, func() { accountSelectionPopup.Hide() }),
		widgets.NewHSpacer(16), widget.NewLabelWithStyle("Receiving account", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer())

	// accountSelectionPopup create a popup that has account names with spendable amount
	accountSelectionPopup = widget.NewPopUp(widget.NewVBox(
		widgets.NewVSpacer(5), hidePopUp, widgets.NewVSpacer(5), canvas.NewLine(color.Black), accountsBox),
		window.Canvas())
	accountSelectionPopup.Hide()

	// accountTab shows the selected account
	accountTab := widget.NewHBox(widget.NewIcon(icons[assets.ReceiveAccountIcon]), widgets.NewHSpacer(16),
		selectedAccountLabel, widgets.NewHSpacer(76), selectedAccountBalanceLabel, widgets.NewHSpacer(8), widget.NewIcon(icons[assets.CollapseIcon]))

	var accountDropdown *widgets.ClickableBox
	accountDropdown = widgets.NewClickableBox(accountTab, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountDropdown).Add(fyne.NewPos(0, accountDropdown.Size().Height)))
		accountSelectionPopup.Show()
	})

	accountCopiedLabel := widget.NewLabelWithStyle("Address copied", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	accountCopiedLabel.Hide()

	// copyAddressAction enables address copying
	copyAddressAction := widgets.NewClickableBox(widget.NewVBox(widget.NewLabelWithStyle("(Tap to copy)", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})),
		func() {
			clipboard := window.Clipboard()
			clipboard.SetContent(receiveHandler.generatedReceiveAddress)

			accountCopiedLabel.Show()
			widget.Refresh(accountCopiedLabel)

			// only hide accountCopiedLabel text if user is currently on the page after 2secs
			if accountCopiedLabel.Hidden == false {
				time.AfterFunc(time.Second*2, func() {
					if tabmenu.CurrentTabIndex() == 3 {
						accountCopiedLabel.Hide()
					}
				})
			}
		},
	)

	pageTitleLabel := widget.NewLabelWithStyle(pageTitleText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	output := widget.NewVBox(
		widgets.NewVSpacer(5),
		widget.NewHBox(pageTitleLabel, widgets.NewHSpacer(110), clickableInfoIcon, widgets.NewHSpacer(26), clickableMoreIcon),
		accountDropdown,
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(70), accountCopiedLabel, layout.NewSpacer()),
		// due to original width content on page being small layout spacing isn't efficient
		// therefore requiring an additional spacing by a 100 width
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(70), fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(300, 300)), qrImage), layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(70), generatedAddressLabel, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(70), copyAddressAction, layout.NewSpacer()),
		errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(18), output)
}

func generateAddressAndQrCode(generateNewAddress bool) error {
	if generateNewAddress {
		generatedReceiveAddress, err := receiveHandler.wallet.NextAddress(int32(receiveHandler.accountNumber))
		if err != nil {
			return err
		}
		receiveHandler.generatedReceiveAddress = generatedReceiveAddress
	} else {
		generatedReceiveAddress, err := receiveHandler.wallet.CurrentAddress(int32(receiveHandler.accountNumber))
		if err != nil {
			return err
		}
		receiveHandler.generatedReceiveAddress = generatedReceiveAddress
	}

	generatedQrCode, err := generateQrcode(receiveHandler.generatedReceiveAddress)
	if err != nil {
		return err
	}

	receiveHandler.generatedQrCode = generatedQrCode
	return nil
}

func generateQrcode(generatedReceiveAddress string) ([]byte, error) {
	// generate qrcode
	png, err := qrcode.Encode(generatedReceiveAddress, qrcode.High, 256)
	if err != nil {
		return nil, err
	}

	return png, nil
}
