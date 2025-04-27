package comufox

import (
	"errors"

	"github.com/playwright-community/playwright-go"
)

var (
	ErrElementNotExists  = errors.New("element not exists")
	ErrElementsNotExists = errors.New("elements not exists")
)

type (
	Element struct {
		element playwright.Locator
	}

	Elements []*Element
)

// get element from inside given element by locator
func (el *Element) Element(locator string, retry ...int) (*Element, error) {
	element := el.element.Locator(locator)
	if element.Err() != nil {
		return nil, element.Err()
	}
	count, err := element.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrElementNotExists
	}
	return &Element{element: element.First()}, nil
}

// get elements from inside given element by locator
func (el *Element) Elements(locator string) (*Elements, error) {
	els := el.element.Locator(locator)
	if els.Err() != nil {
		return nil, els.Err()
	}
	count, err := els.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrElementsNotExists
	}
	all, err := els.All()
	if err != nil {
		return nil, err
	}
	elements := new(Elements)
	for _, e := range all {
		*elements = append(*elements, &Element{element: e})
	}
	return elements, nil
}
