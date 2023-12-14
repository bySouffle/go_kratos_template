package client

import (
	"github.com/guonaihong/gout"
	"time"
)

var DefaultRequestTimeout = time.Second * 30
var GoutNoTLS = gout.NewWithOpt(gout.WithInsecureSkipVerify(), gout.WithTimeout(DefaultRequestTimeout))
