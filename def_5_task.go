package hxscheduler

import "fmt"

type RunDependency struct {
	ResourceType     ResourceType
	ResourceQuantity uint8
}

type RunLoad struct {
	Load     float32
	LoadUnit uint8
}

type Run struct {
	Name         string
	Dependencies []RunDependency

	RunLoad

	ID                RunID
	InitiatorID       int64
	EstimatedDuration int64
}

func (r *Run) GetRunCostUsing(resource *ResourceInfo) (float32, error) {
	costPerUnit, ok := resource.CostPerLoadUnit[r.RunLoad.LoadUnit]
	if !ok {
		return 0,
			fmt.Errorf(
				"resource does not support load unit %d",
				r.RunLoad.LoadUnit,
			)
	}

	return r.RunLoad.Load * costPerUnit,
		nil
}

func (r *Run) GetNeededResourceTypes() []ResourceType {
	resourceTypes := make(map[ResourceType]bool)

	for _, dependency := range r.Dependencies {
		resourceTypes[dependency.ResourceType] = true
	}

	result := make([]ResourceType, len(resourceTypes), len(resourceTypes))

	var ix uint16

	for rt := range resourceTypes {
		result[ix] = rt

		ix++
	}

	return result
}

func (r *Run) GetNeededResourcesPerType() map[ResourceType]uint16 {
	result := make(map[ResourceType]uint16)

	for _, dependency := range r.Dependencies {
		if currentNumberNeeded, exists := result[dependency.ResourceType]; exists {
			result[dependency.ResourceType] = currentNumberNeeded + uint16(dependency.ResourceQuantity)

			continue
		}

		result[dependency.ResourceType] = uint16(dependency.ResourceQuantity)
	}

	return result
}
