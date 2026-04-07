//to test if true color fallback works
package main
import(
//	"fmt"
	"github.com/ph4mished/crayon"
)

func main(){
	col := crayon.Parse("[fg=#ff5fd7]HELLO [fg=#d7875f]WORLD[reset]")
	 cor := crayon.Parse("[fg=rgb(0,215,95)]HELLO[reset]")
	col.Println()
	cor.Println()
}
