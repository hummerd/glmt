package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gitlab.com/gitlab-merge-tool/glmt/internal/gerr"
	"gitlab.com/gitlab-merge-tool/glmt/internal/team"
)

type LocalSource struct {
	path string
}

func (s *LocalSource) Team(ctx context.Context) (*team.Team, error) {
	const maxBodySize = 1 << 20

	f, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read team file: %w", err)
	}

	defer func() { err = gerr.NewMultiError(err, f.Close()) }()

	var t team.Team
	err = json.NewDecoder(io.LimitReader(f, maxBodySize)).Decode(&t)
	if err != nil {
		return nil, fmt.Errorf("decoding body: %w", err)
	}

	return &t, nil
}
