# DRLM2T (DRLM v2 Testing)

## infrastructure.yaml

```go
type DRLMTestingConfig struct {
 Name      string `mpastructure:"name"`
 Prefix    string `mapstructure:"prefix"`
 Templates string `mpastructure:"templates"`
 URL       string `mpastructure:"url"`
 DefIP     string `mpastructure:"defip"`
 DefMask   string `mpastructure:"defmask"`
 DefDNS    string `mpastructure:"defdns"`
 DefTem    string `mpastructure:"deftmp"`
 Kvms  []Kvm     `mpastructure:"kvms"`
 Nets  []Network `mapstructure:"nets"`
 Hosts []Host    `mapstructure:"hosts"`
}
```

```go
type Kvm struct {
 HostName  string `mapstructure:"hostname"`
 User      string `mapstructure:"user"`
 URI       string `mapstructure:"uri"`
 Prefix    string `mapstructure:"prefix"`
 Templates string `mapstructure:"templates"`
 DefIP     string `mpastructure:"defip"`
 DefMask   string `mpastructure:"defmask"`
 DefDNS    string `mpastructure:"defdns"`
 DefTmp    string `mpastructure:"deftmp"`
}
```

```go
type Network struct {
 Name        string `mapstructure:"name"`
 Kvm         string `mapstructure:"kvm"`
 Mac         string `mapstructure:"mac"`
 IP          string `mapstructure:"ip"`
 Mask        string `mapstructure:"mask"`
 Gateway     string `mapstructure:"gateway"`
 DNS         string `mapstructure:"dns"`
 DhcpStartIP string `mapstructure:"dhcpstartip"`
 DhcpEndIP   string `mapstructure:"dhcpendip"`
 Prefix      string `mapstructure:"prefix"`
}
```

```go
type Host struct {
 Name     string    `mapstructure:"name"`
 Kvm      string    `mapstructure:"kvm"`
 Template string    `mapstructure:"template"`
 Prefix   string    `mapstructure:"prefix"`
 Nets     []Network `mapstructure:"nets"`
 Tests    []Test    `mapstructure:"tests"`
}
```

```go
type Test struct {
 Index        int       `mapstructure:"index"`
 Status       int       `mapstructure:"status"`
 Name         string    `mapstructure:"name"`
 TestType     Tipus     `mapstructure:"testtype"`
 Mode         Mode      `mapstructure:"mode"`
 CommandToRun string    `mapstructure:"commandtorun"`
 FileToRun    string    `mapstructure:"filetorun"`
 Args         []string  `mapstructure:"args"`
 Expect       string    `mapstructure:"expect"`
 Output       string    `mapstructure:"output"`
 LandMark     bool      `mapstructure:"landmark"`
 Dependencies []Deps    `mapstructure:"dependencies"`
}

type Deps struct {
 Host string `mapstructure:"host"`
 Test int    `mapstructure:"test"`
}
```
