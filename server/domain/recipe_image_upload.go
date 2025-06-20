package domain

import (
	"context"
	"fmt"
	"io"
	"path"
	"strconv"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/domain"
	uuid "github.com/satori/go.uuid"
)

const (
	RecipeImageRoot = "recipe-images"
)

// UploadRecipeImage -
func (d *Domain) UploadRecipeImage(ctx context.Context, parent model.RecipeParent, id model.RecipeId, imageReader io.Reader) (imageURI string, err error) {
	if parent.UserId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "id required"}
	}

	recipient, err := d.repo.GetRecipeRecipient(ctx, parent, id)
	if err != nil {
		return "", err
	}
	if recipient.PermissionLevel != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
		return "", domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}
	if parent.CircleId != 0 {
		permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
		if err != nil {
			return "", err
		}
		if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return "", domain.ErrPermissionDenied{Msg: "circle does not have write permission"}
		}

	}

	recipe, err := d.GetRecipe(ctx, parent, id, []string{model.RecipeFields.ImageURI})
	if err != nil {
		return "", err
	}
	oldImageURI := recipe.ImageURI

	imageURI, err = d.uploadRecipeImage(ctx, id, imageReader)
	if err != nil {
		return "", err
	}

	_, err = d.repo.UpdateRecipe(ctx, model.Recipe{
		Parent:   parent,
		Id:       id,
		ImageURI: imageURI,
	}, []string{model.RecipeFields.ImageURI})
	if err != nil {
		return "", err
	}

	go d.fileStore.DeleteFile(context.Background(), oldImageURI)

	return imageURI, nil
}

func (d *Domain) uploadRecipeImage(ctx context.Context, id model.RecipeId, imageReader io.Reader) (imageURL string, err error) {
	file, err := d.fileInspector.Inspect(ctx, imageReader)
	if err != nil {
		return "", err
	}
	defer file.Close()

	imagePath := path.Join(RecipeImageRoot, strconv.FormatInt(id.RecipeId, 10), uuid.NewV4().String())
	imagePath = fmt.Sprintf("%s%s", imagePath, file.Extension)

	iamgeURI, err := d.fileStore.UploadPublicFile(ctx, imagePath, file)
	if err != nil {
		return "", err
	}

	return iamgeURI, nil
}
