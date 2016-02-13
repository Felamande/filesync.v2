package syncer

import (
	"github.com/Felamande/filesync.v2/uri"
)

type cacher map[string]uri.Uri

func (c cacher) Get(u string) (uri.Uri, error) {
	if urig, exist := c[u]; exist {
		return urig, nil
	}
	urig, err := uri.Parse(u)
	if err != nil {
		return nil, err
	}
	c[u] = urig
	return urig, nil
	// c[u] = u
}
