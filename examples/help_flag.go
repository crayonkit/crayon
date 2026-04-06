package main

import (
    "fmt"
    "github.com/ph4mished/crayon"
)
  var(
  header = crayon.Parse("[bold fg=cyan][0][reset]")
  command  = crayon.Parse("[fg=yellow][0:<25][fg=green][1][reset]")
  flag = crayon.Parse("[fg=yellow][0][reset], [fg=yellow][1:<20] [fg=green][2][reset]")
  )

func ShowHelp() {
    header.Println("MyApp Help")
    fmt.Println()
    
    header.Println("Usage:")
    fmt.Println("  myapp [command] [options]")
    fmt.Println()
    
    header.Println("Commands:")
    command.Println("start", "Start the application")
    command.Println("stop", "Stop the application")
    command.Println("status", "Check application status")
    command.Println()
    
    header.Println("Options:")
    flag.Println("-h", "--help", "Show this help")
    flag.Println("-v", "--version", "Show version")
    flag.Println("-d", "--debug", "Enable debug mode")
}

func main(){
	ShowHelp()
}
