package main

import (
    "fmt"
    "strings"
    "github.com/ph4mished/crayon"
)

func main() {
    
    // Table with colored headers
    headerTemplate := crayon.Parse("[bold fg=cyan][0][reset]")
    rowTemplate := crayon.Parse("[0]  [fg=yellow][1][reset]  [fg=green][2][reset]")
    
    headerTemplate.Println(strings.Repeat("─", 40))
    headerTemplate.Println("USER MANAGEMENT")
    headerTemplate.Println(strings.Repeat("─", 40))
    
    rowTemplate.Println("Alice", "admin", "active")
    rowTemplate.Println("Bob", "user", "active")
    rowTemplate.Println("Charlie", "guest", "inactive")
    
    // Nested templates
    errorTemplate := crayon.Parse("[bold fg=red][0][reset]: [1]")
    suggestionTemplate := crayon.Parse("[fg=yellow]Suggestion: [0][reset]")
    
    errors := []struct{
        code string
        msg string
        suggestion string
    }{
        {"E001", "File not found", "Check the file path"},
        {"E002", "Permission denied", "Run with sudo or check permissions"},
        {"E003", "Out of memory", "Close other applications"},
    }
    
    for _, err := range errors {
        errorTemplate.Println(err.code, err.msg)
        fmt.Println("  " + suggestionTemplate.Sprint(err.suggestion))
        fmt.Println()
    }
}
