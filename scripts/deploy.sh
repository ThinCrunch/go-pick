#!/bin/bash

rm go-pick
go build -o go-pick main.go bundled.go

rm ~/.local/bin/go-pick
cp ./go-pick ~/.local/bin/go-pick