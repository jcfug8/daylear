package gorm

import (
	"errors"
	"strings"

	"github.com/jcfug8/daylear/server/core/errz"

	"gorm.io/gorm"
)

// ErrzError wraps a gorm error with appropriate status codes.
func ErrzError(e *errz.Errorz, msg string, err error) error {
	if err == nil {
		return nil
	}

	err = simplifyErrorMessage(err)

	for _, e := range []struct {
		err    error
		msg    string
		code   errz.Code
		res    func(string, ...any) *errz.Error
		resMsg string
	}{
		{err: gorm.ErrPrimaryKeyRequired, msg: "primary key required", code: errz.InvalidArgument},
		{err: gorm.ErrInvalidField, msg: "invalid field", code: errz.InvalidArgument},
		{err: gorm.ErrEmptySlice, msg: "empty slice", code: errz.InvalidArgument},
		{err: gorm.ErrInvalidValue, msg: "invalid value", code: errz.InvalidArgument},
		{err: gorm.ErrInvalidValueOfLength, msg: "invalid value of length", code: errz.InvalidArgument},

		{err: gorm.ErrDuplicatedKey, msg: "duplicate key value violates unique constraint", res: errz.NewAlreadyExists, resMsg: "record already exists"},

		{err: gorm.ErrNotImplemented, msg: "not implemented", res: errz.NewUnimplemented, resMsg: "not implemented"},

		{err: gorm.ErrRecordNotFound, msg: "record not found", res: errz.NewNotFound, resMsg: "record not found"},

		{err: gorm.ErrInvalidTransaction, msg: "invalid transaction", code: errz.Internal},
		{err: gorm.ErrMissingWhereClause, msg: "missing WHERE clause", code: errz.Internal},
		{err: gorm.ErrUnsupportedRelation, msg: "unsupported relation", code: errz.Internal},
		{err: gorm.ErrModelValueRequired, msg: "model value required", code: errz.Internal},
		{err: gorm.ErrModelAccessibleFieldsRequired, msg: "model accessible fields required", code: errz.Internal},
		{err: gorm.ErrSubQueryRequired, msg: "subquery required", code: errz.Internal},
		{err: gorm.ErrInvalidData, msg: "invalid data", code: errz.Internal},
		{err: gorm.ErrUnsupportedDriver, msg: "unsupported driver", code: errz.Internal},
		{err: gorm.ErrRegistered, msg: "registered", code: errz.Internal},
		{err: gorm.ErrDryRunModeUnsupported, msg: "dry run mode unsupported", code: errz.Internal},
		{err: gorm.ErrInvalidDB, msg: "invalid db", code: errz.Internal},
		{err: gorm.ErrPreloadNotAllowed, msg: "preload is not allowed", code: errz.Internal},
	} {
		if errors.Is(err, e.err) || strings.Contains(err.Error(), e.msg) {
			if e.res != nil {
				err = e.res(e.resMsg)
			} else {
				err = errz.Wrapf(e.msg).WithCode(e.code)
			}

			break
		}
	}

	if msg == "" {
		return e.Wrap(err)
	}

	return e.Wrapf("%s: %v", msg, err)
}

func simplifyErrorMessage(err error) error {
	if err == nil {
		return nil
	}

	return errz.Wrapf(strings.Split(err.Error(), ";")[0])
}
