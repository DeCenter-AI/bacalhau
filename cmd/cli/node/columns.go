package node

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"github.com/bacalhau-project/bacalhau/cmd/util/output"
	"github.com/bacalhau-project/bacalhau/pkg/models"
	"github.com/bacalhau-project/bacalhau/pkg/util/idgen"
)

var alwaysColumns = []output.TableColumn[*models.NodeState]{
	{
		ColumnConfig: table.ColumnConfig{
			Name:             "id",
			WidthMax:         idgen.ShortIDLengthWithPrefix,
			WidthMaxEnforcer: func(col string, maxLen int) string { return idgen.ShortNodeID(col) }},
		Value: func(node *models.NodeState) string { return node.Info.ID() },
	},
	{
		ColumnConfig: table.ColumnConfig{Name: "type"},
		Value:        func(ni *models.NodeState) string { return ni.Info.NodeType.String() },
	},
	{
		ColumnConfig: table.ColumnConfig{Name: "approval"},
		Value:        func(ni *models.NodeState) string { return ni.Membership.String() },
	},
	{
		ColumnConfig: table.ColumnConfig{Name: "status"},
		Value: func(ni *models.NodeState) string {
			return ni.ConnectionState.Status.String()
		},
	},
}

var toggleColumns = map[string][]output.TableColumn[*models.NodeState]{
	"labels": {
		{
			ColumnConfig: table.ColumnConfig{Name: "labels", WidthMax: 50, WidthMaxEnforcer: text.WrapSoft},
			Value: func(ni *models.NodeState) string {
				labels := lo.MapToSlice(ni.Info.Labels, func(key, val string) string { return fmt.Sprintf("%s=%s", key, val) })
				slices.Sort(labels)
				return strings.Join(labels, " ")
			},
		},
	},
	"version": {
		{
			ColumnConfig: table.ColumnConfig{Name: "version"},
			Value: func(ni *models.NodeState) string {
				return ni.Info.BacalhauVersion.GitVersion
			},
		},
		{
			ColumnConfig: table.ColumnConfig{Name: "architecture"},
			Value: func(ni *models.NodeState) string {
				return ni.Info.BacalhauVersion.GOARCH
			},
		},
		{
			ColumnConfig: table.ColumnConfig{Name: "os"},
			Value: func(ni *models.NodeState) string {
				return ni.Info.BacalhauVersion.GOOS
			},
		},
	},
	"features": {
		{
			ColumnConfig: table.ColumnConfig{Name: "engines", WidthMax: maxLen(models.EngineNames), WidthMaxEnforcer: text.WrapSoft},
			Value: ifComputeNode(func(cni models.ComputeNodeInfo) string {
				return strings.Join(cni.ExecutionEngines, " ")
			}),
		},
		{
			ColumnConfig: table.ColumnConfig{Name: "inputs from", WidthMax: maxLen(models.StoragesNames), WidthMaxEnforcer: text.WrapSoft},
			Value: ifComputeNode(func(cni models.ComputeNodeInfo) string {
				return strings.Join(cni.StorageSources, " ")
			}),
		},
		{
			ColumnConfig: table.ColumnConfig{Name: "outputs", WidthMax: maxLen(models.PublisherNames), WidthMaxEnforcer: text.WrapSoft},
			Value: ifComputeNode(func(cni models.ComputeNodeInfo) string {
				return strings.Join(cni.Publishers, " ")
			}),
		},
	},
	"capacity": {
		{
			ColumnConfig: table.ColumnConfig{Name: "cpu", WidthMax: len("1.0 / "), WidthMaxEnforcer: text.WrapSoft},
			Value: ifComputeNode(func(cni models.ComputeNodeInfo) string {
				return fmt.Sprintf("%.1f / %.1f", cni.AvailableCapacity.CPU, cni.MaxCapacity.CPU)
			}),
		},
		{
			ColumnConfig: table.ColumnConfig{Name: "memory", WidthMax: len("10.0 GB / "), WidthMaxEnforcer: text.WrapSoft},
			Value: ifComputeNode(func(cni models.ComputeNodeInfo) string {
				return fmt.Sprintf("%s / %s",
					humanize.Bytes(cni.AvailableCapacity.Memory),
					humanize.Bytes(cni.MaxCapacity.Memory))
			}),
		},
		{
			ColumnConfig: table.ColumnConfig{Name: "disk", WidthMax: len("100.0 GB / "), WidthMaxEnforcer: text.WrapSoft},
			Value: ifComputeNode(func(cni models.ComputeNodeInfo) string {
				return fmt.Sprintf("%s / %s",
					humanize.Bytes(cni.AvailableCapacity.Disk),
					humanize.Bytes(cni.MaxCapacity.Disk))
			}),
		},
		{
			ColumnConfig: table.ColumnConfig{Name: "gpu", WidthMax: len("1 / "), WidthMaxEnforcer: text.WrapSoft},
			Value: ifComputeNode(func(cni models.ComputeNodeInfo) string {
				return fmt.Sprintf("%d / %d", cni.AvailableCapacity.GPU, cni.MaxCapacity.GPU)
			}),
		},
	},
}

func maxLen(val []string) int {
	return lo.Max(lo.Map[string, int](val, func(item string, index int) int { return len(item) })) + 1
}

func ifComputeNode(getFromCNInfo func(models.ComputeNodeInfo) string) func(state *models.NodeState) string {
	return func(ni *models.NodeState) string {
		if !ni.Info.IsComputeNode() {
			return ""
		}
		return getFromCNInfo(ni.Info.ComputeNodeInfo)
	}
}
