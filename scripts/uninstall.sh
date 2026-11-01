#!/bin/bash

rm ~/.local/bin/go-pick
rm ~/.local/share/applications/go-pick.desktop
update-desktop-database ~/.local/share/applications/

echo "Uninstalled go-pick. Feel free to delete all files."
