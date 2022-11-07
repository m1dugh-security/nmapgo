package main

import (
    // "os/exec"
    // "errors"
    "fmt"
    "log"
    "encoding/json"
)


func main() {
    file := "/tmp/scan.xml"
    /*path, err := exec.LookPath("nmap")
    if errors.Is(err, exec.ErrDot) {
        log.Fatal(err)
        err = nil
    }

    if err != nil {
        log.Fatal(err)
    }


    cmd := exec.Command(path, "-Pn", "-A", "-oX", file, "scanme.nmap.org")

    if err = cmd.Run(); err != nil {
        log.Fatal(err)
    }*/

    host, err := ExtractInfo(file)
    if err != nil {
        log.Fatal(err)
    }
    bytes, _ := json.Marshal(host)
    fmt.Println(string(bytes))

}
