# GoPick
A simple solution to opening links when using multiple browsers, written in Go.

<img width="443" height="249" alt="image" src="https://github.com/user-attachments/assets/ae3050a2-d914-4c72-99c9-96606f0dbfc4" />

## Installation
1. make sure Go is installed
2. clone the repository and navigate to it
4. run `scripts/deploy.sh` to build and install

## Usage
After installing, the program is now available as a `go-pick` command, and as a desktop entry `GoPick`.

To open a link, simply pass it as an argument to the program (e.g. `go-pick "www.example.com"`).

If you chose to set it as default in the installer, you can just click a link and try it out!

## Uninstall
Run `scripts/uninstall.sh`, set your default browser again (with your preferred commands, or manage `.config/mimeapps.list` manually). Then you can just delete the cloned directory.
