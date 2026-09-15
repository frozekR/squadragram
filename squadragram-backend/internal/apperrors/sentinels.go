package apperrors

import "squadraton-backend/pkg/xerrors"

var (
	ErrInvalidRole  = xerrors.NewSentinel("invalid role")
	ErrInvalidSkill = xerrors.NewSentinel("invalid skill")
	ErrInvalidOwner = xerrors.NewSentinel("invalid owner")
	ErrInvalidMedia = xerrors.NewSentinel("invalid media")
)
