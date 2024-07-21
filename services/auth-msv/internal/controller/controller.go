package controller

import "errors"

var (
	ErrNoJWTClaims  = errors.New("no jwt claims")
	ErrNoPayload    = errors.New("no payload (no data)")
	ErrAccessDenied = errors.New("access is denied")
)
