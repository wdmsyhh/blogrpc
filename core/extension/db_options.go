package extension

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/qiniu/qmgo"
)

const (
	DEFAULT_CURSOR_BATCH_SIZE = 100
)

type IterateOption struct {
	Fields    primitive.M
	BatchSize int64
	Sortor    []string
}

func (i *IterateOption) AppendOptionToQueryI(q qmgo.QueryI) qmgo.QueryI {
	newQueryI := q
	// append sortor
	if i.Sortor != nil && len(i.Sortor) > 0 {
		newQueryI = q.Sort(i.Sortor...)
	}
	// append fields
	if i.Fields != nil {
		newQueryI = newQueryI.Select(i.Fields)
	}
	// append batch size
	newQueryI = q.BatchSize(DEFAULT_CURSOR_BATCH_SIZE)
	if i.BatchSize != 0 {
		newQueryI = newQueryI.BatchSize(i.BatchSize)
	}
	return newQueryI
}
