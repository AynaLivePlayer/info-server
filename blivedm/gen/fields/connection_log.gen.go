package fields

import "github.com/rhine-tech/scene/model/filter"

var ConnectionLog = &struct {
	RoomID filter.Field
	Source filter.Field
	Time   filter.Field
}{
	RoomID: "ConnectionLog.RoomID",
	Source: "ConnectionLog.Source",
	Time:   "ConnectionLog.Time",
}
