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
		filter = fmt.Sprintf(`%s = "%s"`, model.UserFields.AmazonId, user.AmazonId)
	case user.GoogleId != "":
		filter = fmt.Sprintf(`%s = "%s"`, model.UserFields.GoogleId, user.GoogleId)
	case user.FacebookId != "":
		filter = fmt.Sprintf(`%s = "%s"`, model.UserFields.FacebookId, user.FacebookId)
	}

	users, err := d.repo.ListUsers(ctx, model.AuthAccount{}, 1, 0, filter, model.UserFields.Mask())
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

	if parent.CircleId != 0 && authAccount.CircleId != 0 && parent.CircleId != authAccount.CircleId {
		log.Warn().Msg("both circle ids set but do not match")
		return nil, domain.ErrInvalidArgument{Msg: "both circle ids set but do not match"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
	} else {
		authAccount.UserId = authAccount.AuthUserId
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getCircleAccessLevels failed")
			return nil, err
		}
	} else if authAccount.UserId != 0 {
		authAccount.PermissionLevel, err = d.getUserAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getUserAccessLevels failed")
			return nil, err
		}
	}

	users, err = d.repo.ListUsers(ctx, authAccount, pageSize, pageOffset, filter, fieldMask)
	if err != nil {
		return nil, err
	}

	// manually add the user access for the current user
	for i, user := range users {
		if user.Id.UserId == authAccount.AuthUserId {
			users[i].UserAccess = model.UserAccess{
				UserAccessParent: model.UserAccessParent{
					UserId: user.Id,
				},
				Requester: model.UserId{
					UserId: authAccount.AuthUserId,
				},
				Recipient: model.UserId{
					UserId: authAccount.AuthUserId,
				},
				Level: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
			}
		}
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

	authAccount.UserId = user.Id.UserId

	authAccount.PermissionLevel, err = d.getUserAccessLevels(ctx, authAccount)
	if err != nil {
		log.Error().Err(err).Msg("getUserAccessLevels failed")
		return model.User{}, err
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_ADMIN {
		log.Warn().Msg("user does not have admin permission")
		return model.User{}, domain.ErrPermissionDenied{Msg: "user does not have admin permission"}
	}

	dbUser, err = d.repo.UpdateUser(ctx, authAccount, user, updateMask)
	if err != nil {
		return model.User{}, err
	}

	return dbUser, nil
}

// GetUser gets a user.
func (d *Domain) GetUser(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, id model.UserId, fieldMask []string) (dbUser model.User, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if id.UserId == 0 {
		log.Warn().Msg("id required")
		return model.User{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.UserId = id.UserId

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
	}

	authAccount.PermissionLevel, err = d.checkUserAccess(ctx, authAccount, id, 0)
	if err != nil {
		log.Error().Err(err).Msg("getUserAccessLevels failed")
		return model.User{}, err
	}

	dbUser, err = d.repo.GetUser(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetUser failed")
		return model.User{}, err
	}

	// manually add the user access for the current user
	if authAccount.AuthUserId == authAccount.UserId {
		dbUser.UserAccess = model.UserAccess{
			UserAccessParent: model.UserAccessParent{
				UserId: id,
			},
			Requester: model.UserId{
				UserId: authAccount.AuthUserId,
			},
			Recipient: model.UserId{
				UserId: authAccount.AuthUserId,
			},
			Level: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
		}
	}

	return dbUser, nil
}

// GetOwnUser gets the current user.
func (d *Domain) GetOwnUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId, fieldMask []string) (model.User, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if id.UserId != authAccount.AuthUserId {
		log.Warn().Msg("id does not match auth account")
		return model.User{}, domain.ErrInvalidArgument{Msg: "id does not match auth account"}
	}

	authAccount.UserId = authAccount.AuthUserId

	dbUser, err := d.repo.GetUser(ctx, authAccount, id)
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

	dbUser, err := d.repo.CreateUser(ctx, user)
	if err != nil {
		d.log.Error().Err(err).Msg("failed to create user")
		return model.User{}, domain.ErrInternal{Msg: "failed to create user"}
	}

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
	if authAccount.AuthUserId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "parent required"}
	}
	if id.UserId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.UserId = id.UserId

	authAccount.PermissionLevel, err = d.getUserAccessLevels(ctx, authAccount)
	if err != nil {
		return "", err
	}
	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_ADMIN {
		return "", domain.ErrPermissionDenied{Msg: "user does not have access"}
	}
	user, err := d.repo.GetUser(ctx, authAccount, id)
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
	}, []string{model.UserFields.ImageUri})
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
