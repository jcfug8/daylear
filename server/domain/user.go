package domain

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"io"
	"path"
	"strconv"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"

	uuid "github.com/satori/go.uuid"
)

const UserImageRoot = "users"

func (d *Domain) IdentifyUser(ctx context.Context, user model.User) (model.User, error) {
	filter := ""
	switch {
	case user.AmazonId != "":
		filter = fmt.Sprintf(`%s = "%s"`, model.UserField_AmazonId, user.AmazonId)
	case user.GoogleId != "":
		filter = fmt.Sprintf(`%s = "%s"`, model.UserField_GoogleId, user.GoogleId)
	case user.FacebookId != "":
		filter = fmt.Sprintf(`%s = "%s"`, model.UserField_FacebookId, user.FacebookId)
	}

	users, err := d.repo.ListUsers(ctx, model.AuthAccount{}, 1, 0, filter, nil)
	if err != nil {
		return model.User{}, err
	}

	if len(users) > 1 {
		return model.User{}, domain.ErrInternal{Msg: "Multiple users found"}
	}

	if len(users) == 0 {
		return model.User{}, domain.ErrNotFound{Msg: "User not found"}
	}

	return users[0], nil
}

// ListUsers lists users.
func (d *Domain) ListUsers(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, pageSize int32, pageOffset int64, filter string, fieldMask []string) (users []model.User, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user_id required")
		return nil, domain.ErrInvalidArgument{Msg: "user_id required"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
		dbCircle, err := d.repo.GetCircle(ctx, authAccount, model.CircleId{CircleId: parent.CircleId}, []string{model.CircleField_Visibility})
		if err != nil {
			log.Error().Err(err).Msg("unable to get circle when listing calendars")
			return nil, domain.ErrInternal{Msg: "unable to get circle when listing calendars"}
		}
		_, err = d.determineCircleAccess(
			ctx, authAccount, model.CircleId{CircleId: parent.CircleId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
			withResourceVisibilityLevel(dbCircle.VisibilityLevel),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing calendars")
			return nil, err
		}
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
	}

	users, err = d.repo.ListUsers(ctx, authAccount, pageSize, pageOffset, filter, fieldMask)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates a user.
func (d *Domain) UpdateUser(ctx context.Context, authAccount model.AuthAccount, user model.User, updateMask []string) (dbUser model.User, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain UpdateUser called")

	if user.Id.UserId == 0 {
		log.Warn().Msg("id required")
		return model.User{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	_, err = d.determineUserAccess(
		ctx, authAccount, user.Id,
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when updating user")
		return model.User{}, err
	}

	dbUser, err = d.repo.UpdateUser(ctx, authAccount, user, updateMask)
	if err != nil {
		return model.User{}, err
	}

	return dbUser, nil
}

// GetUser gets a user.
func (d *Domain) GetUser(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, id model.UserId, fields []string) (dbUser model.User, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if id.UserId == 0 {
		log.Warn().Msg("id required")
		return model.User{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	dbUser, err = d.repo.GetUser(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetUser failed")
		return model.User{}, err
	}

	dbUser.UserAccess, err = d.determineUserAccess(
		ctx, authAccount, id,
		withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC),
		withAllowPendingAccess(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when getting user")
		return model.User{}, err
	}

	return dbUser, nil
}

// GetOwnUser gets the current user.
func (d *Domain) GetOwnUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId, fields []string) (model.User, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if id.UserId != authAccount.AuthUserId {
		log.Warn().Msg("id does not match auth account")
		return model.User{}, domain.ErrInvalidArgument{Msg: "id does not match auth account"}
	}

	dbUser, err := d.repo.GetUser(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetUser failed")
		return model.User{}, err
	}

	return dbUser, nil
}

// DeleteUser deletes a user.
func (d *Domain) DeleteUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.User, error) {
	return model.User{}, nil
}

// CreateUser creates a new user.
func (d *Domain) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	if user.Email == "" {
		return model.User{}, domain.ErrInvalidArgument{Msg: "email is required"}
	}

	// if the username is not set, determine it
	if user.Username == "" {
		username, err := d.determineUsername(ctx, user.Email)
		if err != nil {
			return model.User{}, err
		}
		user.Username = username
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		d.log.Error().Err(err).Msg("failed to begin transaction")
		return model.User{}, domain.ErrInternal{Msg: "failed to begin transaction"}
	}
	defer tx.Rollback()

	dbUser, err := tx.CreateUser(ctx, user, nil)
	if err != nil {
		d.log.Error().Err(err).Msg("failed to create user")
		return model.User{}, domain.ErrInternal{Msg: "failed to create user"}
	}

	userAccess := model.UserAccess{
		UserAccessParent: model.UserAccessParent{
			UserId: dbUser.Id,
		},
		Requester:       dbUser.Id,
		Recipient:       dbUser.Id,
		PermissionLevel: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
		State:           types.AccessState_ACCESS_STATE_ACCEPTED,
	}

	dbUserAccess, err := tx.CreateUserAccess(ctx, userAccess, nil)
	if err != nil {
		d.log.Error().Err(err).Msg("failed to create user access")
		return model.User{}, domain.ErrInternal{Msg: "failed to create user access"}
	}

	err = tx.Commit()
	if err != nil {
		d.log.Error().Err(err).Msg("failed to commit transaction")
		return model.User{}, domain.ErrInternal{Msg: "failed to commit transaction"}
	}

	dbUser.UserAccess = dbUserAccess

	return dbUser, nil
}

func (d *Domain) determineUsername(ctx context.Context, email string) (string, error) {
	emailParts := strings.Split(email, "@")
	if len(emailParts) != 2 {
		return "", domain.ErrInvalidArgument{Msg: "invalid email"}
	}
	baseUsername := emailParts[0]
	username := baseUsername

	for tryCount := 0; tryCount < 5; tryCount++ {
		filter := fmt.Sprintf("username = '%s'", username)
		dbUsers, err := d.repo.ListUsers(ctx, model.AuthAccount{}, 1, 0, filter, nil)
		if err != nil {
			d.log.Error().Err(err).Msg("failed to list users")
			return "", domain.ErrInternal{Msg: "failed to list users"}
		}
		if len(dbUsers) == 0 {
			break
		}

		username = fmt.Sprintf("%s_%d", baseUsername, rand.Intn(10000))
	}

	if username == "" {
		return "", domain.ErrInternal{Msg: "failed to determine username"}
	}

	return username, nil
}

// UploadUserImage uploads a user image and returns the image URI
func (d *Domain) UploadUserImage(ctx context.Context, authAccount model.AuthAccount, id model.UserId, imageReader io.Reader) (imageURI string, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "parent required"}
	}
	if id.UserId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "id required"}
	}

	_, err = d.determineUserAccess(ctx, authAccount, model.UserId{UserId: authAccount.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when uploading user image")
		return "", err
	}

	user, err := d.repo.GetUser(ctx, authAccount, id, []string{model.UserField_ImageUri})
	if err != nil {
		return "", err
	}
	oldImageURI := user.ImageUri

	imageURI, err = d.uploadUserImage(ctx, id, imageReader)
	if err != nil {
		return "", err
	}
	_, err = d.repo.UpdateUser(ctx, authAccount, model.User{
		Id:       id,
		ImageUri: imageURI,
	}, []string{model.UserField_ImageUri})
	if err != nil {
		return "", err
	}
	if oldImageURI != "" && oldImageURI != imageURI {
		go d.fileStore.DeleteFile(context.Background(), oldImageURI)
	}
	return imageURI, nil
}

func (d *Domain) uploadUserImage(ctx context.Context, id model.UserId, imageReader io.Reader) (imageURL string, err error) {
	image, err := d.imageClient.CreateImage(ctx, imageReader)
	if err != nil {
		return "", err
	}
	err = image.Convert(ctx, "jpg")
	if err != nil {
		return "", err
	}
	width, height, err := image.GetDimensions(ctx)
	if err != nil {
		return "", err
	}
	if width > maxImageWidth || height > maxImageHeight {
		newWidth, newHeight := resizeToFit(width, height, maxImageWidth)
		err = image.Resize(ctx, newWidth, newHeight)
		if err != nil {
			return "", err
		}
	}
	file, err := image.GetFile()
	if err != nil {
		return "", err
	}
	defer image.Remove(ctx)
	imagePath := path.Join(UserImageRoot, strconv.FormatInt(id.UserId, 10), uuid.NewV4().String())
	imagePath = imagePath + file.Extension
	imageURI, err := d.fileStore.UploadPublicFile(ctx, imagePath, file)
	if err != nil {
		return "", err
	}
	return imageURI, nil
}

// FavoriteUser adds a user to the current user's favorites.
func (d *Domain) FavoriteUser(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, id model.UserId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain FavoriteUser called")

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth_user_id required")
		return domain.ErrInvalidArgument{Msg: "auth_user_id required"}
	}

	if id.UserId == 0 {
		log.Warn().Msg("id required")
		return domain.ErrInvalidArgument{Msg: "id required"}
	}

	// Prevent self-favoriting
	if authAccount.AuthUserId == id.UserId {
		log.Warn().Msg("cannot favorite yourself")
		return domain.ErrInvalidArgument{Msg: "cannot favorite yourself"}
	}

	// If a user parent is specified, ensure the current user has admin access to that user
	if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		_, err := d.determineUserAccess(
			ctx, authAccount, model.UserId{UserId: authAccount.UserId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN),
		)
		if err != nil {
			log.Error().Err(err).Msg("insufficient permissions to favorite user in this context")
			return domain.ErrPermissionDenied{Msg: "insufficient permissions"}
		}
	}

	_, err := d.determineUserAccess(
		ctx, authAccount, model.UserId{UserId: authAccount.UserId},
		withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("insufficient permissions to favorite user in this context")
		return domain.ErrPermissionDenied{Msg: "insufficient permissions"}
	}

	// Create the favorite relationship (personal favorites only)
	err = d.repo.CreateUserFavorite(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to create user favorite")
		return domain.ErrInternal{Msg: "unable to create user favorite"}
	}

	log.Info().Msg("Domain FavoriteUser completed successfully")
	return nil
}

// UnfavoriteUser removes a user from the current user's favorites.
func (d *Domain) UnfavoriteUser(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, id model.UserId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain UnfavoriteUser called")

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth_user_id required")
		return domain.ErrInvalidArgument{Msg: "auth_user_id required"}
	}

	if authAccount.AuthUserId == id.UserId {
		log.Warn().Msg("cannot unfavorite yourself")
		return domain.ErrInvalidArgument{Msg: "cannot unfavorite yourself"}
	}

	if id.UserId == 0 {
		log.Warn().Msg("id required")
		return domain.ErrInvalidArgument{Msg: "id required"}
	}

	// If a user parent is specified, ensure the current user has admin access to that user
	if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		_, err := d.determineUserAccess(
			ctx, authAccount, model.UserId{UserId: authAccount.UserId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN),
		)
		if err != nil {
			log.Error().Err(err).Msg("insufficient permissions to unfavorite user in this context")
			return domain.ErrPermissionDenied{Msg: "insufficient permissions"}
		}
	}

	// Delete the favorite relationship (personal favorites only)
	err := d.repo.DeleteUserFavorite(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete user favorite")
		return domain.ErrInternal{Msg: "unable to delete user favorite"}
	}

	log.Info().Msg("Domain UnfavoriteUser completed successfully")
	return nil
}
