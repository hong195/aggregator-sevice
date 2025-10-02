package repo

import "time"

type DataPacketCriteria struct {
	Start time.Time
	End   time.Time
}

func MakeDataPacketCriteria(start, end time.Time) *DataPacketCriteria {
	return &DataPacketCriteria{
		Start: start,
		End:   end,
	}
}
