package domain

import (
	"sync"

	domainPort "github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/fileinspector"
	"github.com/jcfug8/daylear/server/ports/fileretriever"
	"github.com/jcfug8/daylear/server/ports/filestorage"
	"github.com/jcfug8/daylear/server/ports/recipescraper"
	"github.com/jcfug8/daylear/server/ports/repository"
	"github.com/jcfug8/daylear/server/ports/token"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var _ domainPort.Domain = &Domain{}

// DomainParams defines the dependencies for the domain layer.
type DomainParams struct {
	fx.In

	Log  zerolog.Logger
	Repo repository.Client

	TokenClient   token.Client
	ImageStore    filestorage.Client
	FileInspector fileinspector.Client
	FileRetriever fileretriever.Client

	RecipeScrapers       []recipescraper.HostSpecificClient `group:"recipescrapers"`
	DefaultRecipeScraper recipescraper.DefaultClient
}

// NewDomain creates a new domain.
func NewDomain(params DomainParams) domainPort.Domain {
	recipeScrapers := make(map[string]recipescraper.HostSpecificClient)
	for _, scraper := range params.RecipeScrapers {
		for _, host := range scraper.GetHost() {
			recipeScrapers[host] = scraper
		}
	}

	d := &Domain{
		log:  params.Log,
		repo: params.Repo,

		tokenClient:   params.TokenClient,
		tokenStore:    &sync.Map{},
		fileStore:     params.ImageStore,
		fileInspector: params.FileInspector,
		fileRetriever: params.FileRetriever,

		recipeScrapers:       recipeScrapers,
		defaultRecipeScraper: params.DefaultRecipeScraper,
	}
	return d
}

// Domain defines the domain layer or service.
type Domain struct {
	log  zerolog.Logger
	repo repository.Client

	tokenClient   token.Client
	tokenStore    *sync.Map
	fileStore     filestorage.Client
	fileInspector fileinspector.Client
	fileRetriever fileretriever.Client

	recipeScrapers       map[string]recipescraper.HostSpecificClient
	defaultRecipeScraper recipescraper.DefaultClient
}
