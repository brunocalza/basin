package basin

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

type Basin struct {
	w3sClient w3s.Client
}

type SnapshotItem struct {
	Created time.Time
	Cid     cid.Cid
}

type StatusInfo struct {
	Created time.Time
}

func NewBasin(token string) (*Basin, error) {
	client, err := w3s.NewClient(w3s.WithToken(token))
	if err != nil {
		return nil, fmt.Errorf("w3s new client: %s", err)
	}

	return &Basin{
		w3sClient: client,
	}, nil
}

func (b *Basin) Upload(path string) (cid.Cid, error) {
	file, err := os.Open(path)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("open deal file: %s", err)
	}

	id, err := b.w3sClient.Put(context.Background(), file)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("uploading to web3 storage: %s", err)
	}

	return id, nil
}

func (b *Basin) List() ([]SnapshotItem, error) {
	items := make([]SnapshotItem, 0)

	it, err := b.w3sClient.List(context.Background(), []w3s.ListOption{}...)
	if err != nil {
		return []SnapshotItem{}, fmt.Errorf("list: %s", err)
	}

	for {
		s, err := it.Next()
		if err != nil {
			break
		}

		items = append(items, SnapshotItem{
			Created: s.Created,
			Cid:     s.Cid,
		})

	}

	return items, nil
}

func (b *Basin) DownloadURL(cid cid.Cid) (string, error) {
	r, err := b.w3sClient.Get(context.Background(), cid)
	if err != nil {
		return "", fmt.Errorf("uploading to web3 storage: %s", err)
	}
	if r.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status not ok: %d", r.StatusCode)
	}

	_, fsys, err := r.Files()
	if err != nil {
		return "", fmt.Errorf("files: %s", err)
	}

	var path string
	_ = fs.WalkDir(fsys, "/", func(p string, d fs.DirEntry, err error) error {
		path = p
		return nil
	})

	return fmt.Sprintf("https://%s.ipfs.w3s.link/ipfs/%s%s", cid, cid, path), nil
}

func (b *Basin) Status(cid cid.Cid) (*w3s.Status, error) {
	status, err := b.w3sClient.Status(context.Background(), cid)
	if err != nil {
		return nil, fmt.Errorf("list: %s", err)
	}

	return status, nil
}
