package numerai

import "github.com/sudachen/go-ml/util"

const (
	VersionMajor = 1
	VersionMinor = 0
	VersionPatch = 0
)

const Version util.VersionType = VersionMajor*10000 + VersionMinor*100 + VersionPatch

type Auth struct {
	PublicID  string
	SecretKey string
}
