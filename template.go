// Package crayon provides terminal colors and styles for Go.
// A better go doc is needed to be written
package crayon

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"io"
	"golang.org/x/term"
)

type TempPart struct {
  Text string
  Index int
  FormatStr string
}

type CompiledTemplate struct {
  Parts []TempPart
  TotalLength int
}

type ColorToggle struct {
  EnableColor bool
}


//=============================
// COLOR TOGGLE
//=============================

func autoDetect() bool {
  if _, exists := os.LookupEnv("NO_COLOR"); exists{
    return false	
  }
  return term.IsTerminal(int(os.Stdout.Fd()))
}



//=============================
// PARSE - LOOP
//=============================

func parseLoop(input string, enableColor bool) ([]TempPart, string){
	var (
		parts           []TempPart
		currentText     string
		contentSequence string
		inReadSequence  bool
	)

	for i, ch := range(input){
		char := string(ch)

		switch {
		case char == "[" && !inReadSequence:
			parts, currentText, contentSequence, inReadSequence = handleOpenBracket(i, input, parts, currentText)
			
		case ch == ']' && inReadSequence:
			parts, inReadSequence = handleCloseBracket(contentSequence, parts, enableColor)
			contentSequence = ""

		case inReadSequence:
			contentSequence += char

		default:
			currentText += char
		}
	}
	return parts, currentText
}



//=============================
// PARSE - BRACKET HANDLERS
//=============================

func handleOpenBracket(i int, input string, parts []TempPart, currentText string) ([]TempPart, string, string, bool) {
	//check if the next value is "["
    // [[fg=color]] should never be an escape
    //consider first '[' as a text, move until, content is found. 
	if i+1 < len(input) && input[i+1] == '['{
	  currentText += "["
	  return parts, currentText, "", false
	}
	//flush current text before entering sequence
	parts = flushText(parts, currentText)
	return parts, "", "", true
}

func handleCloseBracket(contentSequence string, parts []TempPart, enableColor bool) ([]TempPart, bool){
	allWords := strings.Fields(contentSequence)

	if isColorSequence(allWords) {
		parts = handleColorSequence(parts, allWords, enableColor)
	} else {
		parts = handleNonColorSequence(parts, contentSequence)
	}
	return parts, false
}

//=============================
// PARSE - SEQUENCE HANDLERS
//=============================

func isColorSequence(words []string) bool {
	if len(words) == 0 {
		return false
	}
	for _, w := range words {
		if !IsSupportedColor(w){
			return false
		}
	}
	return true
}

func handleColorSequence(parts []TempPart, words []string, enableColor bool) []TempPart {
	if enableColor{
		for _, w := range words {
			parts = append(parts, TempPart{Text: ParseColor(w), Index: -1, FormatStr: ""})
		}
	} else {
		parts = append(parts, TempPart{Text: "", Index: -1, FormatStr: ""})
	}
	return parts
}

func handleNonColorSequence(parts []TempPart, contentSequence string) []TempPart {
	if isValidPlaceholder(contentSequence){
		return handlePlaceholder(parts, contentSequence)
	}

	//for padded placeholders
	if strings.Contains(contentSequence, ":") {
			return handlePaddedPlaceholder(parts, contentSequence)
		}
	//unrecognized -  pass through as literal
	return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
}

//=============================
// PARSE - PLACEHOLDER
//=============================
//Monolithic Parse() should be divided into subsections of funcs so that they can be reused for escapes and other things
//extract placeholders
//placeholders will support padding too. [0:<20] = right alignment, [0:>20] = left align
//Overflow handling will slow down crayon. I'm still on the fence of throwing it away or using it
//It will slow down crayon because calculation will be moved to apply, thats not the work of apply
//[0:<20!] = right alignment  with truncation, [0:>20~] = left align with elipsis(...),
//[0:<20?] = right alignment  with warn to stderr

func isValidPlaceholder(input string) bool{
  return len(input) > 0 && allDigits(input)
}

func handlePlaceholder(parts []TempPart, contentSequence string) []TempPart {
    //decided to make it flexible and accept more indices but its still prone to overflow
      //needs a digit boundary guard	
	  index, err := strconv.Atoi(contentSequence)
	  if err == nil && index >= 0 && index <= 999 {
      return append(parts, TempPart{Text: "", Index: index, FormatStr: ""})
	}
	//out of range -treat as literal
	return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
}

func handlePaddedPlaceholder(parts []TempPart, contentSequence string) []TempPart {
	//[0:>20] stripped of its brackets ==> 0:>20
	splitWord := strings.SplitN(contentSequence, ":", 2) // ==> [0 >20]
	if len(splitWord) != 2 {
		return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
	}
	indexStr := splitWord[0]
	padStr := splitWord[1]
	//parse indexStr
	index, err := strconv.Atoi(indexStr)
	  if err != nil || index < 0 || index > 999 {
      return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
	} 

	//parse the padStr
	align, width, err := parseAlignWidth(padStr)
	if err != nil {
		fmt.Println("Error: ", err)
		//return append(parts, TempPart{Text: "[" + contentSequence + "]", Index: -1, FormatStr: ""})
		return nil
	}
	return append(parts, TempPart{Text: "", Index: index, FormatStr: buildFormatStr(align, width)})
}


//=============================
// PARSE - HELPERS
//=============================
func parseAlignWidth(input string) (rune, int, error){
	if len(input) < 2 {
		return 0, 0, fmt.Errorf("invalid padding specification '%s'", input)
	}
	align := rune(input[0])
	widthStr := input[1:]
	
	width, err := strconv.Atoi(widthStr)
	if err != nil || width <= 0 {
		return 0, 0, fmt.Errorf("invalid width: negative width not supported '%d'", width)
	}

	switch align {
	case '<', '>':
		return align, width, nil
	default:
		return 0, 0, fmt.Errorf("invalid alignment char '%s'", string(align))
	}
}

func buildFormatStr(align rune, width int) string {
	switch align{
	case '<':
		return fmt.Sprintf("%%-%ds", width)
	case '>':
		return fmt.Sprintf("%%%ds", width)
	}
	return ""
}

func flushText(parts []TempPart, currentText string) []TempPart {
	if len(currentText) > 0 {
		parts = append(parts, TempPart{Text: currentText, Index: -1})
	}
	return parts
}

func allDigits(s string) bool {
  for _, r := range s{
	if !unicode.IsDigit(r){
	  return false
	}
  }
  return true
}

  
//=============================
// APPLY
//=============================

//apply will be a private func
func (temp CompiledTemplate) apply(args ...any) string {
  //Calculate estimated size for optimization
  var totalArgLength int
  for _, arg := range args{
	totalArgLength += len(fmt.Sprint(arg))
  }

  estimatedSize := temp.TotalLength + totalArgLength
  var result strings.Builder
  result.Grow(estimatedSize)

  for _, part := range temp.Parts{
	if part.Index < 0{
	  result.WriteString(part.Text)
	} else if part.Index < len(args) {
		value := fmt.Sprint(args[part.Index])
		if part.FormatStr != ""{
			value = fmt.Sprintf(part.FormatStr, value)
		}
		result.WriteString(value)
	}
  }
  return result.String()
}



//=======================
// PRINT
//=======================
func (temp CompiledTemplate) Println(args ...any) {
	fmt.Println(temp.apply(args...))
}

func (temp CompiledTemplate) Print(args ...any) {
	fmt.Print(temp.apply(args...))
}




//=======================
// EPRINT
//=======================
func (temp CompiledTemplate) Eprintln(args ...any) {
	fmt.Fprintln(os.Stderr, temp.apply(args...))
}

func (temp CompiledTemplate) Eprint(args ...any) {
	fmt.Fprint(os.Stderr, temp.apply(args...))
}



//=======================
// FPRINT
//=======================
func (temp CompiledTemplate) Fprintln(w io.Writer, args ...any) (n int, err error) {
	return fmt.Fprintln(w, temp.apply(args...))
}

func (temp CompiledTemplate) Fprint(w io.Writer, args ...any) (n int, err error){
	return fmt.Fprint(w, temp.apply(args...))
}


//=======================
// SPRINT
//=======================
func (temp CompiledTemplate) Sprintln(args ...any) string {
	return fmt.Sprintln(temp.apply(args...))
}

func (temp CompiledTemplate) Sprint(args ...any) string {
	return temp.apply(args...)
}
