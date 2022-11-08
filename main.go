package main

import (
    "log"
    "encoding/json"
    "fmt"
)


func main() {
    options := NewOptions()
    options.Aggressive = true
    scanner, err := NewScanner(options)
    if err != nil {
        log.Fatal(err)
    }

    host, err := scanner.ScanHost("scanme.nmap.org")
    if err != nil {
        log.Fatal(err)
    }
    b, _ := json.Marshal(host)
    fmt.Println(string(b))
}
