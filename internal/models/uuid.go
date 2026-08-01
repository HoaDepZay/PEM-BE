package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
)

// MSSQLUUID is a wrapper around uuid.UUID that handles SQL Server's little-endian byte order
// for the first three components of a UNIQUEIDENTIFIER.
type MSSQLUUID uuid.UUID

func (u *MSSQLUUID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, okStr := value.(string)
		if okStr {
			parsed, err := uuid.Parse(str)
			if err == nil {
				*u = MSSQLUUID(parsed)
				return nil
			}
		}
		return fmt.Errorf("could not scan type %T into MSSQLUUID", value)
	}
	if len(bytes) == 16 {
		var newUUID uuid.UUID
		copy(newUUID[:], bytes)
		// Swap bytes back to standard
		newUUID[0], newUUID[3] = newUUID[3], newUUID[0]
		newUUID[1], newUUID[2] = newUUID[2], newUUID[1]
		newUUID[4], newUUID[5] = newUUID[5], newUUID[4]
		newUUID[6], newUUID[7] = newUUID[7], newUUID[6]
		*u = MSSQLUUID(newUUID)
		return nil
	}
	return fmt.Errorf("invalid length for MSSQLUUID")
}

func (u MSSQLUUID) Value() (driver.Value, error) {
	if uuid.UUID(u) == uuid.Nil {
		return nil, nil
	}
	return uuid.UUID(u).String(), nil
}

func (u MSSQLUUID) String() string {
	return uuid.UUID(u).String()
}

func NewMSSQLUUID() MSSQLUUID {
    return MSSQLUUID(uuid.New())
}

var NilMSSQLUUID = MSSQLUUID(uuid.Nil)

func (u MSSQLUUID) MarshalJSON() ([]byte, error) {
	return []byte(" + uuid.UUID(u).String() + "), nil
}

func (u *MSSQLUUID) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == "" {
		return nil
	}
	if len(str) < 2 || str[0] != '"' || str[len(str)-1] != '"' {
		return fmt.Errorf("invalid UUID format")
	}
	parsed, err := uuid.Parse(str[1 : len(str)-1])
	if err != nil {
		return err
	}
	*u = MSSQLUUID(parsed)
	return nil
}
