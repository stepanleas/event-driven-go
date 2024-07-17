package api

import (
	"context"
	"fmt"
)

type FilesApiMock struct {
	Files map[string]string
}

func (m FilesApiMock) UploadFile(ctx context.Context, fileID string, fileContent string) error {
	if m.Files == nil {
		m.Files = make(map[string]string)
	}

	if _, ok := m.Files[fileID]; ok {
		return fmt.Errorf("file %s is already saved", fileID)
	}

	m.Files[fileID] = fileContent

	return nil
}
