package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateList creates a new list.
func (d *Domain) CreateList(ctx context.Context, authAccount model.AuthAccount, list model.List) (dbList model.List, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return model.List{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if list.Parent.CircleId != 0 && authAccount.CircleId != 0 && list.Parent.CircleId != authAccount.CircleId {
		log.Warn().Msg("both circle ids set but do not match")
		return model.List{}, domain.ErrInvalidArgument{Msg: "both circle ids set but do not match"}
	}

	if list.Parent.CircleId != 0 {
		authAccount.CircleId = list.Parent.CircleId
	} else if list.Parent.UserId != 0 {
		authAccount.UserId = list.Parent.UserId
	}

	list.Id.ListId = 0

	if list.Parent.CircleId != 0 {
		_, err = d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: list.Parent.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine circle access")
			return model.List{}, err
		}
	} else if list.Parent.UserId != 0 {
		_, err = d.determineUserAccess(ctx, authAccount, model.UserId{UserId: authAccount.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user access")
			return model.List{}, err
		}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.List{}, err
	}

	dbList, err = d.repo.CreateList(ctx, authAccount, list)
	if err != nil {
		log.Error().Err(err).Msg("repo.CreateList failed")
		return model.List{}, err
	}

	listAccess := model.ListAccess{
		ListAccessParent: model.ListAccessParent{
			ListId: dbList.Id,
		},
		PermissionLevel: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
		State:           types.AccessState_ACCESS_STATE_ACCEPTED,
		Requester: model.ListRecipientOrRequester{
			UserId: authAccount.AuthUserId,
		},
	}
	if authAccount.CircleId != 0 {
		listAccess.Recipient = model.ListRecipientOrRequester{
			CircleId: authAccount.CircleId,
		}
	} else {
		listAccess.Recipient = model.ListRecipientOrRequester{
			UserId: authAccount.AuthUserId,
		}
	}

	dbListAccess, err := tx.CreateListAccess(ctx, listAccess, nil)
	if err != nil {
		log.Error().Err(err).Msg("tx.CreateListAccess failed")
		return model.List{}, err
	}

	dbList.ListAccess = dbListAccess

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.List{}, err
	}

	dbList.Parent = list.Parent

	return dbList, nil
}

// GetList gets a list.
func (d *Domain) GetList(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, id model.ListId, fields []string) (list model.List, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("parent required")
		return model.List{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.ListId == 0 {
		log.Warn().Msg("id required")
		return model.List{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	list, err = d.repo.GetList(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetList failed")
		return model.List{}, err
	}

	list.ListAccess, err = d.determineListAccess(
		ctx, authAccount, id,
		withResourceVisibilityLevel(list.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
		withAllowPendingAccess(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access")
		return model.List{}, err
	}

	return list, nil
}

// ListLists lists lists.
func (d *Domain) ListLists(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, pageSize int32, pageOffset int32, filter string, fields []string) (lists []model.List, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user_id required")
		return nil, domain.ErrInvalidArgument{Msg: "user_id required"}
	}

	authAccount.PermissionLevel = types.PermissionLevel_PERMISSION_LEVEL_ADMIN

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
		dbCircle, err := d.repo.GetCircle(ctx, authAccount, model.CircleId{CircleId: parent.CircleId}, []string{model.CircleField_Visibility})
		if err != nil {
			log.Error().Err(err).Msg("unable to get circle when listing lists")
			return nil, domain.ErrInternal{Msg: "unable to get circle when listing lists"}
		}
		determinedCircleAccess, err := d.determineCircleAccess(
			ctx, authAccount, model.CircleId{CircleId: parent.CircleId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
			withResourceVisibilityLevel(dbCircle.VisibilityLevel),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing lists")
			return nil, err
		}
		authAccount.PermissionLevel = determinedCircleAccess.GetPermissionLevel()
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		determinedUserAccess, err := d.determineUserAccess(
			ctx, authAccount, model.UserId{UserId: authAccount.UserId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
			withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing lists")
			return nil, err
		}
		authAccount.PermissionLevel = determinedUserAccess.GetPermissionLevel()
	}

	lists, err = d.repo.ListLists(ctx, authAccount, pageSize, pageOffset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.ListLists failed")
		return nil, err
	}

	return lists, nil
}

// UpdateList updates a list.
func (d *Domain) UpdateList(ctx context.Context, authAccount model.AuthAccount, list model.List, fields []string) (dbList model.List, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("parent required")
		return model.List{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if list.Id.ListId == 0 {
		log.Warn().Msg("id required")
		return model.List{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	previousDbList, err := d.repo.GetList(ctx, authAccount, list.Id, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetList failed")
		return model.List{}, err
	}

	_, err = d.determineListAccess(
		ctx, authAccount, list.Id,
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access")
		return model.List{}, err
	}

	previousSectionMap := make(map[int64]model.ListSection)
	for _, section := range previousDbList.Sections {
		previousSectionMap[section.Id] = section
	}

	for i, section := range list.Sections {
		if section.Title == "" {
			log.Warn().Msg("title is required when updating a list section")
			return model.List{}, domain.ErrInvalidArgument{Msg: "title is required"}
		}
		if _, ok := previousSectionMap[section.Id]; !ok && section.Id != 0 {
			log.Warn().Msg("invalid section id when updating a list")
			return model.List{}, domain.ErrInvalidArgument{Msg: "invalid section id"}
		}

		if section.Id == 0 {
			section.Id = dbList.CreateTime.UnixMilli() + int64(i)
		}
	}

	dbList, err = d.repo.UpdateList(ctx, authAccount, list, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateList failed")
		return model.List{}, err
	}

	dbList.Parent = list.Parent

	return dbList, nil
}

// DeleteList deletes a list.
func (d *Domain) DeleteList(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, id model.ListId) (err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("parent required")
		return domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.ListId == 0 {
		log.Warn().Msg("id required")
		return domain.ErrInvalidArgument{Msg: "id required"}
	}

	_, err = d.determineListAccess(
		ctx, authAccount, id,
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access")
		return err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return err
	}

	_, err = d.repo.DeleteList(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.DeleteList failed")
		return err
	}

	// TODO:delete accesses
	err = d.repo.BulkDeleteListAccess(ctx, model.ListAccessParent{ListId: id})
	if err != nil {
		log.Error().Err(err).Msg("unable to bulk delete list accesses")
		return err
	}

	// TODO:delete items
	err = d.repo.BulkDeleteListItems(ctx, model.ListItemParent{ListId: id})
	if err != nil {
		log.Error().Err(err).Msg("unable to bulk delete list items")
		return err
	}
	// TODO:delete favorites
	err = d.repo.BulkDeleteListFavorites(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to bulk delete list favorites")
		return err
	}
	// TODO:delete completions

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return err
	}

	return nil
}

// FavoriteList adds a list to the user's favorites.
func (d *Domain) FavoriteList(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, id model.ListId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
		_, err := d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: authAccount.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when favoriting a list")
			return err
		}
		// TODO: we make want to check if the circle has access to the list as well
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		_, err := d.determineUserAccess(ctx, authAccount, model.UserId{UserId: authAccount.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when favoriting a list")
			return err
		}
		// TODO: we make want to check if the user has access to the list as well
	}

	dbList, err := d.repo.GetList(ctx, authAccount, id, []string{model.ListField_VisibilityLevel})
	if err != nil {
		log.Error().Err(err).Msg("unable to get list for favoriting")
		return err
	}

	_, err = d.determineListAccess(
		ctx, authAccount, id,
		withResourceVisibilityLevel(dbList.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access")
		return err
	}

	err = d.repo.CreateListFavorite(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to create list favorite")
		return err
	}

	log.Info().Int64("listId", id.ListId).Msg("list favorited successfully")
	return nil
}

// UnfavoriteList removes a list from the user's favorites.
func (d *Domain) UnfavoriteList(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, id model.ListId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
		_, err := d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: authAccount.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when unfavoriting a list")
			return err
		}
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		_, err := d.determineUserAccess(ctx, authAccount, model.UserId{UserId: authAccount.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when unfavoriting a list")
			return err
		}
	}

	err := d.repo.DeleteListFavorite(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete list favorite")
		return err
	}

	log.Info().Int64("listId", id.ListId).Msg("list unfavorited successfully")
	return nil
}
