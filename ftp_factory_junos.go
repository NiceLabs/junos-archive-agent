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
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"github.com/thoas/go-funk"

	"goftp.io/server"
)

var (
	errUnable     = errors.New("junos: not supported")
	errFileName   = errors.New("junos: filename format error")
	reFileName    = regexp.MustCompile(`^(?P<name>.+)_\d{8}_\d{6}_juniper(?P<suffix>\.conf(?:\.\d+)?)$`)
	reNameIndex   = funk.IndexOf(reFileName.SubexpNames(), "name")
	reSuffixIndex = funk.IndexOf(reFileName.SubexpNames(), "suffix")
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

	filePath := d.makeFilePath(destPath)
	if filePath == "" {
		err = errFileName
		return
	}

	ctx := context.Background()
	content, _, _, _ := d.client.Repositories.GetContents(ctx, d.Owner, d.Repo, filePath, nil)
	opts := &github.RepositoryContentFileOptions{Message: &destPath, Content: block}
	if content != nil {
		opts.SHA = content.SHA
		if d.isChanged(content, string(opts.Content)) {
			log.Println("not changed, ignored")
			return
		}
	}
	_, _, err = d.client.Repositories.UpdateFile(ctx, d.Owner, d.Repo, filePath, opts)
	return
}

func (d *junosDriver) isChanged(content *github.RepositoryContent, original string) bool {
	payload, err := content.GetContent()
	if err != nil {
		return false
	}
	return strings.EqualFold(
		payload[strings.Index(payload, "\n"):],
		original[strings.Index(original, "\n"):],
	)
}

func (d *junosDriver) decompress(data io.Reader) ([]byte, error) {
	reader, err := gzip.NewReader(data)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func (d *junosDriver) makeFilePath(name string) string {
	base, fileName := path.Split(name)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	matched := reFileName.FindStringSubmatch(fileName)
	if len(matched) < 2 {
		return ""
	}
	routerName := matched[reNameIndex]
	suffix := matched[reSuffixIndex]
	return path.Join(d.Prefix, base, routerName+suffix)
}
