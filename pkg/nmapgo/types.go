package nmapgo

import (
    "io/ioutil"
    "os"
    "encoding/xml"
    "errors"
)

type nmapservice struct {
    XMLName xml.Name `xml:"service"`
    Name string `xml:"name,attr"`
    Product string `xml:"product,attr"`
    Version string `xml:"version,attr"`
    Additionals string `xml:"extrainfo,attr"`
}

func (s nmapservice) toService() Service {
    return Service{
        s.Name,
        s.Product,
        s.Version,
        s.Additionals,
    }
}

type nmapport struct {
    XMLName xml.Name `xml:"port"`
    Protocol string `xml:"protocol,attr"`
    Port int `xml:"portid,attr"`
    State struct{
        State string `xml:"state,attr"`
    } `xml:"state"`
    Service nmapservice `xml:"service"`
}

func (p nmapport) toPort() Port {
    return Port{
        p.Protocol,
        p.Port,
        p.State.State,
        p.Service.toService(),
    }
}

type nmapports struct {
    XMLName xml.Name `xml:"ports"`
    Ports   []nmapport `xml:"port"`
}

func (p nmapports) toPorts() []Port {
    res := make([]Port, len(p.Ports))
    for i := 0; i < len(p.Ports);i++ {
        res[i] = p.Ports[i].toPort()
    }

    return res
}

type nmaphostnames struct {
    XMLName xml.Name `xml:"hostnames"`
    Hostnames []struct{
        Name string `xml:"name,attr"`
        Type string `xml:"type,attr"`
    }  `xml:"hostname"`
}

func (h nmaphostnames) toHostname() string {
    for _, v := range h.Hostnames {
        if v.Type == "user" {
            return v.Name
        }
    }

    return ""
}

type nmaphost struct {
    XMLName xml.Name `xml:"host"`
    Address struct{
        XMLName xml.Name `xml:"address"`
        Addr string `xml:"addr,attr"`
    } `xml:"address"`
    Hostnames nmaphostnames `xml:"hostnames"`
    Ports nmapports `xml:"ports"`
}

func (h nmaphost) toHost() Host {
    return Host{
        h.Address.Addr,
        h.Hostnames.toHostname(),
        h.Ports.toPorts(),
    }
}

type nmapresults struct {
    Hosts   []nmaphost  `xml:"host"`
}

func (r nmapresults) toHostList() []Host {
    res := make([]Host, len(r.Hosts))
    for i, v := range r.Hosts {
        res[i] = v.toHost()
    }

    return res
}

type Service struct {
    Name string     `json:"name"`
    Product string  `json:"product"`
    Version string  `json:"version"`
    Additionals string  `json:"additionals"`
}

type Port struct {
    Protocol string `json:"protocol"`
    Port int        `json:"portid"`
    State string    `json:"state"`
    Service Service `json:"service"`
}

type Host struct {
    Address string `json:"address"`
    Hostname    string `json:"hostname"`
    Ports   []Port `json:"ports"`
}

func ExtractInfo(path string) ([]Host, error) {
    file, err := os.Open(path)

    if err != nil {
        return nil, errors.New("could not open file")
    }

    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, errors.New("could not read file")
    }

    var res nmapresults
    err = xml.Unmarshal(bytes, &res)
    if err != nil {
        return nil, err
    }

    return res.toHostList(), nil
}

