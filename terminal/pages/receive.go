package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
	qrcode "github.com/skip2/go-qrcode"
)

func receivePage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(primitives.TitleTextView("Generate Receive Address"), 1, 0, false)

	hintText := primitives.WordWrappedTextView("(TIP: Navigate with Tab and Shift+Tab, hit ENTER to generate Address. Return with Esc)")
	hintText.SetTextColor(tcell.ColorGray)
	body.AddItem(hintText, 3, 0, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}
	if len(accounts) == 1 {
		address, qr, err := generateAddress(wallet, accounts[0].Number)
		if err != nil {
			return body.AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
		}
		body.AddItem(primitives.NewLeftAlignedTextView(fmt.Sprintf("Address: %s", address)).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				clearFocus()
			}
		}), 2, 1, true).
			AddItem(primitives.NewLeftAlignedTextView(fmt.Sprintf(qr.ToSmallString(false))).SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEscape {
					clearFocus()
				}
			}), 0, 1, true)
	} else {
		form := tview.NewForm()
		form.SetBorderPadding(0, 0, 0, 0)

		var accountNum uint32
		accountN := make([]uint32, len(accounts))
		accountNames := make([]string, len(accounts))
		for index, account := range accounts {
			accountNames[index] = account.Name
			body.AddItem(form.AddDropDown("Account", accountNames, 0, func(option string, optionIndex int) {
				accountNum = accountN[optionIndex]
			}).
				AddButton("Generate", func() {
					address, qr, err := generateAddress(wallet, accountNum)
					if err != nil {
						body.AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf("Error: %s", err.Error())), 3, 1, false)
						return
					}
					body.AddItem(primitives.NewLeftAlignedTextView(fmt.Sprintf("Address: %s", address)), 2, 1, false).
						AddItem(primitives.NewLeftAlignedTextView(fmt.Sprintf(qr.ToSmallString(false))), 0, 1, false)
				}).SetLabelColor(tcell.ColorWhite).SetItemPadding(15).SetHorizontal(true).SetCancelFunc(func() {
				clearFocus()
			}), 4, 1, true)
		}
	}

	body.SetBorderPadding(1, 0, 1, 0)

	setFocus(body)
	return body
}

func generateAddress(wallet walletcore.Wallet, accountNumber uint32) (string, *qrcode.QRCode, error) {
	generatedAddress, err := wallet.ReceiveAddress(accountNumber)
	if err != nil {
		return "", nil, err
	}

	// generate qrcode
	qr, err := qrcode.New(generatedAddress, qrcode.Medium)
	if err != nil {
		return "", nil, err
	}

	return generatedAddress, qr, nil
}
