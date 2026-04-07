package crayon

// =============================
// COLOR TOGGLE
// =============================
// Should auto detect tty by default
func NewColorToggle(enableColor ...bool) *ColorToggle {
	var colorEnabled bool
	if len(enableColor) > 0 {
		colorEnabled = enableColor[0]
	} else {
		colorEnabled = autoDetect()
	}
	return &ColorToggle{
		EnableColor: colorEnabled,
	}
}

//=============================
// PARSE
//=============================

// Parse template string using auto-detected color support
// User don't need to explicitly define color toggle
func Parse(input string) CompiledTemplate {
	return NewColorToggle().Parse(input)
}

// Parse template string using the toggle's color setting
func (toggle *ColorToggle) Parse(input string) CompiledTemplate {
	if toggle == nil {
		toggle = NewColorToggle()
	}
	parts, currentText := parseLoop(input, toggle.EnableColor)

	//flush any remainig text
	parts = flushText(parts, currentText)

	return CompiledTemplate{
		Parts:       parts,
		TotalLength: len(input),
	}
}
