package numerai

import (
	"time"
)

const currentRoundInfoQuery = `
query {
  rounds(status:OPEN) {
	number
	status
	openTime
	closeTime
  }
}`

const roundInfoQuery = `
query $number: Int!{
  rounds(number: $number) {
	number
	status
	openTime
	closeTime
  }
}`

type RoundStatus int

const (
	CLOSE RoundStatus = iota
	OPEN
)

type RoundInfo struct {
	Number    int
	Status    RoundStatus
	OpenTime  time.Time
	CloseTime time.Time
}

func GetRoundInfo(round int) RoundInfo {
	if round <= 0 {
		return CurrentRound()
	}
	r, err := RawQuery(roundInfoQuery, QueryArgs{"round": round})
	if err != nil {
		panic(err.Error())
	}
	return asRoundInfo(r)
}

func CurrentRound() RoundInfo {
	r, err := RawQuery(currentRoundInfoQuery, QueryArgs{})
	if err != nil {
		panic(err.Error())
	}
	return asRoundInfo(r)
}

func asRoundInfo(r QueryResult) RoundInfo {
	v := r.Q("data").Q("rounds").Q(0)
	return RoundInfo{
		Number:    v.Int("number"),
		Status:    asRoundStatus(v.String("status")),
		OpenTime:  v.Time("openTime"),
		CloseTime: v.Time("closeTime"),
	}
}

func asRoundStatus(s string) RoundStatus {
	if s == "OPEN" {
		return OPEN
	}
	return CLOSE
}
