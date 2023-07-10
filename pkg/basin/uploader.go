package basin

import (
	"github.com/brunocalza/basin/pkg/backup"
	zlog "github.com/rs/zerolog/log"
)

// Uploader uploads to backups to Basin.
type Uploader struct {
	ch    chan backup.BackupResult
	basin *Basin
}

// NewUploader creates a new Basin implementation.
func NewUploader(ch chan backup.BackupResult, basin *Basin) *Uploader {
	return &Uploader{
		ch:    ch,
		basin: basin,
	}
}

func (u *Uploader) Start() {
	for backupResult := range u.ch {
		id, err := u.basin.Upload(backupResult.Path)
		if err != nil {
			zlog.Error().Err(err).Msg("upload to basin")
			continue
		}

		zlog.Info().Str("cid", id.String()).Msg("uploaded to basin")
	}
}
