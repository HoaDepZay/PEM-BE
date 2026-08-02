# Backend AI Coding Rules for Visual Finance (Golang + MSSQL)

## 1. SQL Server UNIQUEIDENTIFIER (UUID) Endianness
**Problem**: Microsoft SQL Server stores the first 3 components of a UNIQUEIDENTIFIER in little-endian format. The standard github.com/google/uuid package assumes big-endian (RFC 4122 standard). If you use uuid.UUID directly in GORM models with SQL Server, reading from the DB will flip the bytes, resulting in a corrupted UUID string (e.g., e5a7... becomes 305a...). This completely breaks Foreign Key constraints when inserting related records.
**Rule**: NEVER use uuid.UUID directly in GORM models for SQL Server. ALWAYS use the custom wrapper models.MSSQLUUID which explicitly implements the Scan and Value interfaces to reverse the bytes back to standard network byte order.

## 2. Go Type Aliasing and JSON Marshalling
**Problem**: When creating a custom type based on an existing type (e.g., 	ype MSSQLUUID uuid.UUID), Go DOES NOT inherit the methods of the underlying type.
**Rule**: Because MSSQLUUID is fundamentally a [16]byte array, if you do not explicitly implement MarshalJSON and UnmarshalJSON, Go's encoding/json will serialize it as a raw array of 16 integers (e.g., [19, 4, 73, ...]) instead of a UUID string. Whenever you create a custom type alias for an API response field, YOU MUST explicitly implement the json.Marshaler and json.Unmarshaler interfaces.
