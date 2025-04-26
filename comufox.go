package comufox

import (
	"github.com/SardorShoh/comufox/dirs"
	"github.com/SardorShoh/comufox/launch"

	"github.com/playwright-community/playwright-go"
)

type ComofuxOptions struct {
	LaunchOptions  *playwright.BrowserTypeLaunchOptions
	ContextOptions *playwright.BrowserNewContextOptions
}

type Comofux struct {
	browser playwright.Browser
	ctx     playwright.BrowserContext
}

func Launch(pw *playwright.Playwright, options ...ComofuxOptions) (*Comofux, error) {
	if pw == nil {
		plw, err := playwright.Run()
		if err != nil {
			return nil, err
		}
		pw = plw
	}
	opt := ComofuxOptions{}
	if len(options) > 0 {
		opt = options[0]
	}
	if opt.LaunchOptions == nil {
		opt.LaunchOptions = &playwright.BrowserTypeLaunchOptions{}
	}
	pathName := dirs.GetExecutableName()
	if err := launch.SetExecutablePermissions(pathName); err != nil {
		return nil, err
	}
	opt.LaunchOptions.ExecutablePath = playwright.String(pathName)
	browser, err := pw.Firefox.Launch(*opt.LaunchOptions)
	if err != nil {
		return nil, err
	}
	var ctx playwright.BrowserContext
	if opt.ContextOptions != nil {
		ctx, err = browser.NewContext(*opt.ContextOptions)
	} else {
		ctx, err = browser.NewContext()
	}
	if err != nil {
		return nil, err
	}
	return &Comofux{browser: browser, ctx: ctx}, nil
}

func (c *Comofux) Close() error {
	if err := c.ctx.Close(); err != nil {
		return err
	}
	return c.browser.Close()
}

func (c *Comofux) Browser() playwright.Browser        { return c.browser }
func (c *Comofux) Context() playwright.BrowserContext { return c.ctx }
func (c *Comofux) NewPage(url string, opts ...playwright.PageGotoOptions) (playwright.Page, error) {
	page, err := c.ctx.NewPage()
	if err != nil {
		return nil, err
	}
	if url != "" {
		if _, err := page.Goto(url, opts...); err != nil {
			return nil, err
		}
	}
	return page, nil
}
