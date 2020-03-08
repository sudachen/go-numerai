package numerai

import "github.com/sudachen/go-ml/mlutil"

const (
	VersionMajor = 1
	VersionMinor = 0
	VersionPatch = 0
)

const Version mlutil.VersionType = VersionMajor*10000 + VersionMinor*100 + VersionPatch
const numeraiUrl = "https://api-tournament.numer.ai"
