package main

import (
    "asciimatcher/match"
    "fmt"
    "log"
    "os"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: file_what_to_find file_where_to_find")
        os.Exit(0)
    }
    patternMap, err := match.ParseIncomingPattern(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    found, err := match.FindByPattern(patternMap, os.Args[2])
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Matched: %d\n", found)
}
