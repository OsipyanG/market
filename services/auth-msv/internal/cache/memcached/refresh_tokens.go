package memcached

import (
	"github.com/OsipyanG/market/services/auth-msv/internal/cache"
	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/bradfitz/gomemcache/memcache"
)

// Creating a new key value pair. If an entry with such a refresh token already exists,
// then it is overwritten (and the deadline is reset to the initial value).
func (rt *RefreshTokenCache) SetRefreshToken(refreshToken string, claims *model.JWTClaims) error {
	value := pack(claims)

	err := rt.client.Set(&memcache.Item{
		Key:        refreshToken,
		Value:      []byte(value),
		Expiration: rt.getExpiration(),
	})

	return errwrap.WrapIfErr(cache.ErrSetRefreshToken, err)
}

func (rt *RefreshTokenCache) GetJWTClaims(refreshToken string) (*model.JWTClaims, error) {
	item, err := rt.client.Get(refreshToken)
	if err != nil {
		return nil, errwrap.Wrap(cache.ErrGetJWTClaims, err)
	}

	claims, err := unpack(string(item.Value))
	if err != nil {
		return nil, errwrap.Wrap(cache.ErrGetJWTClaims, err)
	}

	return claims, nil
}

func (rt *RefreshTokenCache) DeleteRefreshToken(refreshToken string) error {
	err := rt.client.Delete(refreshToken)

	return errwrap.WrapIfErr(cache.ErrDeleteRefreshToken, err)
}
