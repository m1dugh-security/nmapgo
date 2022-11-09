package nmapgo

import (
    "io/ioutil"
    "os/exec"
    "errors"
    "fmt"
    "os"
    "runtime"
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

type Scanner struct {
    Options *Options
    binPath string
    outputDir string
}

func createTempDir(path string) error {

    info, err := os.Stat(path)
    if errors.Is(err, os.ErrNotExist) {
        err = os.Mkdir(path, os.ModePerm)
        if err != nil {
            return errors.New(fmt.Sprintf("could not create directory %s", path))
        }
    } else if err != nil {
        return err
    } else if !info.IsDir() {
        if err = os.Remove(path);err != nil {
            return err
        }
        if err = os.Mkdir(path, os.ModePerm);err != nil {
            return errors.New(fmt.Sprintf("could not create directory %s", path))
        }
    }
    return nil
}

func getBinPath() (string, error) {
    path, err := exec.LookPath("nmap")
    return path, err
}

func NewScanner(options *Options) (*Scanner, error) {
    dir, err := GetTempDir()
    if err != nil {
        return nil, err
    }
    if err = createTempDir(dir); err != nil {
        return nil, err
    }
    path, err := getBinPath()
    if err != nil {
        return nil, err
    }
    return &Scanner{options, path, dir}, nil
}

func (s *Scanner) getRandomFile() string {
    return fmt.Sprintf("%s/%s", s.outputDir, GenerateRandomName())
}

func (s *Scanner) ScanHost(addr string) (*Host, error) {

    file := s.getRandomFile();
    cmd := exec.Command(s.binPath,
        s.Options.ToString(),
        "-oX",
        file,
        addr,
    )
    stderr, err := cmd.StderrPipe()
    if err != nil {
        return nil, err
    }

    fmt.Println(cmd)
    if err != nil {
        return nil, err
    }
    if err := cmd.Start(); err != nil {
        return nil, err
    }

    data, err := ioutil.ReadAll(stderr)
    if err != nil {
        return nil, err
    }

    if err = cmd.Wait(); err != nil {
        return nil, errors.New(fmt.Sprintf("An error occured while running command: \n%s", string(data)))
    }

    hosts, err := ExtractInfo(file)
    if err != nil {
        return nil, err
    }

    err = os.Remove(file)
    if err != nil {
        return nil, err
    }
    if len(hosts) > 0 {
        res := &Host{}
        *res = hosts[0]
        return res, nil
    } else {
        return nil, nil
    }
}
