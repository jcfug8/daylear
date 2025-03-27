package namer

import (
	"fmt"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"

	"github.com/jcfug8/daylear/server/core/errz"
	"go.einride.tech/aip/resourcename"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	// IRIOMO:CUSTOM_CODE_SLOT_START userServiceNamerImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

var _ UserNamer = &defaultUserNamer{}

// UserNamer is an interface for creating and validating user names.
type UserNamer interface {
	Format(model.UserId) (string, error)
	IsMatch(name string) bool

	Parse(string) (model.UserId, error)
}

type userNamer struct {
	pattern  string
	varCount int
}

type defaultUserNamer struct {
	// IRIOMO:CUSTOM_CODE_SLOT_START UserNamerFields
	namer *userNamer
	// IRIOMO:CUSTOM_CODE_SLOT_END
}

// NewUserNamer creates a new UserNamer.
func NewUserNamer() (UserNamer, error) {
	t := new(pb.User)
	resourceOption := proto.GetExtension(
		t.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)
	patterns := resourceOption.Pattern
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", t)
	}

	namers := make([]*userNamer, 0, len(patterns))
	for _, pattern := range patterns {
		var patternScanner resourcename.Scanner
		var varCount int
		patternScanner.Init(patterns[0])
		for patternScanner.Scan() {
			segment := patternScanner.Segment()
			if segment.IsVariable() {
				varCount++
			}
		}

		namers = append(namers, &userNamer{
			pattern:  pattern,
			varCount: varCount,
		})
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START UserNamers
	return &defaultUserNamer{
		namer: namers[0],
	}, nil
	// IRIOMO:CUSTOM_CODE_SLOT_END
}

// Format formats a user name.
func (n *defaultUserNamer) Format(id model.UserId) (string, error) {
	return resourcename.Sprint(n.namer.pattern, fmt.Sprintf("%v", id.UserId)), nil
}

// IsMatch checks if a name matches the user pattern.
func (n *defaultUserNamer) IsMatch(name string) bool {
	return resourcename.Match(n.namer.pattern, name)
}

// Parse parses a user name.
func (n *defaultUserNamer) Parse(name string) (id model.UserId, err error) {

	var userIdStr string

	err = resourcename.Sscan(name, n.namer.pattern, &userIdStr)
	if err != nil {
		return id, err
	}

	id.UserId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return id, errz.NewInvalidArgument("invalid format")
	}

	return id, nil
}
