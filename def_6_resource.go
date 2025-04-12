package hxscheduler

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	goerrors "github.com/TudorHulban/go-errors"
)

type ResourceInfo struct {
	Name            string
	CostPerLoadUnit map[uint8]float32 // load unit | cost per unit
	ID              int
	ResourceType    ResourceType
	ServedQuantity  uint16 // ex. apartment w 2 rooms serves 2, room serves 1
}

func (info ResourceInfo) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("ID: %d,", info.ID))
	sb.WriteString(fmt.Sprintf("Name: %q,", info.Name))
	sb.WriteString(fmt.Sprintf("ResourceType: %d", info.ResourceType))

	return sb.String()
}

// ResourceScheduled is mutex protected through Location ops.
type ResourceScheduled struct {
	ResourceInfo

	mu sync.RWMutex

	schedule map[TimeInterval]RunID
}

type ParamsNewResource struct {
	Name            string
	CostPerLoadUnit map[uint8]float32
	ID              int
	ResourceType    uint8
}

func (param *ParamsNewResource) IsValid() error {
	if len(param.Name) == 0 {
		return goerrors.ErrValidation{
			Caller: "IsValid - ParamsNewResource",
			Issue: goerrors.ErrNilInput{
				InputName: "Name",
			},
		}
	}

	if param.ResourceType <= 0 {
		return goerrors.ErrValidation{
			Caller: "IsValid - ParamsNewResource",
			Issue: goerrors.ErrInvalidInput{
				InputName: "ResourceType",
			},
		}
	}

	if param.CostPerLoadUnit == nil {
		return goerrors.ErrValidation{
			Caller: "IsValid - ParamsNewResource",
			Issue: goerrors.ErrNilInput{
				InputName: "CostPerLoadUnit",
			},
		}
	}

	for _, cost := range param.CostPerLoadUnit {
		if cost < 0 {
			return goerrors.ErrValidation{
				Caller: "IsValid - ParamsNewResource",
				Issue: goerrors.ErrNegativeInput{
					InputName: "CostPerLoadUnit",
				},
			}
		}
	}

	return nil
}

func NewResource(params *ParamsNewResource) (*ResourceScheduled, error) {
	if errValidation := params.IsValid(); errValidation != nil {
		return nil,
			errValidation
	}

	return &ResourceScheduled{
			ResourceInfo: ResourceInfo{
				ID:           params.ID,
				Name:         params.Name,
				ResourceType: ResourceType(params.ResourceType),

				CostPerLoadUnit: params.CostPerLoadUnit,
			},

			schedule: make(map[TimeInterval]RunID),
		},
		nil
}

func (res *ResourceScheduled) GetSchedule() string {
	res.mu.RLock()

	if len(res.schedule) == 0 {
		res.mu.RUnlock()

		return "Schedule: (empty)"
	}

	// Extract and sort intervals
	intervals := make([]TimeInterval, 0, len(res.schedule))

	for interval := range res.schedule {
		intervals = append(intervals, interval)
	}

	res.mu.RUnlock()

	sort.Slice(
		intervals,
		func(i, j int) bool {
			return intervals[i].TimeStart < intervals[j].TimeStart
		},
	)

	var sb strings.Builder
	sb.WriteString("Schedule:\n")

	for _, interval := range intervals {
		taskID := res.schedule[interval]

		sb.WriteString(
			fmt.Sprintf(
				"- [%d-%d] (UTC %d-%d) Offset %.1fh → Task %d\n",

				interval.TimeStart,
				interval.TimeEnd,
				interval.GetUTCTimeStart(),
				interval.GetUTCTimeEnd(),
				float64(interval.SecondsOffset)/3600,
				taskID,
			),
		)
	}

	return sb.String()
}

type ParamsRun struct {
	TimeInterval

	ID RunID // ID = 0 reserved for Maintenance.
}

func (params *ParamsRun) IsValidDuration() bool {
	return params.TimeStart >= params.TimeEnd
}

func (params *ParamsRun) IsValidID() bool {
	return params.ID > 0
}

func (res *ResourceScheduled) AddRun(_ context.Context, params *ParamsRun) error {
	if !params.IsValidDuration() {
		return goerrors.ErrInvalidInput{
			Caller:     "AddRun",
			InputName:  "TimeEnd",
			InputValue: params.TimeEnd,
			Issue: errors.New(
				"time start greater or equal to time end",
			),
		}
	}

	if !params.IsValidID() {
		return goerrors.ErrInvalidInput{
			Caller:     "AddRun",
			InputName:  "ID",
			InputValue: params.ID,
			Issue: goerrors.ErrNegativeInput{
				InputName: "ID",
			},
		}
	}

	res.mu.Lock()

	if _, alreadyScheduled := res.schedule[params.TimeInterval]; alreadyScheduled {
		res.mu.Unlock()

		return fmt.Errorf(
			"time interval %v already scheduled",
			params.TimeInterval,
		)
	}

	res.schedule[params.TimeInterval] = params.ID

	res.mu.Unlock()

	return nil
}

// removeRun should be called through Location which is mutex protected.
func (res *ResourceScheduled) removeRun(runID RunID) error {
	for interval, id := range res.schedule {
		if id == runID {
			delete(res.schedule, interval)

			return nil
		}
	}

	return fmt.Errorf(
		"run %d not found in schedule",
		runID,
	)
}

func (res *ResourceScheduled) IsAvailableIn(interval *TimeInterval) bool {
	res.mu.RLock()
	if _, exists := res.schedule[*interval]; exists {
		return false
	}
	res.mu.RUnlock()

	return true
}
