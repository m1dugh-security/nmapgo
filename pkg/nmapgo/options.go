package nmapgo

import (
    "fmt"
    "strings"
)

var (
    TcpSynScan string = "S"
    ConnectScan = "T"
    ACKScan = "A"
    WindowScan = "W"
    UDPScan = "U"
)

var (
    DefaultScan string = ""
    DefaultAggressive bool = false
    DefaultPing bool = false
)

const (
    aggressiveFlag string = "-A"
    pingFlag = "-Pn"
    scanFlag = "-s"
)

type Options struct {
    Scan string
    Ping bool
    Aggressive bool
}

func NewOptions() *Options {
    return &Options{
        Scan: DefaultScan,
        Ping: DefaultPing,
        Aggressive: DefaultAggressive,
    }
}

func (o *Options) ToString() string {
    var res []string
    if len(o.Scan) > 0 {
        res = append(res, fmt.Sprintf("%s%s", scanFlag, o.Scan))
    }

    if o.Ping {
        res = append(res, pingFlag)
    }
    if o.Aggressive {
        res = append(res, aggressiveFlag)
    }

    return strings.Join(res, " ")
}
