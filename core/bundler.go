package core

type Item struct {
	Title string
	Link  string
}

type Bundler struct {
	items      []Item
	runtimeDir string
}

func NewBundler(items []Item, runtimeDir string) *Bundler {
	return &Bundler{
		items:      items,
		runtimeDir: runtimeDir,
	}
}

func (b *Bundler) GetItems() []Item {
	return b.items
}
