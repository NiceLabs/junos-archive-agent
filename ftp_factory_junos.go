package main

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"

	"goftp.io/server"
)

var (
	errUnable = errors.New("github: not supported")
)

type junosDriver struct {
	*GitHubConfigure
	client *github.Client
}

func (d *junosDriver) NewDriver() (server.Driver, error)                   { return d, nil }
func (d *junosDriver) Stat(string) (server.FileInfo, error)                { return &junosFileInfo{}, nil }
func (d *junosDriver) ListDir(string, func(server.FileInfo) error) error   { return errUnable }
func (d *junosDriver) DeleteDir(string) error                              { return errUnable }
func (d *junosDriver) DeleteFile(string) error                             { return errUnable }
func (d *junosDriver) Rename(string, string) error                         { return errUnable }
func (d *junosDriver) MakeDir(string) error                                { return errUnable }
func (d *junosDriver) GetFile(string, int64) (int64, io.ReadCloser, error) { return 0, nil, errUnable }

func (d *junosDriver) PutFile(destPath string, data io.Reader, appendData bool) (n int64, err error) {
	if appendData {
		err = errUnable
		return
	}
	block, err := d.decompress(data)
	if err != nil {
		return
	}
	n = int64(len(block))
	storagePath := d.makeStoragePath(destPath)

	ctx := context.Background()
	opts := &github.RepositoryContentFileOptions{Message: &destPath, Content: block}
	content, _, _, _ := d.client.Repositories.GetContents(ctx, d.Owner, d.Repo, storagePath, nil)
	if content != nil {
		opts.SHA = content.SHA
	}
	_, _, err = d.client.Repositories.UpdateFile(context.Background(), d.Owner, d.Repo, storagePath, opts)
	return
}

func (d *junosDriver) decompress(data io.Reader) ([]byte, error) {
	reader, err := gzip.NewReader(data)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func (d *junosDriver) makeStoragePath(name string) string {
	base, fileName := path.Split(name)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	routerName := strings.SplitN(fileName, "_", 3)[0]
	log.Println(base, name, fileName, routerName)
	return path.Join(d.Prefix, base, routerName+filepath.Ext(fileName))
}
