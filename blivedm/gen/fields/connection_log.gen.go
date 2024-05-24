package fields

import (
	"github.com/rhine-tech/scene/model/query"
)

var ConnectionLog = &struct {
	RoomID query.Field
	Source query.Field
	Time   query.Field
}{
	RoomID: "ConnectionLog.RoomID",
	Source: "ConnectionLog.Source",
	Time:   "ConnectionLog.Time",
}
