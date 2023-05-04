package tui

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"strings"
)

type TUITitle struct {
}

func (t *TUITitle) ShowTitleAndDescription(title, description string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(description)
}

func (t *TUITitle) ShowTitle(title string) {
	titleNormalised := strings.TrimSpace(strings.ToUpper(title))
	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(titleNormalised)).
		Srender()
	pterm.DefaultCenter.Println(s)
}

func (t *TUITitle) ShowDescription(description string) {
	subtitleNormalised := strings.TrimSpace(strings.ToUpper(description))
	pterm.Println()
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("--------------------------------")
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(subtitleNormalised)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println("--------------------------------")
	pterm.Println()
}

func (t *TUITitle) ShowSubTitle(mainTitle string, subTitle string) {
	_ = pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle(strings.ToUpper(mainTitle), pterm.NewStyle(pterm.FgCyan)),
		putils.LettersFromStringWithStyle(strings.ToUpper(subTitle), pterm.NewStyle(pterm.FgLightMagenta))).
		Render()
}

func NewTitle() TUIDisplayer {
	return &TUITitle{}
}
