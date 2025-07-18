package domain

import (
	"sync"

	domainPort "github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/fileretriever"
	"github.com/jcfug8/daylear/server/ports/filestorage"
	"github.com/jcfug8/daylear/server/ports/image"
	"github.com/jcfug8/daylear/server/ports/imagegenerator"
	"github.com/jcfug8/daylear/server/ports/recipescraper"
	"github.com/jcfug8/daylear/server/ports/repository"
	"github.com/jcfug8/daylear/server/ports/token"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var _ domainPort.Domain = &Domain{}

const (
	maxImageWidth  = 1000
	maxImageHeight = 1000
)

// DomainParams defines the dependencies for the domain layer.
type DomainParams struct {
	fx.In

	Log  zerolog.Logger
	Repo repository.Client

	TokenClient    token.Client
	ImageStore     filestorage.Client
	ImageClient    image.Client
	FileRetriever  fileretriever.Client
	ImageGenerator imagegenerator.Client
	RecipeScraper  recipescraper.Client
}

// NewDomain creates a new domain.
func NewDomain(params DomainParams) domainPort.Domain {
	d := &Domain{
		log:  params.Log,
		repo: params.Repo,

		tokenClient:    params.TokenClient,
		tokenStore:     &sync.Map{},
		fileStore:      params.ImageStore,
		imageClient:    params.ImageClient,
		fileRetriever:  params.FileRetriever,
		imageGenerator: params.ImageGenerator,
		recipeScraper:  params.RecipeScraper,
	}
	return d
}

// Domain defines the domain layer or service.
type Domain struct {
	log  zerolog.Logger
	repo repository.Client

	tokenClient    token.Client
	tokenStore     *sync.Map
	fileStore      filestorage.Client
	imageClient    image.Client
	fileRetriever  fileretriever.Client
	imageGenerator imagegenerator.Client
	recipeScraper  recipescraper.Client
}
