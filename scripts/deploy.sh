#!/bin/bash

echo "Building..."

rm go-pick
go build -o go-pick main.go bundled.go

rm ~/.local/bin/go-pick
cp ./go-pick ~/.local/bin/go-pick

rm ~/.local/share/applications/go-pick.desktop
cp ./packaging/go-pick.desktop ~/.local/share/applications/go-pick.desktop
update-desktop-database ~/.local/share/applications/

echo "Installation done!"

read -p "Set GoPick as default for opening links? (y/n) " decision
if [[ $decision == y || $decision == Y ]]; then
    xdg-mime default go-pick.desktop x-scheme-handler/http
    xdg-mime default go-pick.desktop x-scheme-handler/https
fi
