package main

import (
    "github.com/ph4mished/crayon"
)

func main(){
    row := crayon.Parse("[fg=cyan bold][0:<20][fg=yellow][1:<10][reset]")
    
    row.Println("Alice", "admin")
    row.Println("Bob", "user")
    row.Println("Charlie", "guest")
}
