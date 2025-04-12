package hxscheduler

// One per time interval.
func (loc *Location) GetOneSchedulingOption(params *ParamsCanRun) (OptionsSchedule, error) {
	intervalsSought := params.TimeInterval.BreakDown(params.TaskRun.EstimatedDuration)

	resourcesNeededPerType := params.TaskRun.GetNeededResourcesPerType()
	neededTypes := params.TaskRun.GetNeededResourceTypes()

	result := make([]*OptionSchedule, 0)

	for _, interval := range intervalsSought {
		intervalResourcesNeeded := make(ResourcesPerType)

		for _, neededType := range neededTypes {
			resourcesNeededPerCurrentType := resourcesNeededPerType[neededType]

			intervalResourcesPerCurrentType := make([]*ResourceScheduled, 0)

			var qty uint16

			for _, resource := range loc.Resources[neededType] {
				if isAvailable := resource.IsAvailableIn(&interval); !isAvailable {
					continue
				}

				intervalResourcesPerCurrentType = append(intervalResourcesPerCurrentType, resource)

				qty = qty + resource.ServedQuantity

				if qty == resourcesNeededPerCurrentType {
					break
				}
			}

			if qty < resourcesNeededPerCurrentType {
				break //interval cannot provide all resources
			}

			intervalResourcesNeeded[neededType] = intervalResourcesPerCurrentType
		}

		result = append(
			result,

			&OptionSchedule{
				WhenCanStart: interval.TimeStart,
				Resources:    intervalResourcesNeeded,
			},
		)
	}

	return result,
		nil
}
