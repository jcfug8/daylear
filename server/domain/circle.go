package domain

import (
	"context"
	"fmt"
	"io"
	"path"
	"strconv"

	// "fmt"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
	uuid "github.com/satori/go.uuid"
)

// TODO: consolidate recipe image upload and circle image upload

const CircleImageRoot = "circles"

// CreateCircle creates a new circle.
func (d *Domain) CreateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle) (model.Circle, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain CreateCircle called")
	if authAccount.UserId == 0 {
		log.Warn().Msg("parent required: missing user id")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	circle.Id.CircleId = 0

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Circle{}, err
	}
	defer tx.Rollback()

	circle.ImageURI, err = d.createCircleImageURI(ctx, circle)
	if err != nil {
		log.Error().Err(err).Msg("createCircleImageURI failed")
		return model.Circle{}, err
	}

	dbCircle, err := tx.CreateCircle(ctx, circle)
	if err != nil {
		log.Error().Err(err).Msg("CreateCircle in repo failed")
		return model.Circle{}, err
	}

	circleAccess := model.CircleAccess{
		CircleAccessParent: model.CircleAccessParent{
			CircleId: dbCircle.Id,
		},
		Requester: model.AuthAccount{
			UserId: authAccount.UserId,
		},
		Recipient:       authAccount.UserId,
		PermissionLevel: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
		State:           types.AccessState_ACCESS_STATE_ACCEPTED,
	}

	dbCircleAccess, err := tx.CreateCircleAccess(ctx, circleAccess)
	if err != nil {
		log.Error().Err(err).Msg("CreateCircleAccess in repo failed")
		return model.Circle{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Circle{}, err
	}

	dbCircle.CircleAccess = dbCircleAccess
	log.Info().Msg("Domain CreateCircle returning successfully")
	return dbCircle, nil
}

// DeleteCircle deletes a circle.
func (d *Domain) DeleteCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (circle model.Circle, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain DeleteCircle called")
	if authAccount.UserId == 0 || id.CircleId == 0 {
		log.Warn().Msg("parent and id required")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent and id required"}
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
	if err != nil {
		log.Error().Err(err).Msg("getCircleAccessLevels failed")
		return model.Circle{}, err
	}

	if authAccount.PermissionLevel <= types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		log.Warn().Msg("user does not have write permission")
		return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Circle{}, err
	}
	defer tx.Rollback()

	circle, err = tx.DeleteCircle(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("DeleteCircle in repo failed")
		return model.Circle{}, err
	}

	err = tx.BulkDeleteCircleAccess(ctx, model.CircleAccessParent{CircleId: id})
	if err != nil {
		log.Error().Err(err).Msg("BulkDeleteCircleAccess in repo failed")
		return model.Circle{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Circle{}, err
	}

	log.Info().Msg("Domain DeleteCircle returning successfully")
	return circle, nil
}

// GetCircle gets a circle.
func (d *Domain) GetCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (circle model.Circle, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain GetCircle called")
	if authAccount.UserId == 0 {
		log.Warn().Msg("parent required: missing user id")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.CircleId == 0 {
		log.Warn().Msg("id required: missing circle id")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
	if err != nil {
		log.Error().Err(err).Msg("getCircleAccessLevels failed")
		return model.Circle{}, err
	}

	if authAccount.VisibilityLevel != types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC && authAccount.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		log.Warn().Msg("user does not have write permission")
		return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}

	circle, err = d.repo.GetCircle(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("GetCircle in repo failed")
		return model.Circle{}, err
	}

	log.Info().Msg("Domain GetCircle returning successfully")
	return circle, nil
}

// ListCircles lists circles for a parent.
func (d *Domain) ListCircles(ctx context.Context, authAccount model.AuthAccount, pageSize int32, pageOffset int64, filter string, fieldMask []string) ([]model.Circle, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain ListCircles called")
	if authAccount.UserId == 0 {
		log.Warn().Msg("parent required: missing user id")
		return nil, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	circles, err := d.repo.ListCircles(ctx, authAccount, pageSize, pageOffset, filter, fieldMask)
	if err != nil {
		log.Error().Err(err).Msg("ListCircles in repo failed")
		return nil, err
	}

	log.Info().Msg("Domain ListCircles returning successfully")
	return circles, nil
}

// UpdateCircle updates a circle.
func (d *Domain) UpdateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle, updateMask []string) (dbCircle model.Circle, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain UpdateCircle called")
	if authAccount.UserId == 0 {
		log.Warn().Msg("parent required: missing user id")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if circle.Id.CircleId == 0 {
		log.Warn().Msg("id required: missing circle id")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.CircleId = circle.Id.CircleId

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
	if err != nil {
		log.Error().Err(err).Msg("getCircleAccessLevels failed")
		return model.Circle{}, err
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		log.Warn().Msg("user does not have write permission")
		return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}

	previousDbCircle, err := d.repo.GetCircle(ctx, authAccount, circle.Id)
	if err != nil {
		log.Error().Err(err).Msg("GetCircle in repo failed")
		return model.Circle{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Circle{}, err
	}
	defer tx.Rollback()

	for _, updateMaskField := range updateMask {
		if updateMaskField == model.CircleFields.ImageURI && circle.ImageURI != previousDbCircle.ImageURI {
			circle.ImageURI, err = d.updateCircleImageURI(ctx, authAccount, circle)
			if err != nil {
				log.Error().Err(err).Msg("updateCircleImageURI failed")
				return model.Circle{}, err
			}
		}
	}

	updated, err := d.repo.UpdateCircle(ctx, authAccount, circle, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("UpdateCircle in repo failed")
		return model.Circle{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Circle{}, err
	}

	log.Info().Msg("Domain UpdateCircle returning successfully")
	return updated, nil
}

// UploadCircleImage uploads a circle image.
func (d *Domain) UploadCircleImage(ctx context.Context, authAccount model.AuthAccount, id model.CircleId, imageReader io.Reader) (imageURI string, err error) {
	if authAccount.UserId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.CircleId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.CircleId = id.CircleId

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
	if err != nil {
		return "", err
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return "", domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	circle, err := d.repo.GetCircle(ctx, authAccount, id)
	if err != nil {
		return "", err
	}
	oldImageURI := circle.ImageURI

	imageURI, err = d.uploadCircleImage(ctx, id, imageReader)
	if err != nil {
		return "", err
	}

	_, err = d.repo.UpdateCircle(ctx, authAccount, model.Circle{
		Id:       id,
		ImageURI: imageURI,
	}, []string{model.CircleFields.ImageURI})
	if err != nil {
		return "", err
	}

	go d.fileStore.DeleteFile(context.Background(), oldImageURI)

	return imageURI, nil
}

// Helper methods for image handling
func (d *Domain) createCircleImageURI(ctx context.Context, circle model.Circle) (string, error) {
	if circle.ImageURI == "" {
		return "", nil
	}

	fileContents, err := d.fileRetriever.GetFileContents(ctx, circle.ImageURI)
	if err != nil {
		return "", err
	}
	defer fileContents.Close()

	imageURI, err := d.uploadCircleImage(ctx, circle.Id, fileContents)
	if err != nil {
		return "", err
	}

	return imageURI, nil
}

func (d *Domain) updateCircleImageURI(ctx context.Context, authAccount model.AuthAccount, circle model.Circle) (string, error) {
	if circle.ImageURI == "" {
		err := d.removeCircleImage(ctx, authAccount, circle.Id)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	fileContents, err := d.fileRetriever.GetFileContents(ctx, circle.ImageURI)
	if err != nil {
		return "", err
	}
	defer fileContents.Close()

	imageURI, err := d.uploadCircleImage(ctx, circle.Id, fileContents)
	if err != nil {
		return "", err
	}

	return imageURI, nil
}

func (d *Domain) uploadCircleImage(ctx context.Context, id model.CircleId, imageReader io.Reader) (imageURL string, err error) {
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

	imagePath := path.Join(CircleImageRoot, strconv.FormatInt(id.CircleId, 10), uuid.NewV4().String())
	imagePath = fmt.Sprintf("%s%s", imagePath, file.Extension)

	imageURI, err := d.fileStore.UploadPublicFile(ctx, imagePath, file)
	if err != nil {
		return "", err
	}

	return imageURI, nil
}

func (d *Domain) removeCircleImage(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (err error) {
	circle, err := d.GetCircle(ctx, authAccount, id)
	if err != nil {
		return err
	}

	if circle.ImageURI == "" {
		return nil
	}

	err = d.fileStore.DeleteFile(ctx, circle.ImageURI)
	if err != nil {
		return err
	}

	return nil
}
