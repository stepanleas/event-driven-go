package api

import (
	"context"
	"sync"
)

type SpreadsheetsAPIClientMock struct {
	mock sync.Mutex

	rows map[string][]string
}

func (m *SpreadsheetsAPIClientMock) AppendRow(ctx context.Context, sheetName string, row []string) error {
	m.mock.Lock()
	defer m.mock.Unlock()

	m.rows[sheetName] = row

	return nil
}
