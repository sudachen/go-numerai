package nuai

import (
	"fmt"
	"github.com/sudachen/go-ml/tables"
	"reflect"
	"time"
)

const laderboardQuery = `query($number: Int!) {
  rounds(number: $number) {
	leaderboard {
	  liveLogloss
	  submissionId
	  username
	  validationLogloss
	  paymentGeneral {
		nmrAmount
		usdAmount
	  }
	  paymentStaking {
		nmrAmount
		usdAmount
	  }
	  totalPayments {
		nmrAmount
		usdAmount
	  }
	  stake {
		insertedAt
		soc
		confidence
		value
		txHash
	  }
	}
  }
}`

type Round int

func (nuai *Nuai) GetLeaderboard(round int) *tables.Table {
	r, err := nuai.RawQuery(laderboardQuery, QueryArgs{"number":round})
	if err != nil {
		panic(err.Error())
	}
	stakes := r.Q("data").Q("rounds").Q(0).Q("leaderboard")
	fmt.Println(reflect.Value(stakes))
	type S struct {
		Username string
		TotalNMR float32
		TotalUSD float32
		StakingNMR float32
		StakingUSD float32
		Conidence float32
		Soc float32
		Value float32
		TxHash string
		InsertedAt time.Time
		SubmissionId string
	}
	c := make(chan S)
	go func(){
		for _,y := range stakes.List() {
			c <- S{
			y.String("username"),
			y.Q("totalPayments").Float("nmrAmount"),
			y.Q("totalPayments").Float("usdAmount"),
			y.Q("paymentStaking").Float("nmrAmount"),
			y.Q("paymentStaking").Float("usdAmount"),
			y.Q("stake").Float("confidence"),
			y.Q("stake").Float("soc"),
			y.Q("stake").Float("value"),
			y.Q("stake").String("txHash"),
			y.Q("stake").Time("insertedAt"),
			y.String("submissionId"),
			}
		}
		close(c)
	}()
	return tables.New(c)
}
