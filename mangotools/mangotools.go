package mangotools

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/siddhartham/box2mango/lib"
)

type MangoService struct {
}

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
	UserID   int64          `json:"user_id"`
	DomainID int64          `json:"domain_id"`
}
type File struct {
	ID         int64          `json:"id"`
	Filename   string         `json:"filename"`
	ParentID   sql.NullString `json:"parent_id"`
	UserID     int64          `json:"user_id"`
	DomainID   int64          `json:"domain_id"`
	AccessType string         `json:"access_type"`
	Storage    string         `json:"storage"`
	StorageURL string         `json:"storage_url"`
	IsVisible  string         `json:"is_visible"`
}

func (ma *MangoService) CreateFollowListEntry(folderID int64, emailID string, roleID int) (int64, error) {
	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		lib.Err("CreateFollowListEntry", err)
		return 0, err
	}
	defer db.Close()

	var user User
	userCnt := 0
	sql1 := fmt.Sprintf("SELECT id, name FROM users WHERE domain_id in (SELECT id FROM domains where domain_key='%v') and email_id='%v' limit 1", os.Getenv("DOMAINKEY"), emailID)
	userResults, _ := db.Query(sql1)
	for userResults.Next() {
		userCnt++
		userResults.Scan(&user.ID, &user.Name)
	}
	if userCnt != 0 {
		sql2 := fmt.Sprintf("INSERT INTO follow_list (attachment_id, user_id, created_at, updated_at, role_id) VALUES (%v, %v, NOW(), NOW(), %v)", folderID, user.ID, roleID)
		followResults, err1 := db.Exec(sql2)
		if err1 == nil {
			followID, _ := followResults.LastInsertId()
			lib.Info(fmt.Sprintf("Created followlist ID : %v", followID))
			return followID, nil
		}
	}
	return 0, nil
}

func (ma *MangoService) CreateBoxChildFileEntry(folderID int64, fileName string, sanPath string, externalID string) (int64, error) {
	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		lib.Err("CreateBoxChildFolderEntry", err)
		return 0, err
	}
	defer db.Close()

	var pFolder Folder
	sql := fmt.Sprintf("SELECT id, filename, parent_id, user_id, domain_id FROM attachments WHERE id = %v and is_folder=1 and is_deleted=0 limit 1", folderID)
	parentFolderResults, _ := db.Query(sql)
	parentFolderResults.Next()
	parentFolderResults.Scan(&pFolder.ID, &pFolder.Filename, &pFolder.ParentID, &pFolder.UserID, &pFolder.DomainID)

	var uFile File
	sql1 := fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE filename='%v' and parent_id=%v and external_id='%v' and is_deleted=0 limit 1", fileName, pFolder.ID, externalID)
	fileResults, _ := db.Query(sql1)
	fileCnt := 0
	for fileResults.Next() {
		fileCnt++
		fileResults.Scan(&uFile.ID, &uFile.Filename, &uFile.ParentID)
		lib.Info(fmt.Sprintf("File already exists ID : %v", uFile.ID))
		return uFile.ID, nil
	}
	if fileCnt == 0 {
		sql2 := fmt.Sprintf("INSERT INTO attachments (filename, name, kind, storage, storage_url, user_id, domain_id, access_type, created_at, updated_at, privacy_type, parent_id, is_folder, last_uploaded_by, folder_type, is_visible, follow_lists_count, modified_on, external_id) VALUES ('%v', '%v', 'FA', 'SAN', '%v', %v, %v, 'P', NOW(), NOW(), 'R', %v, 0, %v, 'U', 1, 1, NOW(), '%v')", fileName, fileName, sanPath, pFolder.UserID, pFolder.DomainID, pFolder.ID, pFolder.UserID, externalID)
		fResults, _ := db.Exec(sql2)
		fID, _ := fResults.LastInsertId()

		sql3 := fmt.Sprintf("INSERT INTO follow_list (attachment_id, user_id, created_at, updated_at, role_id) VALUES (%v, %v, NOW(), NOW(), 1)", fID, pFolder.UserID)
		followResults, _ := db.Exec(sql3)
		followID, _ := followResults.LastInsertId()
		lib.Info(fmt.Sprintf("Created file. ID : %v", fID))
		lib.Info(fmt.Sprintf("Created followlist. ID : %v", followID))
		return fID, nil
	}

	return 0, nil
}
func (ma *MangoService) CreateBoxChildFolderEntry(userExternalID string, parentFolderExternalID string, folderName string, externalID string) (int64, error) {
	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		lib.Err("CreateBoxChildFolderEntry", err)
		return 0, err
	}
	defer db.Close()

	var sql string
	var pFolder Folder
	if parentFolderExternalID == "0" {
		sql = fmt.Sprintf("SELECT id, filename, parent_id, user_id, domain_id FROM attachments WHERE external_id = '%v' and is_folder=1 and is_deleted=0 limit 1", userExternalID)
	} else {
		sql = fmt.Sprintf("SELECT id, filename, parent_id, user_id, domain_id FROM attachments WHERE external_id = '%v' and is_folder=1 and is_deleted=0 limit 1", parentFolderExternalID)
	}
	parentFolderResults, _ := db.Query(sql)
	parentFolderResults.Next()
	parentFolderResults.Scan(&pFolder.ID, &pFolder.Filename, &pFolder.ParentID, &pFolder.UserID, &pFolder.DomainID)

	var uFolder Folder
	sql1 := fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE filename='%v' and parent_id=%v and external_id='%v' and is_deleted=0 limit 1", folderName, pFolder.ID, externalID)
	folderResults, _ := db.Query(sql1)
	folderCnt := 0

	for folderResults.Next() {
		folderCnt++
		folderResults.Scan(&uFolder.ID, &uFolder.Filename, &uFolder.ParentID)
		lib.Info(fmt.Sprintf("Folder already exists ID : %v", uFolder.ID))
		return uFolder.ID, nil
	}
	if folderCnt == 0 {
		sql2 := fmt.Sprintf("INSERT INTO attachments (filename, name, kind, storage, storage_url, user_id, domain_id, access_type, created_at, updated_at, privacy_type, parent_id, is_folder, last_uploaded_by, folder_type, is_visible, follow_lists_count, modified_on, external_id) VALUES ('%v', '%v', 'FL', 'DB', 'http://', %v, %v, 'P', NOW(), NOW(), 'R', %v, 1, %v, 'U', 1, 1, NOW(), '%v')", folderName, folderName, pFolder.UserID, pFolder.DomainID, pFolder.ID, pFolder.UserID, externalID)
		fResults, _ := db.Exec(sql2)
		fID, _ := fResults.LastInsertId()

		sql = fmt.Sprintf("INSERT INTO follow_list (attachment_id, user_id, created_at, updated_at, role_id) VALUES (%v, %v, NOW(), NOW(), 5)", fID, pFolder.UserID)
		followResults, _ := db.Exec(sql)
		followID, _ := followResults.LastInsertId()
		lib.Info(fmt.Sprintf("Created folder ID : %v", fID))
		lib.Info(fmt.Sprintf("Created followlist ID : %v", followID))
		return fID, nil
	}

	return 0, nil
}

func (ma *MangoService) CreateUserBoxFolderEntry(emailID string, name string, externalID string) (int64, error) {
	var newFolderName string

	db, err := sql.Open("mysql", os.Getenv("MYSQL"))
	if err != nil {
		lib.Err("CreateUserBoxFolderEntry", err)
		return 0, err
	}
	defer db.Close()

	var domain Domain
	sql := fmt.Sprintf("SELECT id, domain_key FROM domains where domain_key='%v' limit 1", os.Getenv("DOMAINKEY"))
	domainResults, _ := db.Query(sql)
	domainResults.Next()
	domainResults.Scan(&domain.ID, &domain.DomainKey)

	var ndFolder Folder
	sql = fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE parent_id is null and is_folder=1 and is_deleted=0 and folder_type='A' and scope='AF' and filename='Network Drive' and domain_id=%v limit 1", domain.ID)
	ndFolderResults, _ := db.Query(sql)
	ndFolderResults.Next()
	ndFolderResults.Scan(&ndFolder.ID, &ndFolder.Filename, &ndFolder.ParentID)

	var user User
	userCnt := 0
	sql = fmt.Sprintf("SELECT id, name FROM users WHERE domain_id=%v and email_id='%v' limit 1", domain.ID, emailID)
	userResults, _ := db.Query(sql)
	for userResults.Next() {
		userCnt++
		userResults.Scan(&user.ID, &user.Name)
	}

	var sql1 string
	var sql2 string
	var uFolder Folder
	if userCnt == 0 {
		sql = fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE parent_id = %v and is_folder=1 and is_deleted=0 and filename='%v' and domain_id=%v limit 1", ndFolder.ID, os.Getenv("ORPHANFOLDERNAME"), domain.ID)
		oFolderResults, _ := db.Query(sql)
		var oFolder Folder
		oFolderResults.Next()
		oFolderResults.Scan(&oFolder.ID, &oFolder.Filename, &oFolder.ParentID)

		sql = fmt.Sprintf("SELECT id, name FROM users WHERE domain_id=%v and email_id='%v' limit 1", domain.ID, os.Getenv("ORPHANFOLDEROWNER"))
		userResults, _ := db.Query(sql)
		userResults.Next()
		userResults.Scan(&user.ID, &user.Name)

		newFolderName = fmt.Sprintf("%v%v", name, os.Getenv("BOXFOLDERSUFFIX"))
		sql1 = fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE domain_id=%v and filename='%v' and parent_id=%v and external_id='%v' and is_deleted=0 limit 1", domain.ID, newFolderName, oFolder.ID, externalID)
		sql2 = fmt.Sprintf("INSERT INTO attachments (filename, name, kind, storage, storage_url, user_id, domain_id, access_type, created_at, updated_at, privacy_type, parent_id, is_folder, last_uploaded_by, folder_type, is_visible, follow_lists_count, modified_on, external_id) VALUES ('%v', '%v', 'FL', 'DB', 'http://', %v, %v, 'P', NOW(), NOW(), 'R', %v, 1, %v, 'U', 1, 1, NOW(), '%v')", newFolderName, newFolderName, user.ID, domain.ID, oFolder.ID, user.ID, externalID)
	} else {
		newFolderName = fmt.Sprintf("%v%v", user.Name, os.Getenv("BOXFOLDERSUFFIX"))
		sql1 = fmt.Sprintf("SELECT id, filename, parent_id FROM attachments WHERE domain_id=%v and filename='%v' and parent_id=%v and external_id='%v' and is_deleted=0 limit 1", domain.ID, newFolderName, ndFolder.ID, externalID)
		sql2 = fmt.Sprintf("INSERT INTO attachments (filename, name, kind, storage, storage_url, user_id, domain_id, access_type, created_at, updated_at, privacy_type, parent_id, is_folder, last_uploaded_by, folder_type, is_visible, follow_lists_count, modified_on, external_id) VALUES ('%v', '%v', 'FL', 'DB', 'http://', %v, %v, 'P', NOW(), NOW(), 'R', %v, 1, %v, 'U', 1, 1, NOW(), '%v')", newFolderName, newFolderName, user.ID, domain.ID, ndFolder.ID, user.ID, externalID)
	}
	folderResults, _ := db.Query(sql1)
	folderCnt := 0
	for folderResults.Next() {
		folderCnt++
		folderResults.Scan(&uFolder.ID, &uFolder.Filename, &uFolder.ParentID)
		lib.Info(fmt.Sprintf("Folder already exists ID : %v", uFolder.ID))
		return uFolder.ID, nil
	}
	if folderCnt == 0 {
		fResults, _ := db.Exec(sql2)
		fID, _ := fResults.LastInsertId()
		uFolder.ID = fID
		uFolder.Filename = user.Name

		sql = fmt.Sprintf("INSERT INTO follow_list (attachment_id, user_id, created_at, updated_at, role_id) VALUES (%v, %v, NOW(), NOW(), 5)", uFolder.ID, user.ID)
		followListResults, _ := db.Exec(sql)
		followListID, _ := followListResults.LastInsertId()
		lib.Info(fmt.Sprintf("Created folder ID : %v", uFolder.ID))
		lib.Info(fmt.Sprintf("Created followlist ID : %v", followListID))
		return uFolder.ID, nil
	}

	return 0, nil
}
