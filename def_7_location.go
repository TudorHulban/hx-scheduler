package hxscheduler

import (
	"sync"

	goerrors "github.com/TudorHulban/go-errors"
	"github.com/asaskevich/govalidator"
)

type Location struct {
	Name      string
	Resources ResourcesPerType

	mu sync.Mutex

	ID             int64
	LocationOffset int64
}

type ParamsNewLocation struct {
	Name string `valid:"required"`

	ID             int64 `valid:"required"`
	LocationOffset int64
}

func NewLocation(params *ParamsNewLocation) (*Location, error) {
	if _, errValidation := govalidator.ValidateStruct(params); errValidation != nil {
		return nil,
			goerrors.ErrServiceValidation{
				ServiceName: "Organigram",
				Caller:      "NewLocation",
				Issue:       errValidation,
			}
	}

	return &Location{
			ID:             params.ID,
			Name:           params.Name,
			LocationOffset: params.LocationOffset,
			Resources:      make(ResourcesPerType),
		},
		nil
}

func (loc *Location) AddResource(resource *ResourceInfo) {
	loc.mu.Lock()
	defer loc.mu.Unlock()

	if _, exists := loc.Resources[resource.ResourceType]; exists {
		loc.Resources[resource.ResourceType] = append(
			loc.Resources[resource.ResourceType],
			&ResourceScheduled{
				ResourceInfo: *resource,
				schedule:     make(map[TimeInterval]RunID),
			},
		)

		return
	}

	loc.Resources[resource.ResourceType] = []*ResourceScheduled{
		{
			ResourceInfo: *resource,
			schedule:     make(map[TimeInterval]RunID),
		},
	}
}
