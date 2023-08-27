package dataprovider

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"gotrack/pkg/app/dataprovider"
)

func NewCsvFileDataProvider(path string) dataprovider.Provider {
	return &csvFileDataProvider{path: path}
}

type csvFileDataProvider struct {
	path string
}

func (provider *csvFileDataProvider) Read(ctx context.Context) ([]dataprovider.Record, error) {
	file, err := os.Open(provider.path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	r := csv.NewReader(file)
	_, err = r.Read() // skip first header line
	if err != nil {
		return nil, err
	}

	var res []dataprovider.Record
	for {
		record, err2 := r.Read()
		if err2 != nil {
			if err2 == io.EOF {
				break
			}
			return nil, errors.Wrap(err2, "failed to read record")
		}

		startTime, err2 := time.Parse(time.TimeOnly, appendEmptySeconds(record[2]))
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to parse startTime %s", record[2])
		}

		endTime, err2 := time.Parse(time.TimeOnly, appendEmptySeconds(record[3]))
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to parse endTime %s", record[3])
		}

		tracked, err2 := strconv.ParseBool(record[4])
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to parse tracked %s", record[4])
		}

		res = append(res, dataprovider.Record{
			IssueID: record[0],
			Comment: record[1],
			Start:   startTime,
			End:     endTime,
			Tracked: tracked,
		})
	}

	return res, nil
}

func appendEmptySeconds(t string) string {
	return fmt.Sprintf("%s:00", t)
}
