package mongo

const (
	DB_NAME             = "goauth"
	USER_COLLECTION     = "users"
	AUDITLOG_COLLECTION = "auditlog"
)

// import (
// 	"github.com/calvine/goauth/models/nullable"
// )

// type NullableTimeBSON nullable.NullableTime

// func (ntb *NullableTimeBSON) MarshalBSON(date []byte) error {
// 	ntb.IsNull = true
// 	return nil
// }
