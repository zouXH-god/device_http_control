rsrc -ico app.ico -o rsrc.syso
go build -ldflags="-s -w -H=windowsgui" -trimpath