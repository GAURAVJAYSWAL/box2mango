package mango

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Domain struct {
	ID        int    `json:"id"`
	DomainKey string `json:"domain_key"`
}
type Folder struct {
	ID       int            `json:"id"`
	Filename string         `json:"filename"`
	ParentID sql.NullString `json:"parent_id"`
}

func createFolderEntry(name string) error {
	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		return err
	}
	defer db.Close()

	domainResults, err := db.Query(fmt.Sprintf("SELECT id, domain_key FROM domains where domain_key='%v' limit 1", os.Getenv("DOMAINKEY")))
	if err != nil {
		return err
	}
	var domain Domain
	domainResults.Next()
	err = domainResults.Scan(&domain.ID, &domain.DomainKey)
	if err != nil {
		return err
	}
	log.Printf("Got domain %v", domain)

	ndFolderResults, err := db.Query(fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE parent_id is null and is_folder=1 and is_deleted=0 and folder_type='A' and scope='AF' and filename='Network Drive' and domain_id=%v limit 1", domain.ID))
	var ndFolder Folder
	ndFolderResults.Next()
	err = ndFolderResults.Scan(&ndFolder.ID, &ndFolder.Filename, &ndFolder.ParentID)
	if err != nil {
		return err
	}
	log.Printf("Got network drive %v", ndFolder)

	return nil
}
