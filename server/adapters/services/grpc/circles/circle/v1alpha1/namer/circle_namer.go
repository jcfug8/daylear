package namer

import (
	"fmt"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"github.com/jcfug8/daylear/server/ports/fileretriever"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	circleParentSegmentCount = 2
)

var _ CircleNamer = &defaultCircleNamer{}

// CircleNamer is an interface for creating and validating circle names.
type CircleNamer interface {
	Format(model.CircleParent, model.CircleId) (string, error)
	IsMatch(name string) bool
	IsParent(parent string) bool
	Parse(string) (model.CircleParent, model.CircleId, error)
	ParseParent(string) (model.CircleParent, error)
}

type circleNamer struct {
	pattern  string
	varCount int
}

type defaultCircleNamer struct {
	// IRIOMO:CUSTOM_CODE_SLOT_START CircleNamerFields
	namer *circleNamer
	// IRIOMO:CUSTOM_CODE_SLOT_END
}

// NewCircleNamer creates a new CircleNamer.
func NewCircleNamer() (CircleNamer, error) {
	t := new(pb.Circle)
	resourceOption := proto.GetExtension(
		t.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)
	patterns := resourceOption.Pattern
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", t)
	}

	namers := make([]*circleNamer, 0, len(patterns))
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

		namers = append(namers, &circleNamer{
			pattern:  pattern,
			varCount: varCount,
		})
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START CircleNamers
	return &defaultCircleNamer{
		namer: namers[0],
	}, nil
	// IRIOMO:CUSTOM_CODE_SLOT_END
}

// Format formats a circle name.
func (n *defaultCircleNamer) Format(parent model.CircleParent, id model.CircleId) (string, error) {
	return resourcename.Sprint(n.namer.pattern, fmt.Sprintf("%v", parent.UserId), fmt.Sprintf("%v", id.CircleId)), nil
}

// IsMatch checks if a name matches the circle pattern.
func (n *defaultCircleNamer) IsMatch(name string) bool {
	return resourcename.Match(n.namer.pattern, name)
}

// IsParent checks if a name matches the circle parent pattern.
func (n *defaultCircleNamer) IsParent(parent string) bool {
	isParent := false
	foundSegments := 1
	resourcename.RangeParents(n.namer.pattern, func(p string) bool {
		if resourcename.Match(p, parent) && circleParentSegmentCount == foundSegments {
			isParent = true
			return false
		}
		foundSegments++
		return true
	})
	return isParent
}

// Parse parses a circle name.
func (n *defaultCircleNamer) Parse(name string) (parent model.CircleParent, id model.CircleId, err error) {

	var userIdStr string

	var circleIdStr string

	err = resourcename.Sscan(name, n.namer.pattern, &userIdStr, &circleIdStr)
	if err != nil {
		return parent, id, err
	}

	parent.UserId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return parent, id, fileretriever.ErrInvalidArgument{Msg: "invalid parent format"}
	}
	id.CircleId, err = strconv.ParseInt(circleIdStr, 10, 64)
	if err != nil {
		return parent, id, fileretriever.ErrInvalidArgument{Msg: "invalid format"}
	}

	return parent, id, nil
}

// ParseParent parses a circle parent name.
func (n *defaultCircleNamer) ParseParent(name string) (parent model.CircleParent, err error) {
	if !n.IsParent(name) {
		return parent, fmt.Errorf("invalid parent %s", name)
	}

	var userIdStr string

	resourcename.RangeParents(n.namer.pattern, func(p string) bool {
		if !resourcename.Match(p, name) {
			return true
		}

		err = resourcename.Sscan(name, p, &userIdStr)
		if err != nil {
			return false
		}

		return false
	})

	parent.UserId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return parent, fileretriever.ErrInvalidArgument{Msg: "invalid parent format"}
	}

	return parent, err
}
