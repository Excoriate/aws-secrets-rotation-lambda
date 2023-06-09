package tui

type TUIMessenger interface {
	ShowError(title, msg string, err error)
	ShowInfo(title, msg string)
	ShowSuccess(title, msg string)
	ShowWarning(title, msg string)
}

type TUIDisplayer interface {
	ShowTitleAndDescription(title, description string)
	ShowTitle(title string)
	ShowSubTitle(mainTitle, subtitle string)
}
