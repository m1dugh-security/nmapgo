package main

import (
    "os/exec"
    "errors"
    "fmt"
    "os"
    "runtime"
    "log"
    "encoding/json"
    "math/rand"
)

const SUFFIX = "nmapgo"

func GetTempDir() (string, error) {
    var res string
    if runtime.GOOS == "windows" {
        res = `C:\WINDOWS\Temp\`
    } else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
        res = `/tmp/`
    } else {
        return "", errors.New("could not temp directory")
    }

    info, err := os.Stat(res)
    if !info.IsDir() {
        return "", errors.New("temp path is not a directory")
    }
    if errors.Is(err, os.ErrNotExist) {
        return "", errors.New("could not find temp dir")
    }

    return res + SUFFIX, nil
}

func GenerateRandomString(len int) string {
    res := ""
    for i := 0; i < len; i++ {
        value := rand.Intn(62)
        start := 'a'
        if value >= 52 {
            start = '0'
            value -= 52
        } else if value >= 26 {
            start = 'A'
            value -= 26
        }
        b := byte(start) + byte(value)
        res += string(b)
    }

    return res
}

func GenerateRandomName() string {
    res := GenerateRandomString(16)
    for i := 0; i < 3;i++ {
        res += "_" + GenerateRandomString(16)
    }

    return res
}

func main() {
    dir, err := GetTempDir()
    if err != nil {
        log.Fatal(err)
    }
    info, err := os.Stat(dir)
    if errors.Is(err, os.ErrNotExist) {
        err = os.Mkdir(dir, os.ModePerm)
        if err != nil {
            log.Fatal(err)
        }
    } else if err != nil {
        log.Fatal(err)
    } else if !info.IsDir() {
        log.Fatal(fmt.Sprintf("file %s exists and creates collision with nmapgo\n", dir))
    }

    file := fmt.Sprintf("%s/%s", dir, GenerateRandomName())
    path, err := exec.LookPath("nmap")
    if errors.Is(err, exec.ErrDot) {
        log.Fatal(err)
        err = nil
    }

    if err != nil {
        log.Fatal(err)
    }


    cmd := exec.Command(path,
        "-Pn",
        "-A",
        "-oX",
        file,
        "scanme.nmap.org",
    )
    if err = cmd.Run(); err != nil {
        log.Fatal(err)
    }

    host, err := ExtractInfo(file)
    if err != nil {
        log.Fatal(err)
    }
    bytes, _ := json.Marshal(host)
    fmt.Println(string(bytes))

    err = os.Remove(file)
    if err != nil {
        log.Fatal(fmt.Sprintf("could not remove temp file %s\n", file))
    }
}
