package browingdata

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type outPutter struct {
	json bool
	csv  bool
}

func newOutPutter(flag string) *outPutter {
	o := &outPutter{}
	if flag == "json" {
		o.json = true
	} else {
		o.csv = true
	}
	return o
}

func (o *outPutter) Write(data Source, writer io.Writer) error {
	switch o.json {
	case true:
		encoder := json.NewEncoder(writer)
		encoder.SetIndent("  ", "  ")
		encoder.SetEscapeHTML(false)
		return encoder.Encode(data)
	default:
		gocsv.SetCSVWriter(func(w io.Writer) *gocsv.SafeCSVWriter {
			writer := csv.NewWriter(transform.NewWriter(w, unicode.UTF8BOM.NewEncoder()))
			writer.Comma = ','
			return gocsv.NewSafeCSVWriter(writer)
		})
		return gocsv.Marshal(data, writer)
	}
}

func (o *outPutter) CreateFile(dir, filename string) (*os.File, error) {
	if filename == "" {
		return nil, errors.New("empty filename")
	}

	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0o750)
			if err != nil {
				return nil, err
			}
		}
	}

	var file *os.File
	var err error
	p := filepath.Join(dir, filename)
	file, err = os.OpenFile(filepath.Clean(p), os.O_TRUNC|os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (o *outPutter) Ext() string {
	if o.json {
		return "json"
	}
	return "csv"
}

func getJson(data Source) (string, error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("  ", "  ")
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Data encode to json error: %s", err))
	}
	return string(buf.Bytes()), nil
}

func getCSV(data Source) (string, error) {
	gocsv.SetCSVWriter(func(w io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(transform.NewWriter(w, unicode.UTF8BOM.NewEncoder()))
		writer.Comma = ','
		return gocsv.NewSafeCSVWriter(writer)
	})
	buf := new(bytes.Buffer)
	err := gocsv.Marshal(data, buf)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
