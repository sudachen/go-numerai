package numerai

import (
	"github.com/sudachen/go-fp/fu"
	"github.com/sudachen/go-ml/tables"
	"strings"
)

type v2LeaderBoard struct {
	Username                 string
	Tier                     string
	Reputation               float64
	Rank                     int
	StakedRank               int
	MmcRep                   float64
	MmcRank                  int
	MmcStakedRank            int
	PrevRank                 int
	PrevStakedRank           int
	PrevMmcRank              int
	PrevMmcStakedRank        int
	CorrWithMmAvg            float64
	RollingScoreRep          float64
	NmrStaked                float64
	BonusPerc                int
	Bio                      string
	AverageCorrelationPayout float64
	LeaderboardBonus         float64
}

var v2LeaderBoardQuery = `query { v2Leaderboard { ` + strings.Join(fu.FieldsOf((*v2LeaderBoard)(nil)), " ") + ` } }`

func Leaderboard() *tables.Table {
	r, err := RawQuery(v2LeaderBoardQuery, QueryArgs{})
	if err != nil {
		panic(err.Error())
	}
	c := make(chan v2LeaderBoard)
	go func() {
		defer close(c)
		for y := range r.Q("data").Q("v2Leaderboard").Chan() {
			s := v2LeaderBoard{}
			y.Fill(&s)
			c <- s
		}
	}()
	return tables.New(c)
}
