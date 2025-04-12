package hxscheduler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllOptionsPerTimeInterval(t *testing.T) {
	location := Location{
		ID:   1,
		Name: t.Name(),

		Resources: ResourcesPerType{
			1: []*ResourceScheduled{
				{
					ResourceInfo: ResourceInfo{
						ID:              1,
						Name:            "Resource 1",
						CostPerLoadUnit: map[uint8]float32{1: 2.0},
						ResourceType:    1,
						ServedQuantity:  1,
					},

					schedule: map[TimeInterval]RunID{
						{TimeStart: now, TimeEnd: now + halfHour}:                     Maintenance,
						{TimeStart: now + oneHour, TimeEnd: now + oneHour + halfHour}: Maintenance,
					},
				},
				{
					ResourceInfo: ResourceInfo{
						ID:              2,
						Name:            "Resource 2",
						CostPerLoadUnit: map[uint8]float32{1: 3.0},
						ResourceType:    1,
						ServedQuantity:  1,
					},

					schedule: map[TimeInterval]RunID{
						{TimeStart: now + oneHour, TimeEnd: now + oneHour + halfHour}: Maintenance,
					},
				},
				{
					ResourceInfo: ResourceInfo{
						ID:              3,
						Name:            "Resource 3",
						CostPerLoadUnit: map[uint8]float32{1: 2.0},
						ResourceType:    1,
						ServedQuantity:  1,
					},

					schedule: map[TimeInterval]RunID{
						{TimeStart: now, TimeEnd: now + halfHour}: Maintenance,
					},
				},
			},
			2: []*ResourceScheduled{
				{
					ResourceInfo: ResourceInfo{
						ID:              4,
						Name:            "Resource 4",
						CostPerLoadUnit: map[uint8]float32{1: 1.0},
						ResourceType:    2,
						ServedQuantity:  1,
					},

					schedule: map[TimeInterval]RunID{},
				},
			},
		},
	}

	taskRun := Run{
		ID:                1,
		EstimatedDuration: halfHour,

		Dependencies: []RunDependency{
			{
				ResourceType:     1,
				ResourceQuantity: 1,
			},
			{
				ResourceType:     2,
				ResourceQuantity: 1,
			},
		},

		RunLoad: RunLoad{
			Load:     1,
			LoadUnit: 1,
		},
	}

	options, errGetOptions := location.GetAllSchedulingOptions(
		&ParamsCanRun{
			TimeInterval: TimeInterval{
				TimeStart: now,
				TimeEnd:   now + 2*oneHour,
			},

			TaskRun: &taskRun,

			PossibilitiesUpTo: 2,
		},
	)
	require.NoError(t, errGetOptions)
	require.NotEmpty(t, options)
	require.Len(t,
		options,
		6,
	)
	require.NotEmpty(t,
		options[0].Resources,
	)

	fmt.Println(
		options.String(&taskRun),
	)
}
