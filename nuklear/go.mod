module github.com/raedahgroup/godcr/nuklear

go 1.13

require (
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802 // indirect
	github.com/aarzilli/nucular v0.0.0-20180724180927-1d42e3457cef
	github.com/atotto/clipboard v0.1.2
	github.com/decred/dcrd/dcrutil v1.4.0
	github.com/decred/slog v1.0.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/raedahgroup/dcrlibwallet v1.0.1-0.20190831020110-aad933e3f96d
	github.com/raedahgroup/godcr/app v0.0.0-20200107105444-bd23847c1453
	github.com/skip2/go-qrcode v0.0.0-20190110000554-dc11ecdae0a9
	golang.org/x/image v0.0.0-20190501045829-6d32002ffd75
	golang.org/x/mobile v0.0.0-20190318164015-6bd122906c08
)

replace github.com/raedahgroup/godcr/app => ../app
