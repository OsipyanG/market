package memcached

import (
	"errors"
	"math/rand"
	"net"

	"github.com/OsipyanG/market/services/auth-msv/config"
	"github.com/OsipyanG/market/services/auth-msv/internal/cache"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/bradfitz/gomemcache/memcache"
)

var (
	ErrUnpackValue    = errors.New("can't unpack memcached value")
	ErrCloseMemcached = errors.New("can't close memcached")
)

type RefreshTokenCache struct {
	config *config.RefreshTokenConfig
	client *memcache.Client
}

func New(memcacheConf *config.MemcachedConfig, refreshConf *config.RefreshTokenConfig) (*RefreshTokenCache, error) {
	refreshCache := &RefreshTokenCache{
		config: refreshConf,
		client: memcache.New(net.JoinHostPort(memcacheConf.Host, memcacheConf.Port)),
	}

	if err := refreshCache.client.Ping(); err != nil {
		return nil, errwrap.Wrap(cache.ErrConnection, err)
	}

	return refreshCache, nil
}

func (rt *RefreshTokenCache) getExpiration() int32 {
	return int32(rt.config.Timeout.Seconds()) +
		rand.Int31n(int32(rt.config.Jitter.Seconds()))
}

func (rt *RefreshTokenCache) Close() error {
	return errwrap.WrapIfErr(ErrCloseMemcached, rt.client.Close())
}
