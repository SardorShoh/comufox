package comufox

import "github.com/playwright-community/playwright-go"

type Page struct {
	page playwright.Page
}

// get element by given locator path from page
func (page *Page) Element(locator string) (*Element, error) {
	el := page.page.Locator(locator)
	if el.Err() != nil {
		return nil, el.Err()
	}
	count, err := el.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrElementNotExists
	}
	return &Element{element: el.First()}, nil
}

// get elements by given locator path from page
func (page *Page) Elements(locator string) (*Elements, error) {
	els := page.page.Locator(locator)
	if els.Err() != nil {
		return nil, els.Err()
	}
	count, err := els.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrElementNotExists
	}
	all, err := els.All()
	if err != nil {
		return nil, err
	}
	elements := new(Elements)
	for _, el := range all {
		*elements = append(*elements, &Element{element: el})
	}
	return elements, nil
}
