package numerai

import (
	"github.com/sudachen/go-ml/graphql"
	"time"
)

const currentRoundInfoQuery = `
query {
  rounds(status:OPEN) {
	number
	openTime
	closeTime
  }
}`

const roundInfoQuery = `
query($number:Int!){
  rounds(number: $number) {
	number
	openTime
	closeTime
  }
}`

type RoundStatus int

type RoundInfo struct {
	Number    int
	OpenTime  time.Time
	CloseTime time.Time
}

func QueryRoundInfo(round int) (RoundInfo, error) {
	if round <= 0 {
		return QueryCurrentRound()
	}
	r, err := graphql.DoQuery(numeraiUrl, roundInfoQuery, graphql.Args{"number": round})
	if err != nil {
		return RoundInfo{}, nil
	}
	return asRoundInfo(r), nil
}

func QueryCurrentRound() (RoundInfo, error) {
	r, err := graphql.DoQuery(numeraiUrl, currentRoundInfoQuery, graphql.Args{})
	if err != nil {
		return RoundInfo{}, nil
	}
	return asRoundInfo(r), nil
}

func asRoundInfo(r graphql.Result) RoundInfo {
	v := r.Q("data").Q("rounds").Q(0)
	return RoundInfo{
		Number:    v.Int("number"),
		OpenTime:  v.Time("openTime"),
		CloseTime: v.Time("closeTime"),
	}
}
