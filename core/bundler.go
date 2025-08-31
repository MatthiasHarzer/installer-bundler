package core

import "io/fs"

type Item struct {
	Title string
	Link  string
}

type Bundler struct {
	items            []Item
	runtimeProjectFS fs.FS
}

func NewBundler(items []Item, runtimeProject fs.FS) *Bundler {
	return &Bundler{
		items:            items,
		runtimeProjectFS: runtimeProject,
	}
}

func (b *Bundler) GetItems() []Item {
	return b.items
}
