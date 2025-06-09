package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/sqweek/dialog"
)

func main() {
	app := app.New()
	window := app.NewWindow("File Encryptor")
	//window.Resize(fyne.NewSize(300, 200))
	window.SetFixedSize(true)

	fileSelectedLabel := widget.NewLabel("No file selected")
	//fileSelectedLabel.Wrapping = fyne.TextWrap(fyne.TextTruncateEllipsis)

	fileSelectButton := widget.NewButton("Select file", func() {
		filePath, err := dialog.File().Load()
		if err != nil && err.Error() != "Cancelled" {
			dialog.Message("Error seleting file: %s", err).Error()
		}
		if filePath != "" {
			fileSelectedLabel.SetText(filePath)
		}
	})

	fileSelectorContainer := container.NewHBox(fileSelectedLabel, fileSelectButton)

	passwordEntry := widget.NewPasswordEntry()

	encryptButton := widget.NewButton("Encrypt", func() {
		if passwordEntry.Text != "" && fileSelectedLabel.Text != "No file selected" {
			data, err := os.ReadFile(fileSelectedLabel.Text)
			if err != nil {
				dialog.Message("Error opening selected file: %s", err).Error()
			}

			password := passwordEntry.Text

			salt := generateSalt(16)
			key := deriveKey(password, salt, 2, (64 * 1024), 4, 32)

			armoredBytes := encryptBytes(key, data)
			armoredBytesWithAppendedSalt := append(armoredBytes, salt...)
			log.Printf("\n%s", armoredBytes)
			log.Printf("\n%s", armoredBytesWithAppendedSalt[:len(armoredBytesWithAppendedSalt)-16])

			filepath, err := dialog.File().Title("Save encrypted file").Save()
			if err != nil {
				if err.Error() == "Cancelled" {
					return
				} else {
					dialog.Message("Error reading path to save encrypted file: %s", err).Error()
				}
			}

			err = os.WriteFile(filepath, armoredBytesWithAppendedSalt, os.ModePerm)
			if err != nil {
				dialog.Message("Error writing encrypted file: %s", err).Error()
			}
		} else {
			dialog.Message("Select a file and make a password to encrypt").Error()
		}

	})

	decryptButton := widget.NewButton("Decrypt", func() {
		if passwordEntry.Text != "" && fileSelectedLabel.Text != "No file selected" {
			data, err := os.ReadFile(fileSelectedLabel.Text)
			if err != nil {
				dialog.Message("Error opening selected file: %s", err).Error()
			}

			password := passwordEntry.Text
			salt := data[len(data)-16:]
			key := deriveKey(password, salt, 2, (64 * 1024), 4, 32)

			decryptedBytes := decryptBytes(key, data[:len(data)-16])
			log.Printf("\n%s", data[:len(data)-16])

			filepath, err := dialog.File().Title("Save decrypted file").Save()
			if err.Error() == "Cancelled" {
				return
			} else {
				dialog.Message("Error reading path to save decrypted file: %s", err).Error()
			}

			err = os.WriteFile(filepath, decryptedBytes, os.ModePerm)
			if err != nil {
				dialog.Message("Error writing encrypted file: %s", err).Error()
			}
		} else {
			dialog.Message("Select a file and make a password to decrypt").Error()
		}
	})

	encryptDecryptContainer := container.New(layout.NewGridLayoutWithColumns(2), encryptButton, decryptButton)

	content := container.New(layout.NewGridLayoutWithRows(3), fileSelectorContainer, passwordEntry, encryptDecryptContainer)
	centeredContent := container.NewCenter(content)
	window.SetContent(container.NewPadded(centeredContent))

	window.ShowAndRun()
}
