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

	var (
		res []dataprovider.Record
		i   int
	)
	for {
		i++

		record, err2 := r.Read()
		if err2 != nil {
			if err2 == io.EOF {
				break
			}
			return nil, errors.Wrapf(err2, "failed to read record %d", i)
		}

		id := record[0]
		if id == "" {
			// if issue id empty, skip record parsing
			continue
		}

		startTime, err2 := time.Parse(time.TimeOnly, appendEmptySeconds(record[2]))
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to parse startTime %s at %d record", record[2], i)
		}

		endTime, err2 := time.Parse(time.TimeOnly, appendEmptySeconds(record[3]))
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to parse endTime %s at %d record", record[3], i)
		}

		tracked, err2 := strconv.ParseBool(record[4])
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to parse tracked %s at %d record", record[4], i)
		}

		res = append(res, dataprovider.Record{
			IssueID: id,
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
