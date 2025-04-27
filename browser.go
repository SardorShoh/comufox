package comufox

import "github.com/playwright-community/playwright-go"

type Browser struct {
	browser playwright.Browser
	ctx     playwright.BrowserContext
	page    Page
}

func New() (*Browser, error) {
	return nil, nil
}
