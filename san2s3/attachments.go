package san2s3

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Domain struct {
	ID        int    `json:"id"`
	DomainKey string `json:"domain_key"`
}
type Folder struct {
	ID            int            `json:"id"`
	Filename      string         `json:"filename"`
	ParentID      sql.NullString `json:"parent_id"`
	SanStorageUrl string         `json:"san_storage_url"`
}
