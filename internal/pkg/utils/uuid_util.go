package utils

import (
	"fmt"
	"github.com/google/uuid"
)

// ToMSSQLUUIDString converts a standard uuid.UUID into the mixed-endian string 
// format that SQL Server uses for uniqueidentifier columns. 
// This is required when querying by ID because GORM/google-uuid uses standard big-endian format.
func ToMSSQLUUIDString(u uuid.UUID) string {
	return fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
		u[3], u[2], u[1], u[0],
		u[5], u[4],
		u[7], u[6],
		u[8], u[9],
		u[10], u[11], u[12], u[13], u[14], u[15])
}
