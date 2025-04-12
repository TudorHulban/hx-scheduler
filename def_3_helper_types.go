package hxscheduler

import (
	"fmt"
	"slices"
	"strings"
)

type ResourcesPerType map[ResourceType][]*ResourceScheduled

func (rpt ResourcesPerType) GetResourceTypesSorted() []ResourceType {
	result := make([]ResourceType, 0)

	for resourceType := range rpt {
		result = append(result, resourceType)
	}

	slices.Sort(result)

	return result
}

func (rpt ResourcesPerType) String() string {
	var sb strings.Builder
	sb.WriteString("ResourcesPerType{\n")

	// Sort keys for consistent output
	types := make([]ResourceType, 0, len(rpt))

	for t := range rpt {
		types = append(types, t)
	}

	slices.Sort(types)

	for _, t := range types {
		resources := rpt[t]
		sb.WriteString(fmt.Sprintf("\t%d: []*Resource{\n", t))

		for _, resource := range resources {
			if resource != nil {
				// Indent the resource string and add newlines after each line
				resourceStr := resource.String()
				resourceStr = strings.ReplaceAll(resourceStr, "\n", "\n\t\t")
				sb.WriteString("\t\t" + resourceStr + ",\n")
			} else {
				sb.WriteString("\t\tnil,\n")
			}
		}

		sb.WriteString("\t},\n")
	}

	sb.WriteString("}")

	return sb.String()
}
