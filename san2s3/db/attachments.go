package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type File struct {
	ID            int            `json:"id"`
	Filename      string         `json:"filename"`
	ParentID      sql.NullString `json:"parent_id"`
	UserID        int            `json:"user_id"`
	DomainID      int            `json:"domain_id"`
	SanStorageUrl string         `json:"san_storage_url"`
	DomainKey     string         `json:"domain_key"`
}

func getFileFromDb() ([]File, error) {
	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var sql string = "select san_storage_url,a.id,a.user_id,domain_id,parent_id,filename,domain_key from attachments a inner join domains d on d.id=a.domain_id where storage='SAN' and is_deleted=0 and is_visible=0 and workflow_status is null and is_folder=0 limit 5"
	parentFolderResults, err := db.Query(sql)
	var pFileList []File
	for parentFolderResults.Next() {
		var pFile File
		err := parentFolderResults.Scan(&pFile.SanStorageUrl, &pFile.ID, &pFile.UserID, &pFile.DomainID, &pFile.ParentID, &pFile.Filename, &pFile.DomainKey)
		if err != nil {
			fmt.Printf("mysql: could not read row: %v", err)
			return nil, err
		} else {
			pFileList = append(pFileList, pFile)
		}
	}
	fmt.Println(pFileList)
	return pFileList, nil
}
