package version

import (
	"github.com/blang/semver/v4"
)

const RawVersion = "0.0.1-dev"

var Version = semver.MustParse(RawVersion)
