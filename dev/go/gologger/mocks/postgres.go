package mocks

import (
	"errors"
	"time"

	"github.com/modell-aachen/gologger/interfaces"
)

type mockData struct {
	Stored   []interfaces.LogRow
	IsClosed bool
}

type PostgresMock struct {
	Data *mockData
}

func (instance PostgresMock) Store(row interfaces.LogRow) (err error) {
	if instance.Data.IsClosed {
		return errors.New("Handle has been closed")
	}
	instance.Data.Stored = append(instance.Data.Stored, row)
	return nil
}

func (instance PostgresMock) Read(startTime time.Time, endTime time.Time, source interfaces.SourceString, levels []interfaces.LevelString) (rows []interfaces.LogRow, err error) {
	if instance.Data.IsClosed {
		return rows, errors.New("Handle has been closed")
	}
	var filteredRows []interfaces.LogRow
	for _, row := range instance.Data.Stored {
		if row.Time.After(startTime) && row.Time.Before(endTime) {
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

func (instance PostgresMock) Close() {
	instance.Data.IsClosed = true
}

func (instance PostgresMock) CleanUp() (err error) {
	if instance.Data.IsClosed {
		return errors.New("Handle has been closed")
	}

	instance.Data.Stored = make([]interfaces.LogRow, 0)
	return nil
}
