package store

import (
	"errors"
	"time"

	"github.com/modell-aachen/gologger/internal/interfaces"
	"github.com/modell-aachen/gologger/internal/mocks/dummydata"
)

func CreateStoreMock() *StoreMock {
	mock := StoreMock{}

	return &mock
}

type StoreMock struct {
	Stored   []interfaces.LogRow
	isClosed bool
}

func (instance *StoreMock) Fill() {
	instance.Store(dummydata.Row0)
	instance.Store(dummydata.Row1)
	instance.Store(dummydata.Row2)

}
func (instance *StoreMock) Store(row interfaces.LogRow) (err error) {
	if instance.isClosed {
		return errors.New("Handle has been closed")
	}
	instance.Stored = append(instance.Stored, row)
	return nil
}

func (instance *StoreMock) Read(startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString) (rows []interfaces.LogRow, err error) {
	if instance.isClosed {
		return rows, errors.New("Handle has been closed")
	}
	var filteredRows []interfaces.LogRow
	for _, row := range instance.Stored {
		inStartRange := row.Time.Equal(startTime) || row.Time.After(startTime)
		inEndRange := row.Time.Equal(endTime) || row.Time.Before(endTime)
		if inStartRange && inEndRange {
			if source == interfaces.SourceString("") || row.Source == source {
				levelFits := false
				for _, level := range levels {
					if level == row.Level {
						levelFits = true
						break
					}
				}
				if levelFits {
					filteredRows = append(filteredRows, row)
				}
			}
		}
	}
	return filteredRows, nil
}

func (instance *StoreMock) Close() {
	instance.isClosed = true
}

func (instance *StoreMock) IsClosed() bool {
	return instance.isClosed
}

func (instance *StoreMock) CleanUp() (err error) {
	if instance.isClosed {
		return errors.New("Handle has been closed")
	}

	instance.Stored = make([]interfaces.LogRow, 0)
	return nil
}

var _ interfaces.LogStore = &StoreMock{}
