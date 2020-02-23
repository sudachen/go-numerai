package numerai

import (
	"archive/zip"
	"github.com/sudachen/go-ml/tables"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GetTurnamentData(dirname string) *tables.Table {
	if !HasDatasetsIn(dirname) {
		if err := DownloadCurrentDatasets(dirname); err != nil {
			panic(err.Error())
		}
	}
	return nil
}

func GetTrainingData(dirname string) *tables.Table {
	if !HasDatasetsIn(dirname) {
		if err := DownloadCurrentDatasets(dirname); err != nil {
			panic(err.Error())
		}
	}
	return nil
}

const datasetZipFile = "numerai_datasets.zip"
const datasetTrainingCSV = "numerai_training_data.csv"
const datasetTurnamentCSV = "numerai_turnament_data.csv"

func HasDatasetsIn(dirname string) bool {
	_, err := os.Stat(filepath.Join(dirname, datasetTrainingCSV))
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(dirname, datasetTurnamentCSV))
	return err == nil
}

func downloadCurrentDatasetsZip(dirname string) (err error) {
	r, err := RawQuery(`query {dataset}`, QueryArgs{})
	if err != nil {
		return
	}
	url := r.Q("data").String("dataset")
	b, err := http.Get(url)
	if err != nil {
		return
	}
	defer b.Body.Close()
	f, err := os.Create(filepath.Join(dirname, datasetZipFile))
	if err != nil {
		return
	}
	defer f.Close()
	_, err = io.Copy(f, b.Body)
	return err
}

func unzipCurrentDatasetsZip(dirname string) (err error) {
	r, err := zip.OpenReader(filepath.Join(dirname, datasetZipFile))
	if err != nil {
		return
	}
	for _, f := range r.File {
		var x io.ReadCloser
		if x, err = f.Open(); err != nil {
			return
		}
		var o io.WriteCloser
		if o, err = os.Create(filepath.Join(dirname, datasetZipFile)); err != nil {
			return
		}
		defer o.Close()
		if _, err = io.Copy(o, x); err != nil {
			return
		}
	}
	return
}

func DownloadCurrentDatasets(dirname string) (err error) {
	if err = os.MkdirAll(dirname, 0655); err != nil {
		return err
	}
	if err = downloadCurrentDatasetsZip(dirname); err != nil {
		return
	}
	if err = unzipCurrentDatasetsZip(dirname); err != nil {
		_ = os.RemoveAll(dirname)
	}
	return
}
