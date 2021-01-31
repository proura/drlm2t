//////////////////////////////////
//!!!!Actualy not in use....!!!!!
//////////////////////////////////

package cfg

import "github.com/proura/drlm2t/model"

var Config *DRLM2TConfiguration

type DRLM2TConfiguration struct {
	Net model.Network `mapstructure:"network"`
}
