package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// 连接到数据库
	connStr := "host=localhost user=oidc_user password=oidc_password dbname=oidc_db port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// 查询oauth_refresh_tokens表中的数据
	rows, err := db.Query("SELECT id, token_hash, user_id, client_id, scopes, expires_at, revoked_at FROM oauth_refresh_tokens")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Refresh Tokens in database:")
	fmt.Println("ID\tToken Hash\t\t\t\tUser ID\tClient ID\tScopes\t\tExpires At\t\t\tRevoked At")
	fmt.Println("------------------------------------------------------------------------------------------------------------------------")

	for rows.Next() {
		var id int64
		var tokenHash string
		var userID int64
		var clientID string
		var scopes []string
		var expiresAt string
		var revokedAt sql.NullString

		err := rows.Scan(&id, &tokenHash, &userID, &clientID, &scopes, &expiresAt, &revokedAt)
		if err != nil {
			log.Fatal(err)
		}

		revokedStr := "NULL"
		if revokedAt.Valid {
			revokedStr = revokedAt.String
		}

		fmt.Printf("%d\t%s\t%d\t%s\t%v\t%s\t%s\n", id, tokenHash[:20], userID, clientID, scopes, expiresAt, revokedStr)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}