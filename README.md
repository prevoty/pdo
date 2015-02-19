#pdo
Prevoty data objects

A simple object/table mapper for easy access to your table records. Fetch data using raw sql condition clauses for complex record fetching or simple column matching.  Uses tags to denote columns from regular struct members. 

Utilizes go's reflection library for type flexibility, and prompts incorrect types with intelligible panics.


### currently supports mysql

```go
package main

import (
	"log"

	"github.com/prevoty/pdo"
)

type User struct {
	_meta string `table:"user"`

	Id    int    `column:"id"`
	First string `column:"first_name"`
	Last  string `column:"last_name"`
}

func main() {

	db, err = pdo.NewMySQL("user@tcp(localhost:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}

	id, err := db.Create(&User{
		First: "Artemis",
		Last:  "Prime",
	})
	if err != nil {
		log.Fatal(err)
	}

	user := new(User)

	err = db.Find(user, "WHERE `id` = ?", id)
	switch err {
	case sql.ErrNoRows:
		log.Println("no record found...")
	case nil:
		log.Println("found a user")
	default:
		log.Fatal(err)
	}

	user.Last = "Johnson"

	err = db.Update(user)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Delete(user)
	if err != nil {
		log.Fatal(err)
	}

	// for fetching row sets
	users := make([]*User, 0, 0)

	err = db.FindAll(&users, `
		AS u
		INNER JOIN organization o ON u.oid = o.id
		WHERE o.id = ?
		ORDER BY u.last_name
		LIMIT 10
	`, 123)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		log.Printf("%s is what's what", u.First)
	}

	// for direct access to *sql.DB, use the DB member
	// ...result, err := d.DB.QueryRow("select * from...", params...)

}
```
