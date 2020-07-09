package boilingcore

import (
	"sort"

	"github.com/volatiletech/sqlboiler/v4/drivers"
)

// Orders defines orders for the generation run
type Orders struct {
	Columns map[string]int `toml:"columns,omitempty" json:"columns,omitempty"`
}

func (o Orders) sortColumns(cols []drivers.Column) []drivers.Column {
	if o.Columns == nil {
		return cols
	}
	type indexOrder struct {
		index int
		order int
	}
	var indexOrders []indexOrder
	for i, col := range cols {
		if score, ok := o.Columns[col.Name]; ok {
			indexOrders = append(indexOrders, indexOrder{i, score})
			continue
		}
		// unspecified columns to 0
		indexOrders = append(indexOrders, indexOrder{i, 0})
	}
	sort.Slice(indexOrders, func(i, j int) bool {
		return indexOrders[i].order < indexOrders[j].order
	})

	sortedCols := make([]drivers.Column, len(cols))
	for i, indexOrder := range indexOrders {
		sortedCols[i] = cols[indexOrder.index]
	}
	return sortedCols
}
