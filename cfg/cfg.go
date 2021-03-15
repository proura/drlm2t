package cfg

var Config *DRLM2TConfiguration

type DRLM2TConfiguration struct {
	Drlm2tPath string `mapstructure:"drlm2tPath"`
}
