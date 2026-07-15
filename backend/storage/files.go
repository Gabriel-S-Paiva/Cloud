package storage

import "database/sql"

type File struct {
	Id           int
	DisplayName  string
	OwnedBy      int
	Size         int
	UploadedAt   int
	LastModified int
	ParentFolder sql.NullInt64
}
