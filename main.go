package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const IconsPath string = "/usr/share/icons/hicolor/32x32/apps/"
const MimeInfoPath string = "/usr/share/applications/mimeinfo.cache"
const ApplicationsPath string = "/usr/share/applications/"

const WindowPixelWidth int = 400

var BrowserBlacklist = []string{"go-pick.desktop"}

func main() {
	url := ""
	if len(os.Args) >= 2 {
		url = os.Args[1]
	}

	a := app.New()
	w := a.NewWindow("GoPick")
	w.SetFixedSize(true)

	// Disables scaling change when moving between monitors with different DPIs
	err := os.Setenv("FYNE_DISABLE_DPI_DETECTION", "1")
	if err != nil {
		log.Printf("Unable to set environmental variable %s", err)
	}

	windowContent, err := createWindowContent(a, url)
	if err != nil {
		log.Fatalf("Unable to retrieve necessary information. Nothing to show. %s", err)
	}

	w.SetContent(windowContent)

	w.Resize(fyne.NewSize(float32(WindowPixelWidth), windowContent.Size().Height))
	w.CenterOnScreen()

	w.ShowAndRun()
}

func createWindowContent(app fyne.App, url string) (*fyne.Container, error) {
	browsers, err := getBrowsers()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve installed browsers %w", err)
	}

	browserButtons := make([]fyne.CanvasObject, 0)

	for _, browser := range browsers {
		button, err := createBrowserButton(browser, app, url)
		if err != nil {
			log.Printf("Unable to create button for %s %s", browser, err)
			continue
		}

		browserButtons = append(browserButtons, button)
	}

	pickLabel := widget.NewLabel("Pick a browser to open:")
	wrappedUrlLabel := wrapStringToWidth(url, WindowPixelWidth)
	urlLabel := widget.NewLabel(wrappedUrlLabel)

	windowContent := container.NewVBox(pickLabel)
	windowContent.Add(urlLabel)
	for _, button := range browserButtons {
		windowContent.Add(button)
	}

	return windowContent, nil
}

func wrapStringToWidth(text string, maxPixelWidth int) string {
	sBuilder := strings.Builder{}

	fontSize := theme.TextSize()
	fontStyle := fyne.TextStyle{}

	lastWrapIdx := 0
	for i := 1; i <= len(text); i++ {
		currentTextWidth := fyne.MeasureText(text[lastWrapIdx:i], fontSize, fontStyle).Width
		if currentTextWidth > float32(maxPixelWidth) {
			sBuilder.WriteString(text[lastWrapIdx:i])
			sBuilder.WriteString("\n")
			lastWrapIdx = i
		}
	}

	if lastWrapIdx < len(text) {
		sBuilder.WriteString(text[lastWrapIdx:])
	}

	return sBuilder.String()
}

func createBrowserButton(browser string, app fyne.App, url string) (*widget.Button, error) {
	name, iconName := getBrowserNameIcon(browser)
	icon := getBrowserIconResource(iconName)

	button := widget.NewButtonWithIcon(name, icon, func() {
		err := openBrowser(browser, url)
		if err != nil {
			log.Printf("Unable to open browser. %s", err)
		}
		app.Quit()
	})
	button.Alignment = widget.ButtonAlignLeading

	return button, nil
}

func getBrowserIconResource(iconName string) fyne.Resource {
	if iconName == "" {
		return resourceQuestionMarkPng
	}

	iconPath := getBrowserIconPath(iconName)
	icon, err := fyne.LoadResourceFromPath(iconPath)
	if err != nil {
		log.Printf("Unable to load icon %s\n%s", iconName, err)
		return resourceQuestionMarkPng
	}
	return icon
}

// Launches the browser based on the passed desktop file.
func openBrowser(desktopFileName string, url string) error {
	cmd := exec.Command("gtk-launch", desktopFileName, url)
	if url == "" {
		cmd = exec.Command("gtk-launch", desktopFileName)
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to run command. %w", err)
	}
	return nil
}

// Retrieves value from a file, that contains KeyValue pairs.
func getValueFromFile(filePath string, key string) (string, error) {
	types, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open file %s %w", filePath, err)
	}

	for _, line := range strings.Split(string(types), "\n") {
		if strings.HasPrefix(line, key) {
			return strings.Split(line, "=")[1], nil
		}
	}
	return "", errors.New("Key " + key + " not found.")
}

// Retrieves all applications capable of handling http and https links.
// Exceptions are set in <BrowserBlacklist>.
func getBrowsers() ([]string, error) {
	httpTypes := []string{"x-scheme-handler/http", "x-scheme-handler/https"}

	browsers := make([]string, 0)
	for _, type_ := range httpTypes {
		programsString, err := getValueFromFile(MimeInfoPath, type_)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve browsers. %w", err)
		}

		lineBrowsers := strings.Split(programsString, ";")
		for _, browser := range lineBrowsers {
			if !strings.HasSuffix(browser, ".desktop") {
				continue
			}

			if slices.Contains(browsers, browser) {
				continue
			}

			if slices.Contains(BrowserBlacklist, browser) {
				continue
			}

			browsers = append(browsers, browser)
		}
	}

	return browsers, nil
}

func getBrowserNameIcon(desktopFile string) (string, string) {
	file := ApplicationsPath + desktopFile
	icon, iconErr := getValueFromFile(file, "Icon")
	name, nameErr := getValueFromFile(file, "Name")

	if iconErr != nil {
		icon = ""
		log.Printf("Unable to get icon for %s: %s", desktopFile, iconErr)
	}

	if nameErr != nil {
		name = desktopFile
		log.Printf("Unable to get name for %s: %s", desktopFile, nameErr)
	}

	return name, icon
}

func getBrowserIconPath(iconName string) string {
	return IconsPath + iconName + ".png"
}
