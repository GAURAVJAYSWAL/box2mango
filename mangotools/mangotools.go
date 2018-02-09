package mangotools

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Domain struct {
	ID        int64  `json:"id"`
	DomainKey string `json:"domain_key"`
}
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type Folder struct {
	ID       int64          `json:"id"`
	Filename string         `json:"filename"`
	ParentID sql.NullString `json:"parent_id"`
}

func CreateUserBoxFolderEntry(email_id string, name string) error {
	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		return err
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id, domain_key FROM domains where domain_key='%v' limit 1", os.Getenv("DOMAINKEY"))
	// fmt.Println(sql)
	domainResults, err := db.Query(sql)
	if err != nil {
		return err
	}
	var domain Domain
	domainResults.Next()
	err = domainResults.Scan(&domain.ID, &domain.DomainKey)
	if err != nil {
		return err
	}
	// log.Printf("Got domain %v", domain)

	sql = fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE parent_id is null and is_folder=1 and is_deleted=0 and folder_type='A' and scope='AF' and filename='Network Drive' and domain_id=%v limit 1", domain.ID)
	// fmt.Println(sql)
	ndFolderResults, _ := db.Query(sql)
	var ndFolder Folder
	ndFolderResults.Next()
	err = ndFolderResults.Scan(&ndFolder.ID, &ndFolder.Filename, &ndFolder.ParentID)
	if err != nil {
		return err
	}
	// log.Printf("Got network drive %v", ndFolder)

	sql = fmt.Sprintf("SELECT id, name FROM users WHERE domain_id=%v and email_id='%v' limit 1", domain.ID, email_id)
	// fmt.Println(sql)
	userResults, _ := db.Query(sql)
	var user User

	userCnt := 0
	for userResults.Next() {
		userCnt++
		err = userResults.Scan(&user.ID, &user.Name)
		if err != nil {
			return err
		}
		log.Printf("Got user %v", user)
	}
	var sql1 string
	var sql2 string
	var uFolder Folder
	if userCnt == 0 {
		sql = fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE parent_id = %v and is_folder=1 and is_deleted=0 and filename='%v' and domain_id=%v limit 1", ndFolder.ID, os.Getenv("ORPHANFOLDERNAME"), domain.ID)
		// fmt.Println(sql)
		oFolderResults, _ := db.Query(sql)
		var oFolder Folder
		oFolderResults.Next()
		err = oFolderResults.Scan(&oFolder.ID, &oFolder.Filename, &oFolder.ParentID)
		if err != nil {
			return err
		}
		log.Printf("Got orphan folder %v", oFolder)

		sql = fmt.Sprintf("SELECT id, name FROM users WHERE domain_id=%v and email_id='%v' limit 1", domain.ID, os.Getenv("ORPHANFOLDEROWNER"))
		// fmt.Println(sql)
		userResults, _ := db.Query(sql)
		userResults.Next()
		err = userResults.Scan(&user.ID, &user.Name)
		if err != nil {
			return err
		}
		log.Printf("Got orphan owner %v", user)

		sql1 = fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE domain_id=%v and filename='%v' and parent_id=%v and is_deleted=0 limit 1", domain.ID, name, oFolder.ID)
		sql2 = fmt.Sprintf("INSERT INTO attachments (filename, name, kind, storage, storage_url, user_id, domain_id, access_type, created_at, updated_at, privacy_type, parent_id, is_folder, last_uploaded_by, folder_type, is_visible, follow_lists_count, modified_on) VALUES ('%v', '%v', 'FL', 'DB', 'http://', %v, %v, 'P', NOW(), NOW(), 'R', %v, 1, %v, 'U', 1, 1, NOW())", name, name, user.ID, domain.ID, oFolder.ID, user.ID)

	} else {
		sql1 = fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE domain_id=%v and filename='%v' and parent_id=%v and is_deleted=0 limit 1", domain.ID, user.Name, ndFolder.ID)
		sql2 = fmt.Sprintf("INSERT INTO attachments (filename, name, kind, storage, storage_url, user_id, domain_id, access_type, created_at, updated_at, privacy_type, parent_id, is_folder, last_uploaded_by, folder_type, is_visible, follow_lists_count, modified_on) VALUES ('%v', '%v', 'FL', 'DB', 'http://', %v, %v, 'P', NOW(), NOW(), 'R', %v, 1, %v, 'U', 1, 1, NOW())", user.Name, user.Name, user.ID, domain.ID, ndFolder.ID, user.ID)
	}
	// fmt.Println(sql)
	folderResults, _ := db.Query(sql1)
	folderCnt := 0
	for folderResults.Next() {
		folderCnt++
		err = folderResults.Scan(&uFolder.ID, &uFolder.Filename, &uFolder.ParentID)
		if err != nil {
			return err
		}

		log.Printf("Folder already exists. ID : %v", uFolder.ID)
	}

	if folderCnt == 0 {
		fmt.Printf("Creating folder : %v", user.Name)
		// fmt.Println(sql)
		fResults, _ := db.Exec(sql2)
		fID, _ := fResults.LastInsertId()
		uFolder.ID = fID
		uFolder.Filename = user.Name

		sql = fmt.Sprintf("INSERT INTO follow_list (attachment_id, user_id, created_at, updated_at, role_id) VALUES (%v, %v, NOW(), NOW(), 5)", uFolder.ID, user.ID)
		// fmt.Println(sql)
		followListResults, _ := db.Exec(sql)
		followListID, _ := followListResults.LastInsertId()
		log.Printf("Created folder. ID : %v", uFolder.ID)
		log.Printf("Created followlist. ID : %v", followListID)
	}

	return nil
}
