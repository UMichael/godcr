package pages

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/raedahgroup/dcrlibwallet/txhelper"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageData struct {
	fromAccountSelect *widget.Select
	toAccountSelect   *widget.Select
	errorLabel        *widget.Label
}

//both send page and send page update would be in a function
var send sendPageData

func sendPageUpdates(wallet godcrApp.WalletMiddleware) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		send.errorLabel.Text = err.Error()
		send.errorLabel.Show()
		canvas.Refresh(send.errorLabel)
		return
	}
	send.errorLabel.Text = ""
	send.errorLabel.Hide()

	var fullStatus []string
	for _, account := range accounts {
		fullStatus = append(fullStatus, account.String())
	}
	send.fromAccountSelect.Options = fullStatus
	send.toAccountSelect.Options = fullStatus
}

func sendPage(wallet godcrApp.WalletMiddleware, window fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Sending Decred", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	accountLabel := widget.NewLabelWithStyle("From:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	send.errorLabel = widget.NewLabel("")
	widget.Refresh(send.errorLabel)
	send.errorLabel.Hide()

	button := widget.NewButton("Submit", func() {
		fmt.Println("Hello")
	})
	button.Disabled()

	send.fromAccountSelect = widget.NewSelect(nil, func(s string) {
		if button.Disabled() == true {
			button.Enable()
		}
	})
	send.toAccountSelect = widget.NewSelect(nil, nil)
	widget.Refresh(send.toAccountSelect)
	send.toAccountSelect.Hide()
	sendPageUpdates(wallet)
	if send.errorLabel.Text == "" {
		send.toAccountSelect.SetSelected(send.toAccountSelect.Options[0])
		send.fromAccountSelect.SetSelected(send.fromAccountSelect.Options[0])
	}

	//place a menu toolbar to view options
	menuPng, err := ioutil.ReadFile("fyne/pages/png/menu.png")
	if err != nil {
		log.Fatalln("could not read file menu.png", err.Error())
	}
	sendMaxPng, err := ioutil.ReadFile("fyne/pages/png/sendMax.png")
	if err != nil {
		log.Fatalln("could not read file sendMax.png", err.Error())
	}

	fee := widget.NewLabel("")
	estimateSize := widget.NewLabel("")
	balanceAfter := widget.NewLabel("")

	sendProperties := widget.NewForm()

	sendProperties.Append("Fee:", fee)
	sendProperties.Append("Estimate Size:", estimateSize)
	sendProperties.Append("Balance After:", balanceAfter)

	address := widget.NewEntry()
	address.SetPlaceHolder("Destination Address")
	amount := widget.NewEntry()
	amount.SetPlaceHolder("Amount")
	amount.OnChanged = func(text string) {
		amnt, err := strconv.ParseFloat(amount.Text, 64)
		if err != nil {
			return
		}
		splittedWord := strings.Split(send.fromAccountSelect.Selected, " ")
		acctNo, err := wallet.AccountNumber(splittedWord[0])
		if err != nil {
			return
		}

		balance, err := wallet.AccountBalance(acctNo, walletcore.DefaultRequiredConfirmations)
		if err != nil {
			return
		}
		if (balance.Spendable.ToCoin() - amnt) <= 0 {
			send.errorLabel.SetText("Not enough funds (or not connected)")
			return
		}
		txData := txhelper.TransactionDestination{
			Address: address.Text,
			Amount:  amnt,
		}
		output, err := txhelper.MakeTxOutput(txData)
		if err != nil {
			return
		}
		val, ok := interface{}(sendProperties.Items[0].Widget).(*widget.Label)
		fmt.Println(ok)
		val.SetText(strconv.Itoa(int(output.Value)))
		val, ok = interface{}(sendProperties.Items[1].Widget).(*widget.Label)
		fmt.Println(ok)
		val.SetText(strconv.Itoa(output.SerializeSize()))
		val, ok = interface{}(sendProperties.Items[2].Widget).(*widget.Label)
		fmt.Println(ok)
		val.SetText(strconv.FormatFloat(balance.Spendable.ToCoin()-amnt, 'f', 8, 64))
	}

	sendMaxButton := widget.NewButtonWithIcon("", fyne.NewStaticResource("send max", sendMaxPng), func() {

		splittedWord := strings.Split(send.fromAccountSelect.Selected, " ")
		acctNo, err := wallet.AccountNumber(splittedWord[0])
		if err != nil {
			return
		}

		balance, err := wallet.AccountBalance(acctNo, walletcore.DefaultRequiredConfirmations)
		if err != nil {
			return
		}
		amount.SetText(strconv.FormatFloat(balance.Spendable.ToCoin(), 'f', 8, 64))
	})
	amountContainer := widget.NewHBox(amount, sendMaxButton)

	var popup *widget.PopUp
	menuIcon := widget.NewToolbar(widget.NewToolbarAction(fyne.NewStaticResource("menu", menuPng), func() { popup.Show() }))
	transactionInfo := widget.NewHBox(address, menuIcon)

	popup = widget.NewPopUp(widget.NewVBox(
		widget.NewButton("Send between accounts", func() {
			address.Hide()
			send.toAccountSelect.Show()
			transactionInfo.Children = []fyne.CanvasObject{send.toAccountSelect, menuIcon}
			widget.Refresh(transactionInfo)
			popup.Hide()
		}),
		widget.NewButton("Send to others", func() {
			send.toAccountSelect.Hide()
			address.Show()
			transactionInfo.Children = []fyne.CanvasObject{address, menuIcon}
			widget.Refresh(transactionInfo)
			popup.Hide()
		}),
		widget.NewButton("Clear all fields", func() {
			address.SetText("")
			amount.SetText("")
			popup.Hide()
		}),
	), window.Canvas())

	pos := menuIcon.Position()
	mousePos := fyne.NewPos(pos.X+50, pos.Y+50)
	popup.Move(mousePos)
	popup.Hide()

	infoIcon := widget.NewToolbar(widget.NewToolbarAction(theme.InfoIcon(), func() {
		var popUp *widget.PopUp
		button := widget.NewButton("Got it", func() {
			popUp.Hide()
		})

		header := widget.NewLabelWithStyle("Send DCR", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		info := widget.NewLabelWithStyle("Input the destination wallet address and the amount to send funds", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
		data := widget.NewVBox(header, widgets.NewVSpacer(10), info,
			widget.NewHBox(layout.NewSpacer(), button))
		popUp = widget.NewModalPopUp(data, window.Canvas())
	}))

	output := widget.NewVBox(
		widget.NewHBox(label, infoIcon),
		widgets.NewVSpacer(15),
		widget.NewHBox(accountLabel, send.fromAccountSelect),
		widgets.NewVSpacer(15),
		transactionInfo,
		widgets.NewVSpacer(15),
		amountContainer,
		widgets.NewVSpacer(15),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(button.MinSize()), button),
		sendProperties,
		send.errorLabel,
	)
	return widget.NewHBox(widgets.NewHSpacer(10), output)
}
