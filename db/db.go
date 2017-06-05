package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ingmardrewing/gomicSocMed/config"
)

var db *sql.DB

func Initialize() {
	dsn := config.GetDsn()
	db, _ = sql.Open("mysql", dsn)
}

func AddEmailAddress(email string, token string, deletion_token string) {
	stmt, err := db.Prepare("INSERT INTO newsletter_addresses(email, double_opt_mail_token, deletion_token ) VALUES(?, ?, ?)")
	handleErr(err)
	_, err = stmt.Exec(email, token, deletion_token)
	handleErr(err)
}

func VerifySubscription(token string) {
	stmt, err := db.Prepare("UPDATE newsletter_addresses SET double_opt_verified=1 WHERE double_opt_verified != 1 AND double_opt_mail_token=?")
	handleErr(err)
	_, err = stmt.Exec(token)
	handleErr(err)
}

func TokenExists(token string) bool {
	var amount string
	err := db.QueryRow("SELECT count(*) FROM newsletter_addresses WHERE double_opt_verified != 1 AND double_opt_mail_token=?", token).Scan(&amount)
	handleErr(err)
	return amount != "0"
}

func DeletionTokenExists(token string) bool {
	var amount string
	err := db.QueryRow("SELECT count(*) FROM newsletter_addresses WHERE  deletion_token=?", token).Scan(&amount)
	handleErr(err)
	return amount != "0"
}

func DeleteEmailAddressWithToken(deletion_token string) {
	stmt, err := db.Prepare("DELETE FROM newsletter_addresses WHERE deletion_token=?")
	handleErr(err)
	_, err = stmt.Exec(deletion_token)
	handleErr(err)
}

func AddressExists(email string) bool {
	var amount string
	err := db.QueryRow("SELECT count(*) FROM newsletter_addresses WHERE email=?", email).Scan(&amount)
	handleErr(err)
	return amount != "0"
}

func GetNewsletterRecipients() []string {
	rows, err := db.Query("SELECT email FROM newsletter_addresses WHERE double_opt_verified = 1")
	handleErr(err)
	defer rows.Close()
	var (
		email  string
		emails []string
	)

	for rows.Next() {
		err := rows.Scan(&email)
		handleErr(err)
		emails = append(emails, email)
	}
	err = rows.Err()
	handleErr(err)
	return emails
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
