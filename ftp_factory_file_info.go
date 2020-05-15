package main

import (
	"os"
	"time"
)

type junosFileInfo struct{}

func (f *junosFileInfo) Name() string       { return "" }
func (f *junosFileInfo) Size() int64        { return 0 }
func (f *junosFileInfo) Mode() os.FileMode  { return os.ModeTemporary }
func (f *junosFileInfo) ModTime() time.Time { return time.Now() }
func (f *junosFileInfo) IsDir() bool        { return true }
func (f *junosFileInfo) Sys() interface{}   { return nil }
func (f *junosFileInfo) Owner() string      { return "juniper" }
func (f *junosFileInfo) Group() string      { return "staff" }
