package design

import webappdesign "github.com/dinnerdonebetter/backend/internal/platform/webapp/design"

type Color = webappdesign.Color
type Palette = webappdesign.Palette

var StandardPalette = webappdesign.StandardPalette

func Background(c Color) string  { return webappdesign.Background(c) }
func TextColor(c Color) string   { return webappdesign.TextColor(c) }
func BorderColor(c Color) string { return webappdesign.BorderColor(c) }
