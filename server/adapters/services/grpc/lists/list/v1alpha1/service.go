package v1alpha1

import (
	fieldmask "github.com/jcfug8/daylear/server/core/fieldmask"
	namer "github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/lists/list/v1alpha1"
	domain "github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// NewListServiceParams defines the dependencies for the ListService.
type NewListServiceParams struct {
	fx.In

	Domain              domain.Domain
	Log                 zerolog.Logger
	ListFieldMasker     fieldmask.FieldMasker `name:"v1alpha1ListFieldMasker"`
	AccessFieldMasker   fieldmask.FieldMasker `name:"v1alpha1ListAccessFieldMasker"`
	ListItemFieldMasker fieldmask.FieldMasker `name:"v1alpha1ListItemFieldMasker"`
	ListNamer           namer.ReflectNamer    `name:"v1alpha1ListNamer"`
	ListSectionNamer    namer.ReflectNamer    `name:"v1alpha1ListSectionNamer"`
	ListItemNamer       namer.ReflectNamer    `name:"v1alpha1ListItemNamer"`
	AccessNamer         namer.ReflectNamer    `name:"v1alpha1ListAccessNamer"`
	UserNamer           namer.ReflectNamer    `name:"v1alpha1UserNamer"`
	CircleNamer         namer.ReflectNamer    `name:"v1alpha1CircleNamer"`
}

// NewListService creates a new ListService.
func NewListService(params NewListServiceParams) (*ListService, error) {
	return &ListService{
		domain:              params.Domain,
		log:                 params.Log,
		listFieldMasker:     params.ListFieldMasker,
		accessFieldMasker:   params.AccessFieldMasker,
		listItemFieldMasker: params.ListItemFieldMasker,
		listNamer:           params.ListNamer,
		listSectionNamer:    params.ListSectionNamer,
		listItemNamer:       params.ListItemNamer,
		userNamer:           params.UserNamer,
		accessNamer:         params.AccessNamer,
		circleNamer:         params.CircleNamer,
	}, nil
}

// ListService defines the grpc handlers for the ListService.
type ListService struct {
	pb.UnimplementedListServiceServer
	pb.UnimplementedListAccessServiceServer
	pb.UnimplementedListItemServiceServer
	domain              domain.Domain
	log                 zerolog.Logger
	listFieldMasker     fieldmask.FieldMasker
	accessFieldMasker   fieldmask.FieldMasker
	listItemFieldMasker fieldmask.FieldMasker
	listNamer           namer.ReflectNamer
	listSectionNamer    namer.ReflectNamer
	listItemNamer       namer.ReflectNamer
	userNamer           namer.ReflectNamer
	circleNamer         namer.ReflectNamer
	accessNamer         namer.ReflectNamer
}
