package main

// functionality must include: changing tab, opening tab, closing tab, listening to buffer in active tab

type KeyboardClick string

type BrowserAutomater interface {
	StartSession() error
	ChangeTab(tabIndex int32) error
	OpenTab(url string) error
	CloseCurrentTab() error
	ListenToKeyboard() KeyboardClick
	EndSession() error
}

type SeleniumBrowserAutomater struct {
}
