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

var _ PublicCircleNamer = &defaultPublicCircleNamer{}

// PublicCircleNamer is an interface for creating and validating circle names.
type PublicCircleNamer interface {
	Format(model.CircleId) (string, error)
	IsMatch(name string) bool

	Parse(string) (model.CircleId, error)
}

type publicCircleNamer struct {
	pattern  string
	varCount int
}

type defaultPublicCircleNamer struct {
	namer *publicCircleNamer
}

// NewPublicCircleNamer creates a new PublicCircleNamer.
func NewPublicCircleNamer() (PublicCircleNamer, error) {
	t := new(pb.PublicCircle)
	resourceOption := proto.GetExtension(
		t.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)
	patterns := resourceOption.Pattern
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", t)
	}

	namers := make([]*publicCircleNamer, 0, len(patterns))
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

		namers = append(namers, &publicCircleNamer{
			pattern:  pattern,
			varCount: varCount,
		})
	}

	return &defaultPublicCircleNamer{
		namer: namers[0],
	}, nil
}

// Format formats a circle name.
func (n *defaultPublicCircleNamer) Format(id model.CircleId) (string, error) {
	return resourcename.Sprint(n.namer.pattern, fmt.Sprintf("%v", id.CircleId)), nil
}

// IsMatch checks if a name matches the circle pattern.
func (n *defaultPublicCircleNamer) IsMatch(name string) bool {
	return resourcename.Match(n.namer.pattern, name)
}

// Parse parses a circle name.
func (n *defaultPublicCircleNamer) Parse(name string) (id model.CircleId, err error) {

	var circleIdStr string

	err = resourcename.Sscan(name, n.namer.pattern, &circleIdStr)
	if err != nil {
		return id, err
	}

	id.CircleId, err = strconv.ParseInt(circleIdStr, 10, 64)
	if err != nil {
		return id, fileretriever.ErrInvalidArgument{Msg: "invalid format"}
	}

	return id, nil
}
