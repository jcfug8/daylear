package domain

import (
	"context"
	"fmt"
	"io"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	// "fmt"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
	uuid "github.com/satori/go.uuid"
)

// TODO: consolidate recipe image upload and circle image upload
// TODO: also need to remove old image on update ^^^

const CircleImageRoot = "circles"

// CreateCircle creates a new circle.
func (d *Domain) CreateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle) (model.Circle, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when creating a circle")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	circle.Id.CircleId = 0

	// Generate a handle if not provided
	if circle.Handle == "" {
		circle.Handle, _ = d.generateUniqueCircleHandle(ctx, circle.Title)
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to begin creating circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to begin creating circle"}
	}
	defer tx.Rollback()

	circle.ImageURI, err = d.createCircleImageURI(ctx, circle)
	if err != nil {
		log.Error().Err(err).Msg("unable to create circle image")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to create circle image"}
	}

	dbCircle, err := tx.CreateCircle(ctx, circle)
	if err != nil {
		log.Error().Err(err).Msg("unable to create circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to create circle"}
	}

	circleAccess := model.CircleAccess{
		CircleAccessParent: model.CircleAccessParent{
			CircleId: dbCircle.Id,
		},
		Requester: model.CircleRequester{
			UserId: authAccount.AuthUserId,
		},
		Recipient: model.UserId{
			UserId: authAccount.AuthUserId,
		},
		PermissionLevel: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
		State:           types.AccessState_ACCESS_STATE_ACCEPTED,
	}

	dbCircleAccess, err := tx.CreateCircleAccess(ctx, circleAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to create circle access")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to create circle access"}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("unable to finish creating circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to finish creating circle"}
	}

	dbCircle.CircleAccess = dbCircleAccess
	return dbCircle, nil
}

// DeleteCircle deletes a circle.
func (d *Domain) DeleteCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (circle model.Circle, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when deleting a circle")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	if id.CircleId == 0 {
		log.Warn().Msg("circle id required when deleting a circle")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	_, err = d.determineCircleAccess(ctx, authAccount, id, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when deleting a circle")
		return model.Circle{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to begin deleting circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to begin deleting circle"}
	}
	defer tx.Rollback()

	circle, err = tx.DeleteCircle(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to delete circle"}
	}

	if circle.ImageURI != "" {
		err := d.removeCircleImage(ctx, authAccount, circle.Id)
		if err != nil {
			log.Error().Err(err).Msg("unable to remove circle image")
			return model.Circle{}, domain.ErrInternal{Msg: "unable to remove circle image"}
		}
	}

	err = tx.BulkDeleteCircleAccess(ctx, model.CircleAccessParent{CircleId: id})
	if err != nil {
		log.Error().Err(err).Msg("unable to delete circle accesses")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to delete circle accesses"}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("unable to finish deleting circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to finish deleting circle"}
	}

	return circle, nil
}

// GetCircle gets a circle.
func (d *Domain) GetCircle(ctx context.Context, authAccount model.AuthAccount, parent model.CircleParent, id model.CircleId) (dbCircle model.Circle, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when getting a circle")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	if id.CircleId == 0 {
		log.Warn().Msg("circle id required when getting a circle")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	dbCircle, err = d.repo.GetCircle(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to get circle"}
	}

	dbCircle.CircleAccess, err = d.determineCircleAccess(
		ctx, authAccount, id,
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ),
		withResourceVisibilityLevel(dbCircle.VisibilityLevel),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when getting a circle")
		return model.Circle{}, err
	}

	return dbCircle, nil
}

// ListCircles lists circles for a parent.
func (d *Domain) ListCircles(ctx context.Context, authAccount model.AuthAccount, parent model.CircleParent, pageSize int32, pageOffset int64, filter string, fieldMask []string) (dbCircles []model.Circle, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when listing circles")
		return nil, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		_, err = d.determineUserAccess(ctx, authAccount, model.UserId{UserId: parent.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing circles")
			return nil, err
		}
	}

	dbCircles, err = d.repo.ListCircles(ctx, authAccount, pageSize, pageOffset, filter, fieldMask)
	if err != nil {
		log.Error().Err(err).Msg("unable to list circles")
		return nil, domain.ErrInternal{Msg: "unable to list circles"}
	}

	return dbCircles, nil
}

// UpdateCircle updates a circle.
func (d *Domain) UpdateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle, fields []string) (dbCircle model.Circle, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when updating a circle")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	if circle.Id.CircleId == 0 {
		log.Warn().Msg("circle id required when updating a circle")
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	_, err = d.determineCircleAccess(ctx, authAccount, circle.Id, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when updating a circle")
		return model.Circle{}, err
	}

	if slices.Contains(fields, model.CircleFields.Handle) && circle.Handle == "" {
		circle.Handle, _ = d.generateUniqueCircleHandle(ctx, circle.Title)
	}

	previousDbCircle, err := d.repo.GetCircle(ctx, authAccount, circle.Id)
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle when updating a circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to get circle"}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to begin updating circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to begin updating circle"}
	}
	defer tx.Rollback()

	if slices.Contains(fields, model.CircleFields.ImageURI) && circle.ImageURI != previousDbCircle.ImageURI {
		circle.ImageURI, err = d.updateCircleImageURI(ctx, authAccount, circle)
		if err != nil {
			log.Error().Err(err).Msg("unable to update circle image")
			return model.Circle{}, domain.ErrInternal{Msg: "unable to update circle image"}
		}
	}

	updated, err := d.repo.UpdateCircle(ctx, authAccount, circle, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to update circle"}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("unable to finish updating circle")
		return model.Circle{}, domain.ErrInternal{Msg: "unable to finish updating circle"}
	}

	return updated, nil
}

// UploadCircleImage uploads a circle image.
func (d *Domain) UploadCircleImage(ctx context.Context, authAccount model.AuthAccount, id model.CircleId, imageReader io.Reader) (imageURI string, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when uploading a circle image")
		return "", domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	if id.CircleId == 0 {
		log.Warn().Msg("circle id required when uploading a circle image")
		return "", domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	_, err = d.determineCircleAccess(ctx, authAccount, id, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when uploading a circle image")
		return "", err
	}

	circle, err := d.repo.GetCircle(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle when uploading a circle image")
		return "", domain.ErrInternal{Msg: "unable to get circle"}
	}
	oldImageURI := circle.ImageURI

	imageURI, err = d.uploadCircleImage(ctx, id, imageReader)
	if err != nil {
		log.Error().Err(err).Msg("unable to upload circle image")
		return "", domain.ErrInternal{Msg: "unable to upload circle image"}
	}

	_, err = d.repo.UpdateCircle(ctx, authAccount, model.Circle{
		Id:       id,
		ImageURI: imageURI,
	}, []string{model.CircleFields.ImageURI})
	if err != nil {
		log.Error().Err(err).Msg("unable to update circle image")
		return "", domain.ErrInternal{Msg: "unable to update circle image"}
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
	circle, err := d.GetCircle(ctx, authAccount, model.CircleParent{}, id)
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

// generateUniqueCircleHandle generates a unique handle for a circle based on the title or a random string
func (d *Domain) generateUniqueCircleHandle(ctx context.Context, title string) (string, error) {
	// Basic slugify: lowercase, replace spaces with dashes, remove non-alphanum
	base := slugify(title)
	if base == "" {
		base = "circle"
	}
	handle := base
	i := 1
	circleRepo, ok := d.repo.(interface {
		CircleHandleExists(context.Context, string) (bool, error)
	})
	if !ok {
		return "", fmt.Errorf("repo does not implement CircleHandleExists")
	}
	for {
		exists, _ := circleRepo.CircleHandleExists(ctx, handle)
		if !exists {
			return handle, nil
		}
		handle = base + "-" + strconv.Itoa(i)
		i++
	}
}

// slugify is a helper to create a URL-friendly string
func slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`(^-|-$)`).ReplaceAllString(s, "")
	return s
}
