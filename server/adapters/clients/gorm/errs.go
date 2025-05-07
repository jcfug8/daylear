package gorm

import (
	"errors"
	"strings"

	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
)

func ConvertGormError(err error) error {
	if err == nil {
		return nil
	}

	msg := simplifyErrorMessage(err)

	for _, e := range []struct {
		gormErr    error
		msg    string
		err   error
	}{
		{gormErr: gorm.ErrPrimaryKeyRequired, msg: "primary key required", err: repository.ErrInvalidArgument{Msg:msg}},
		{gormErr: gorm.ErrInvalidField, msg: "invalid field", err: repository.ErrInvalidArgument{Msg:msg}},
		{gormErr: gorm.ErrEmptySlice, msg: "empty slice", err: repository.ErrInvalidArgument{Msg:msg}},
		{gormErr: gorm.ErrInvalidValue, msg: "invalid value", err: repository.ErrInvalidArgument{Msg:msg}},
		{gormErr: gorm.ErrInvalidValueOfLength, msg: "invalid value of length", err: repository.ErrInvalidArgument{Msg:msg}},

		{gormErr: gorm.ErrDuplicatedKey, msg: "duplicate key value violates unique constraint", err: repository.ErrNewAlreadyExists{Msg:msg}},

		{gormErr: gorm.ErrNotImplemented, msg: "not implemented", err: repository.ErrNewUnimplemented{Msg:msg}},

		{gormErr: gorm.ErrRecordNotFound, msg: "record not found", err: repository.ErrNotFound{Msg:msg}},

		{gormErr: gorm.ErrInvalidTransaction, msg: "invalid transaction", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrMissingWhereClause, msg: "missing WHERE clause", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrUnsupportedRelation, msg: "unsupported relation", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrModelValueRequired, msg: "model value required", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrModelAccessibleFieldsRequired, msg: "model accessible fields required", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrSubQueryRequired, msg: "subquery required", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrInvalidData, msg: "invalid data", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrUnsupportedDriver, msg: "unsupported driver", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrRegistered, msg: "registered", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrDryRunModeUnsupported, msg: "dry run mode unsupported", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrInvalidDB, msg: "invalid db", err: repository.ErrInternal{Msg:msg}},
		{gormErr: gorm.ErrPreloadNotAllowed, msg: "preload is not allowed", err: repository.ErrInternal{Msg:msg}},
	} {
		if errors.Is(err, e.gormErr) || strings.Contains(msg, e.msg) {
			return e.err
		}
	}

	return err
}

func simplifyErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	return strings.Split(err.Error(), ";")[0]
}
