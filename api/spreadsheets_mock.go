package api

import (
	"context"
	"sync"
)

type SpreadsheetsAPIClientMock struct {
	mock sync.Mutex

	Rows map[string][][]string
}

func (m *SpreadsheetsAPIClientMock) AppendRow(ctx context.Context, sheetName string, row []string) error {
	m.mock.Lock()
	defer m.mock.Unlock()

	if m.Rows == nil {
		m.Rows = make(map[string][][]string)
	}

	m.Rows[sheetName] = append(m.Rows[sheetName], row)

	return nil
}
