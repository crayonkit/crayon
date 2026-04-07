package crayon

import (
	"fmt"
	"testing"
)

func TestExample(t *testing.T) {
	tempAnsi := Parse("FOR ANSI: [bold fg=blue]Hello[reset]")
	temp256 := Parse("FOR 256: [bold fg=115]Hello[reset]")
	tempHex := Parse("FOR HEX: [bold fg=#AABBCC]Hello[reset]")
	tempRGB := Parse("FOR RGB: [bold fg=rgb(15,102,224)]Hello [reset]")
	fmt.Println("STATIC TEMPLATES")
	tempAnsi.Println()
	temp256.Println()
	tempHex.Println()
	tempRGB.Println()
}

func TestWithToggle(t *testing.T) {
	toggle := NewColorToggle(false)

	tempAnsi := toggle.Parse("FOR ANSI: [bold fg=blue]Hello[reset]")
	temp256 := toggle.Parse("FOR 256: [bold fg=115]Hello[reset]")
	tempHex := toggle.Parse("FOR HEX: [bold fg=#AABBCC]Hello[reset]")
	tempRGB := toggle.Parse("FOR RGB: [bold fg=rgb(15,102,224)]Hello [reset]")
	fmt.Println("\n\nTEMPLATES WITH TOGGLE(COLOR OFF)")
	tempAnsi.Println()
	temp256.Println()
	tempHex.Println()
	tempRGB.Println()

}

func TestForInterpolation(t *testing.T) {
	tempAnsi := Parse("FOR ANSI: [bold fg=blue]Hello [fg=yellow][0][reset]")
	temp256 := Parse("FOR 256: [bold fg=115]Hello [fg=13][0][reset]")
	tempHex := Parse("FOR HEX: [bold fg=#AABBCC]Hello [fg=#AAFFCC][0][reset]")
	tempRGB := Parse("FOR RGB: [bold fg=rgb(15,102,224)]Hello [fg=rgb(10,94,104)][0][reset]")
	fmt.Println("\n\nTEMPLATES WITH PLACEHOLDER")
	tempAnsi.Println("World")
	temp256.Println("World")
	tempHex.Println("World")
	tempRGB.Println("World")

}

//test for internal (unexportable functions) will also be made
