package numerai

import (
	"fmt"
	"github.com/sudachen/go-foo/fu"
	"github.com/sudachen/go-foo/lazy"
	"github.com/sudachen/go-ml/graphql"
	"github.com/sudachen/go-ml/tables"
	"github.com/sudachen/go-ml/tables/csv"
)

const (
	TrainSubset = iota
	ValidationSubset
	TestSubset
	LiveSubset
)

const datasetTrainingCSV = "numerai_training_data.csv"
const datasetTournamentCSV = "numerai_tournament_data.csv"
const kazutsugiTarget = "target_kazutsugi"
const kazutsugiPrediction = "prediction_kazutsugi"
const supportedTournament = 8 /*kazutsugi*/
const supportedTournamentName = "kazutsugi"

const datasetQuery = `
query($tournament:Int!){
    rounds(status:OPEN,tournament:$tournament) {
      number
    }
  	dataset(tournament:$tournament)
  }
`
func dataset(file string) (stream lazy.Stream) {
	return lazy.Wrap(
		graphql.IfQuery(numeraiUrl, datasetQuery, graphql.Args{"tournament": supportedTournament},
		func(r graphql.Result) interface{}{
			subset := tables.Enumset{
				"train": TrainSubset,
				"validation": ValidationSubset,
				"test": TestSubset,
				"live": LiveSubset}
			era := tables.Enumset{}
			round := r.Q("data").Q("rounds").Q(0).Int("number")
			url := r.Q("data").String("dataset")
			cachefile := fmt.Sprintf("datasets/numerai/%v/numerai_datasets_%v.zip", supportedTournamentName, round)
			return csv.Source(fu.ZipFile(file, fu.External(url, fu.Cached(cachefile))),
				csv.String("id").As("Id"),
				csv.Meta(subset.Integer(), "data_type").As("Subset"),
				csv.Meta(era.Integer(), "era").As("Era"),
				csv.Float32("feature_intelligence*").As("Feature*i"),
				csv.Float32("feature_charisma*").As("Feature*h"),
				csv.Float32("feature_strength*").As("Feature*s"),
				csv.Float32("feature_dexterity*").As("Feature*d"),
				csv.Float32("feature_constitution*").As("Feature*c"),
				csv.Float32("feature_wisdom*").As("Feature*w"),
				csv.Float32(kazutsugiTarget).As("Label"))()
		}))
}

var Tournament tables.Lazy = func() lazy.Stream { return dataset(datasetTournamentCSV) }
var Training tables.Lazy = func() lazy.Stream { return dataset(datasetTrainingCSV) }
