package domain

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
	"golang.org/x/crypto/bcrypt"
)

// Common words for generating human-readable access keys
var accessKeyWords = []string{
	// Nature & Elements
	"apple", "banana", "cherry", "dragon", "eagle", "forest", "garden", "harbor",
	"island", "jungle", "knight", "lighthouse", "mountain", "ocean", "palace", "queen",
	"river", "sunset", "tiger", "umbrella", "village", "waterfall", "xylophone", "yellow",
	"zebra", "alpine", "breeze", "cascade", "dolphin", "emerald", "flamingo", "glacier",
	"horizon", "infinity", "jasmine", "kangaroo", "lagoon", "meadow", "nebula", "orchid",
	"phoenix", "quartz", "rainbow", "sapphire", "thunder", "unicorn", "volcano", "whisper",
	"xenon", "yacht", "zenith", "aurora", "blossom", "crystal", "diamond", "eclipse",
	"fountain", "galaxy", "harmony", "illusion", "jewel", "kingdom", "legend", "mystic",
	"nectar", "oasis", "paradise", "radiant", "serenity", "treasure", "utopia", "velvet",
	"wonder", "xanadu", "yearning", "zephyr",

	// Colors & Materials
	"amber", "azure", "bronze", "coral", "golden", "ivory", "jade", "lavender",
	"marble", "onyx", "pearl", "ruby", "silver", "topaz", "violet", "warmth",

	// Animals & Creatures
	"butterfly", "cardinal", "elephant", "falcon", "gazelle", "hummingbird", "leopard", "mermaid",
	"nightingale", "octopus", "penguin", "quokka", "raccoon", "seahorse", "toucan", "vulture",
	"walrus", "xenops", "yak", "zebra",

	// Places & Landmarks
	"acropolis", "bamboo", "castle", "desert", "everest", "fuji", "giza", "himalaya",
	"kilimanjaro", "louvre", "monaco", "nepal", "olympus", "paris", "quebec", "rome",
	"sahara", "taj", "venice", "washington", "yosemite", "zurich",

	// Emotions & Qualities
	"adventure", "bravery", "courage", "dignity", "elegance", "freedom", "grace", "honor",
	"imagination", "joy", "kindness", "love", "majesty", "nobility", "passion", "quiet",
	"respect", "strength", "tranquility", "unity", "wisdom", "zeal",

	// Objects & Tools
	"anchor", "basket", "compass", "drum", "envelope", "feather", "guitar", "hammer",
	"ink", "journal", "key", "lantern", "mirror", "notebook", "oar", "paintbrush",
	"quill", "ribbon", "sail", "telescope", "violin", "wreath",

	// Food & Drinks
	"avocado", "blueberry", "chocolate", "dragonfruit", "elderberry", "fig", "grapefruit", "honey",
	"kiwi", "lemon", "mango", "nectarine", "orange", "papaya", "quince", "raspberry",
	"strawberry", "tangerine", "vanilla", "watermelon",

	// Weather & Seasons
	"autumn", "blizzard", "drizzle", "frost", "hail", "mist", "rain", "snow",
	"spring", "summer", "thunder", "winter",

	// Time & Space
	"dawn", "dusk", "midnight", "noon", "twilight", "cosmos", "galaxy", "nebula",
	"planet", "solar", "stellar", "universe",
}

// generateHumanReadableKey creates a human-readable access key using words and a number
func generateHumanReadableKey() (string, error) {
	// Generate 3 random words and 1 random 5-digit number
	var words []string

	// Generate 3 random words
	for i := 0; i < 3; i++ {
		// Generate random index
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(accessKeyWords))))
		if err != nil {
			return "", err
		}
		words = append(words, accessKeyWords[randomIndex.Int64()])
	}

	// Generate a random 5-digit number
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return "", err
	}
	// Format as 5 digits with leading zeros if needed
	numberStr := fmt.Sprintf("%05d", randomNumber.Int64())

	// Insert the number at a random position (0-3)
	insertPosition, err := rand.Int(rand.Reader, big.NewInt(4))
	if err != nil {
		return "", err
	}

	// Insert the number at the random position
	words = append(words[:insertPosition.Int64()], append([]string{numberStr}, words[insertPosition.Int64():]...)...)

	// Format words (keep everything lowercase) and join with hyphens
	var formattedWords []string
	for _, word := range words {
		if len(word) > 0 {
			formattedWords = append(formattedWords, word)
		}
	}

	return strings.Join(formattedWords, "-"), nil
}

// generateRandomKey generates a traditional random key (fallback)
func generateRandomKey() (string, error) {
	randomBytes := make([]byte, 24) // Reduced from 32 for shorter keys
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

// CreateAccessKey creates a new access key
func (d *Domain) CreateAccessKey(ctx context.Context, authAccount model.AuthAccount, accessKey model.AccessKey) (dbAccessKey model.AccessKey, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required when creating an access key")
		return model.AccessKey{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	// Ensure the access key is being created for the authenticated user
	if accessKey.Parent.UserId == 0 {
		log.Warn().Msg("user id required when creating an access key")
		return model.AccessKey{}, domain.ErrInvalidArgument{Msg: "user id required when creating an access key"}
	}

	_, err = d.determineUserAccess(
		ctx, authAccount, model.UserId{UserId: accessKey.Parent.UserId},
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when creating an access key")
		return model.AccessKey{}, err
	}

	accessKey.AccessKeyId.AccessKeyId = 0

	// Generate a human-readable access key using words
	unencryptedKey, err := generateHumanReadableKey()
	if err != nil {
		log.Warn().Err(err).Msg("unable to generate human-readable key, falling back to random key")
		// Fallback to random key if word generation fails
		unencryptedKey, err = generateRandomKey()
		if err != nil {
			log.Error().Err(err).Msg("unable to generate access key")
			return model.AccessKey{}, domain.ErrInternal{Msg: "unable to generate access key"}
		}
	}

	// Encrypt the access key using bcrypt
	encryptedKey, err := bcrypt.GenerateFromPassword([]byte(unencryptedKey), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("unable to encrypt access key")
		return model.AccessKey{}, domain.ErrInternal{Msg: "unable to encrypt access key"}
	}

	// Set the encrypted key for storage
	accessKey.EncryptedAccessKey = string(encryptedKey)
	// Set the unencrypted key for return (only on creation)
	accessKey.UnencryptedAccessKey = unencryptedKey

	// Create the access key in the database
	dbAccessKey, err = d.repo.CreateAccessKey(ctx, accessKey, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to create access key")
		return model.AccessKey{}, domain.ErrInternal{Msg: "unable to create access key"}
	}

	// Ensure the unencrypted key is returned for the response
	dbAccessKey.UnencryptedAccessKey = unencryptedKey
	dbAccessKey.Parent = accessKey.Parent

	return dbAccessKey, nil
}

// DeleteAccessKey deletes an access key
func (d *Domain) DeleteAccessKey(ctx context.Context, authAccount model.AuthAccount, id model.AccessKeyId) (dbAccessKey model.AccessKey, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when deleting an access key")
		return model.AccessKey{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	// Get the access key to verify ownership
	existingAccessKey, err := d.repo.GetAccessKey(ctx, authAccount, id, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to get access key for deletion")
		return model.AccessKey{}, err
	}

	_, err = d.determineUserAccess(
		ctx, authAccount, model.UserId{UserId: existingAccessKey.Parent.UserId},
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when deleting an access key")
		return model.AccessKey{}, err
	}

	// Delete the access key
	dbAccessKey, err = d.repo.DeleteAccessKey(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete access key")
		return model.AccessKey{}, domain.ErrInternal{Msg: "unable to delete access key"}
	}

	dbAccessKey.Parent = existingAccessKey.Parent
	return dbAccessKey, nil
}

// GetAccessKey retrieves an access key
func (d *Domain) GetAccessKey(ctx context.Context, authAccount model.AuthAccount, parent model.AccessKeyParent, id model.AccessKeyId, fields []string) (dbAccessKey model.AccessKey, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when getting an access key")
		return model.AccessKey{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	userId := authAccount.AuthUserId

	// verify the user has access to the parent
	if parent.UserId != 0 {
		_, err := d.determineUserAccess(
			ctx, authAccount, model.UserId{UserId: userId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when getting an access key")
			return model.AccessKey{}, err
		}
		userId = parent.UserId
	}

	// Get the access key
	dbAccessKey, err = d.repo.GetAccessKey(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get access key")
		return model.AccessKey{}, err
	}

	// Verify the user owns this access key
	if dbAccessKey.Parent.UserId != userId {
		log.Warn().Msg("users can only get their own access keys")
		return model.AccessKey{}, domain.ErrPermissionDenied{Msg: "users can only get their own access keys"}
	}

	// Ensure unencrypted key is never returned
	dbAccessKey.UnencryptedAccessKey = ""

	return dbAccessKey, nil
}

// ListAccessKeys lists access keys for a user
func (d *Domain) ListAccessKeys(ctx context.Context, authAccount model.AuthAccount, userId int64, pageSize int32, offset int64, filter string, fields []string) ([]model.AccessKey, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when listing access keys")
		return nil, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	// verify the user has access to the parent
	_, err := d.determineUserAccess(
		ctx, authAccount, model.UserId{UserId: userId},
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when listing access keys")
		return nil, err
	}

	// List access keys
	accessKeys, err := d.repo.ListAccessKeys(ctx, authAccount, userId, pageSize, offset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list access keys")
		return nil, domain.ErrInternal{Msg: "unable to list access keys"}
	}

	// Ensure unencrypted keys are never returned
	for i := range accessKeys {
		accessKeys[i].UnencryptedAccessKey = ""
	}

	return accessKeys, nil
}

// UpdateAccessKey updates an access key
func (d *Domain) UpdateAccessKey(ctx context.Context, authAccount model.AuthAccount, accessKey model.AccessKey, fields []string) (dbAccessKey model.AccessKey, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when updating an access key")
		return model.AccessKey{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	// Get the existing access key to verify ownership
	existingAccessKey, err := d.repo.GetAccessKey(ctx, authAccount, accessKey.AccessKeyId, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to get access key for update")
		return model.AccessKey{}, err
	}

	_, err = d.determineUserAccess(
		ctx, authAccount, model.UserId{UserId: existingAccessKey.Parent.UserId},
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
		withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when updating an access key")
		return model.AccessKey{}, err
	}

	// Ensure only title and description can be updated
	accessKey.Parent = existingAccessKey.Parent
	accessKey.EncryptedAccessKey = existingAccessKey.EncryptedAccessKey
	accessKey.CreateTime = existingAccessKey.CreateTime

	// Update the access key
	dbAccessKey, err = d.repo.UpdateAccessKey(ctx, authAccount, accessKey, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update access key")
		return model.AccessKey{}, domain.ErrInternal{Msg: "unable to update access key"}
	}

	// Ensure unencrypted key is never returned
	dbAccessKey.UnencryptedAccessKey = ""
	dbAccessKey.Parent = existingAccessKey.Parent

	return dbAccessKey, nil
}

// AuthenticateByAccessKey authenticates a user by their access key
func (d *Domain) AuthenticateByAccessKey(ctx context.Context, username string, unencryptedAccessKey string) (model.User, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if username == "" || unencryptedAccessKey == "" {
		log.Warn().Msg("user id and access key are required for authentication")
		return model.User{}, domain.ErrNotFound{Msg: "user id and access key are required"}
	}

	filter := fmt.Sprintf("username = %s", username)

	// Get the user
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

	user := users[0]

	// Get all access keys for the user to check against
	accessKeys, err := d.repo.ListAccessKeys(ctx, model.AuthAccount{}, user.Id.UserId, 100, 0, "", []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to get access keys for authentication")
		return model.User{}, err
	}

	found := false

	// Check each access key against the provided unencrypted key
	for _, accessKey := range accessKeys {
		// Compare the unencrypted key with the stored encrypted key
		err := bcrypt.CompareHashAndPassword([]byte(accessKey.EncryptedAccessKey), []byte(unencryptedAccessKey))
		if err == nil {
			found = true
			break
		}
	}

	if !found {
		return model.User{}, domain.ErrNotFound{Msg: "invalid access key"}
	}

	return user, nil
}
