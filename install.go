package comufox

import (
	"github.com/SardorShoh/comufox/installer"
	"github.com/playwright-community/playwright-go"
)

func Install() error {
	if err := installer.InstallCamoufox(); err != nil {
		return err
	}
	if err := playwright.Install(); err != nil {
		return err
	}
	return nil
}
