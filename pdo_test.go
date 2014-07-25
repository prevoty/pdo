package pdo

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type (
	TestUser struct {
		_meta string `table:"test_table"`

		Id    int64          `column:"id"`
		First string         `column:"first_name"`
		Last  sql.NullString `column:"last_name"`
	}
)

var (
	mysqldb  *MySQL
	sqlitedb *Sqlite

	MySQLbuilt  = false
	Sqlitebuilt = false
)

func TestSqliteFind(t *testing.T) {

	EnsureSqlite()

	user := new(TestUser)
	err := sqlitedb.Find(user, "WHERE `first_name` like ?", "user1")
	if err != nil {
		t.Error(err)
	}

}

func TestSqliteFindAll(t *testing.T) {

	EnsureSqlite()

	rows := make([]*TestUser, 0, 0)

	err := sqlitedb.FindAll(&rows, `
		AS tu
		WHERE tu.first_name != ""
	`)
	if err != nil {
		t.Error(err)
	}

	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got: %d\n", len(rows))
	}

	for i, v := range rows {
		if v.Id != int64(i+1) {
			t.Errorf("expected: %d, got: %d\n", i, v.Id)
		}
		if v.Last.Valid == false || v.Last.String != "johnson" {
			t.Errorf("expected: johnson, got: %s\n", v.Last.String)
		}
	}

}

func TestSqliteCreate(t *testing.T) {

	EnsureSqlite()

	id, err := sqlitedb.Create(&TestUser{
		First: "user5",
		Last:  sql.NullString{String: "johnson", Valid: true},
	})
	if err != nil {
		t.Error(err)
	}

	user := new(TestUser)
	err = sqlitedb.Find(user, "WHERE id = ?", id)
	if err != nil {
		t.Error(err)
	}

}

func TestSqliteDelete(t *testing.T) {

	EnsureSqlite()

	user := new(TestUser)
	err := sqlitedb.Find(user, "WHERE `last_name` = ?", "johnson")
	if err != nil {
		t.Error(err)
	}

	err = sqlitedb.Delete(user)
	if err != nil {
		t.Error(err)
	}

}

func TestSqliteUpdate(t *testing.T) {

	EnsureSqlite()

	// create a user
	id, err := sqlitedb.Create(&TestUser{
		First: "user6",
		Last:  sql.NullString{String: "johnson", Valid: true},
	})
	if err != nil {
		t.Error(err)
	}

	// load the user record
	user := new(TestUser)
	err = sqlitedb.Find(user, "WHERE id = ?", id)
	if err != nil {
		t.Error(err)
	}

	// change a field
	user.Last.String = "jackson"
	user.Last.Valid = true

	// persist the field change
	err = sqlitedb.Update(user)
	if err != nil {
		t.Error(err)
	}

	// load the persisted record
	user = new(TestUser)
	err = sqlitedb.Find(user, "WHERE id = ?", id)
	if err != nil {
		t.Error(err)
	}

	if user.Last.Valid == false || user.Last.String != "jackson" {
		t.Errorf("expected: jackson, got: %s\n", user.Last.String)
	}

	err = sqlitedb.Delete(user)
	if err != nil {
		t.Error(err)
	}

}

func TestMySQLFind(t *testing.T) {

	EnsureMySQL()

	user := new(TestUser)
	err := mysqldb.Find(user, "WHERE `first_name` like ?", "user1")
	if err != nil {
		t.Error(err)
	}

}

func TestMySQLFindAll(t *testing.T) {

	EnsureMySQL()

	rows := make([]*TestUser, 0, 0)

	err := mysqldb.FindAll(&rows, "WHERE `first_name` != ''")
	if err != nil {
		t.Error(err)
	}

	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got: %d\n", len(rows))
	}

	for i, v := range rows {
		if v.Id != int64(i+1) {
			t.Errorf("expected: %d, got: %d\n", i, v.Id)
		}
		if v.Last.Valid == false || v.Last.String != "johnson" {
			t.Errorf("expected: johnson, got: %s\n", v.Last.String)
		}
	}

}

func TestMySQLCreate(t *testing.T) {

	EnsureMySQL()

	id, err := mysqldb.Create(&TestUser{
		First: "user5",
		Last:  sql.NullString{String: "johnson", Valid: true},
	})
	if err != nil {
		t.Error(err)
	}

	user := new(TestUser)
	err = mysqldb.Find(user, "WHERE id = ?", id)
	if err != nil {
		t.Error(err)
	}

}

func TestMySQLDelete(t *testing.T) {

	EnsureMySQL()

	user := new(TestUser)
	err := mysqldb.Find(user, "WHERE `last_name` = ?", "johnson")
	if err != nil {
		t.Error(err)
	}

	err = mysqldb.Delete(user)
	if err != nil {
		t.Error(err)
	}

}

func TestMySQLUpdate(t *testing.T) {

	EnsureMySQL()

	// create a user
	id, err := mysqldb.Create(&TestUser{
		First: "user6",
		Last:  sql.NullString{String: "johnson", Valid: true},
	})
	if err != nil {
		t.Error(err)
	}

	// load the user record
	user := new(TestUser)
	err = mysqldb.Find(user, "WHERE id = ?", id)
	if err != nil {
		t.Error(err)
	}

	// change a field
	user.Last.String = "jackson"
	user.Last.Valid = true

	// persist the field change
	err = mysqldb.Update(user)
	if err != nil {
		t.Error(err)
	}

	// load the persisted record
	user = new(TestUser)
	err = mysqldb.Find(user, "WHERE id = ?", id)
	if err != nil {
		t.Error(err)
	}

	if user.Last.Valid == false || user.Last.String != "jackson" {
		t.Errorf("expected: jackson, got: %s\n", user.Last.String)
	}

	err = mysqldb.Delete(user)
	if err != nil {
		t.Error(err)
	}

}

func EnsureMySQL() {

	if MySQLbuilt {
		return
	}

	database, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/pdo_test")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec("CREATE DATABASE IF NOT EXISTS `pdo_test`")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec("DROP TABLE IF EXISTS `test_table`")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec("CREATE TABLE `pdo_test`.`test_table` (`id` int NOT NULL AUTO_INCREMENT, `first_name` text, `last_name` text, PRIMARY KEY (`id`)) ENGINE=`InnoDB` DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci")
	if err != nil {
		log.Fatal(err)
	}

	database.Exec("INSERT INTO `test_table` (`first_name`,`last_name`) VALUES (\"user1\",\"johnson\"),(\"user2\",\"johnson\")")
	database.Close()
	// done building, continue with the tests

	mysqldb, err = NewMySQL("root:password@tcp(127.0.0.1:3306)/pdo_test")
	if err != nil {
		log.Fatal(err)
	}

	MySQLbuilt = true

}

func EnsureSqlite() {

	if Sqlitebuilt {
		return
	}

	database, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec(`DROP TABLE IF EXISTS "test_table"`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec(`CREATE TABLE "test_table" ("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL , "first_name" VARCHAR, "last_name" VARCHAR)`)
	if err != nil {
		log.Fatal(err)
	}

	database.Exec(`INSERT INTO "test_table" ("first_name","last_name") VALUES ("user1","johnson"),("user2","johnson")`)
	database.Close()

	sqlitedb, err = NewSqlite("test.db")
	if err != nil {
		log.Fatal(err)
	}

	Sqlitebuilt = true
}

func BenchmarkTable(b *testing.B) {

	user := new(TestUser)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		Table(user)

	}

}

func BenchmarkColumns(b *testing.B) {

	user := new(TestUser)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		Columns(user, EmptySkipList)

	}

}

func BenchmarkFieldPointers(b *testing.B) {

	user := new(TestUser)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		FieldPointers(user, EmptySkipList)

	}

}

func BenchmarkSqliteCreate(b *testing.B) {

	EnsureSqlite()

	user := &TestUser{
		First: "user3",
		Last:  sql.NullString{String: "johnson", Valid: true},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err := sqlitedb.Create(user)
		if err != nil {
			b.Fatal(err)
		}

	}

	b.StopTimer()

	_, err := sqlitedb.DB.Exec("DELETE FROM `test_table` WHERE `first_name` like ?", "user3")
	if err != nil {
		b.Fatal(err)
	}

}

func BenchmarkSqliteFind(b *testing.B) {

	EnsureSqlite()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		user := new(TestUser)

		err := sqlitedb.Find(user, "WHERE 1=1")
		if err != nil {
			b.Fatal(err)
		}

	}

}

func BenchmarkSqliteFindAll(b *testing.B) {

	EnsureSqlite()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		rowset := make([]*TestUser, 0, 0)

		err := sqlitedb.FindAll(&rowset, "WHERE 1=1")
		if err != nil {
			b.Fatal(err)
		}

	}

}

func BenchmarkSqliteUpdate(b *testing.B) {

	EnsureSqlite()

	user := new(TestUser)
	err := sqlitedb.Find(user, "WHERE first_name like ?", "user2")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := sqlitedb.Update(user)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkSqliteDelete(b *testing.B) {

	EnsureSqlite()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		b.StopTimer()

		id, err := sqlitedb.Create(&TestUser{
			First: "user4",
			Last:  sql.NullString{String: "johnson", Valid: true},
		})
		if err != nil {
			b.Fatal(err)
		}

		user := new(TestUser)
		err = sqlitedb.Find(user, "WHERE id = ?", id)
		if err != nil {
			b.Fatal(err)
		}

		b.StartTimer()

		err = sqlitedb.Delete(user)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkMySQLCreate(b *testing.B) {

	EnsureMySQL()

	user := &TestUser{
		First: "user3",
		Last:  sql.NullString{String: "johnson", Valid: true},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err := mysqldb.Create(user)
		if err != nil {
			b.Fatal(err)
		}

	}

	b.StopTimer()

	_, err := mysqldb.DB.Exec("DELETE FROM `test_table` WHERE `first_name` like ?", "user3")
	if err != nil {
		b.Fatal(err)
	}

}
func BenchmarkMySQLFind(b *testing.B) {

	EnsureMySQL()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		user := new(TestUser)

		err := mysqldb.Find(user, "WHERE 1=1")
		if err != nil {
			b.Fatal(err)
		}

	}

}
func BenchmarkMySQLFindAll(b *testing.B) {

	EnsureMySQL()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		rowset := make([]*TestUser, 0, 0)

		err := mysqldb.FindAll(&rowset, "WHERE 1=1")
		if err != nil {
			b.Fatal(err)
		}

	}

}

func BenchmarkMySQLUpdate(b *testing.B) {

	EnsureMySQL()

	user := new(TestUser)
	err := mysqldb.Find(user, "WHERE first_name like ?", "user2")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := mysqldb.Update(user)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkMySQLDelete(b *testing.B) {

	EnsureMySQL()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		b.StopTimer()

		id, err := mysqldb.Create(&TestUser{
			First: "user4",
			Last:  sql.NullString{String: "johnson", Valid: true},
		})
		if err != nil {
			b.Fatal(err)
		}

		user := new(TestUser)
		err = mysqldb.Find(user, "WHERE id = ?", id)
		if err != nil {
			b.Fatal(err)
		}

		b.StartTimer()

		err = mysqldb.Delete(user)
		if err != nil {
			b.Fatal(err)
		}
	}

}
