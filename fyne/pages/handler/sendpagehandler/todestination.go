package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/fyne/widgets"

	"github.com/raedahgroup/godcr/fyne/pages/handler/constantvalues"
)

// sendingToDestinationComponents constitutes all components that composes sending coins to users or self.
func (sendPage *SendPageObjects) initToDestinationComponents() error {

	fromLabel := canvas.NewText("To", color.RGBA{61, 88, 115, 255})
	fromLabel.TextStyle.Bold = true

	accountBox, err := sendPage.SelfSending.CreateAccountSelector("Sending account")
	if err != nil {
		return err
	}
	accountBox.Hide()

	sendPage.destinationAddressEntryComponent()

	sendToAccountLabel := canvas.NewText(constantvalues.SwitchToSendToAccount, color.RGBA{R: 41, G: 112, B: 255, A: 255})
	sendToAccountLabel.TextSize = 12

	destinationAddressContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(
		fyne.NewSize(widget.NewLabel(constantvalues.TestAddress).MinSize().Width, sendPage.destinationAddressEntry.MinSize().Height)), sendPage.destinationAddressEntry)

	spacer := widgets.NewVSpacer(10)

	var container *fyne.Container
	switchingComponentButton := widgets.NewClickableBox(widget.NewVBox(sendToAccountLabel), func() {
		if accountBox.Hidden {
			sendToAccountLabel.Text = constantvalues.SwitchToSendToAddress
			accountBox.Show()
			sendPage.destinationAddressEntry.Hide()
			destinationAddressContainer.Hide()
			sendPage.destinationAddressErrorLabel.Hide()
			spacer.Hide()

			if sendPage.amountEntry.Text != "" {
				sendPage.nextButton.Enable()
			} else {
				sendPage.nextButton.Disable()
			}

		} else {
			sendToAccountLabel.Text = constantvalues.SwitchToSendToAccount
			sendPage.destinationAddressEntry.Show()
			destinationAddressContainer.Show()
			accountBox.Hide()
			spacer.Show()

			sendPage.destinationAddressEntry.OnChanged(sendPage.destinationAddressEntry.Text)
			if sendPage.amountEntry.Text != "" && sendPage.destinationAddressEntry.Text != "" && sendPage.destinationAddressErrorLabel.Hidden {
				sendPage.nextButton.Enable()
			} else {
				sendPage.nextButton.Disable()
			}
		}

		container.Refresh()
		sendPage.SendPageContents.Refresh()
		sendPage.amountEntry.OnChanged(sendPage.amountEntry.Text)
	})

	box := widget.NewVBox(
		widget.NewHBox(fromLabel, layout.NewSpacer(), switchingComponentButton, widgets.NewHSpacer(20)),
		accountBox,
		destinationAddressContainer,
		sendPage.destinationAddressErrorLabel,
		spacer)

	container = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(box.MinSize()), box)

	sendPage.SendPageContents.Append(container)

	return err
}

func (sendPage *SendPageObjects) destinationAddressEntryComponent() {

	sendPage.destinationAddressErrorLabel = canvas.NewText("", color.RGBA{237, 109, 71, 255})
	sendPage.destinationAddressErrorLabel.TextSize = 12
	sendPage.destinationAddressErrorLabel.Hide()

	sendPage.destinationAddressEntry = widget.NewEntry()
	sendPage.destinationAddressEntry.SetPlaceHolder(constantvalues.DestinationAddressPlaceHolder)

	sendPage.destinationAddressEntry.OnChanged = func(address string) {
		if sendPage.destinationAddressEntry.Text == "" {
			sendPage.destinationAddressErrorLabel.Hide()
			sendPage.SendPageContents.Refresh()
			return
		}

		_, err := dcrutil.DecodeAddress(address)
		if err != nil {
			sendPage.destinationAddressErrorLabel.Text = constantvalues.InvalidAddress
			sendPage.destinationAddressErrorLabel.Show()
			setLabelText(constantvalues.NilAmount, sendPage.transactionFeeLabel, sendPage.totalCostLabel, sendPage.balanceAfterSendLabel)
			setLabelText(constantvalues.ZeroByte, sendPage.transactionSizeLabel)

		} else {
			sendPage.destinationAddressErrorLabel.Hide()
		}

		if sendPage.amountEntry.Text != "" && sendPage.amountEntryErrorLabel.Hidden && sendPage.destinationAddressErrorLabel.Hidden {
			sendPage.nextButton.Enable()
		} else {
			sendPage.nextButton.Disable()
		}

		sendPage.SendPageContents.Refresh()
	}
}

func setLabelText(Text string, objects ...*widget.Label) {
	for _, object := range objects {
		object.SetText(Text)
		object.Refresh()
	}
}