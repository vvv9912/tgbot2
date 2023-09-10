package storage_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"testing"
	"tgbotv2/internal/model"
	"tgbotv2/internal/storage"
	"time"
)

func TestUsersPostgresStorage_AddUser(t *testing.T) {

	a := model.Users{
		TgID:       222,
		StatusUser: 1,
		StateUser:  2,
		CreatedAt:  time.Now().UTC(),
		Corzine:    []int{1, 2, 3, 4},
	}
	//fmt.Println("db test:", config.Get().DatabaseDSN)
	//db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Printf("failed to connect db:%v", err)
		t.Fatalf("coonect to db_ addUser %v", err)
		return
	}
	defer db.Close()

	s := storage.NewUsersPostgresStorage(db)
	ctx := context.TODO()
	err = s.AddUser(ctx, a)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("TestUsersPostgresStorage_AddUser: %v", err)
	}
}

func TestUsersPostgresStorage_GetOrderByID(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Printf("failed to connect db:%v", err)
		t.Fatalf("coonect to db_ addUser %v", err)
		return
	}
	defer db.Close()

	s := storage.NewUsersPostgresStorage(db)
	ctx := context.TODO()
	id := 2
	corz, err := s.GetOrderByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("TestUsersPostgresStorage_AddUser: %v", err)
	}
	t.Logf("cout GetOrderById id: %v, corz:%v", id, corz)
}

func TestUsersPostgresStorage_GetOrderByTgID(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Printf("failed to connect db:%v", err)
		t.Fatalf("coonect to db_ addUser %v", err)
		return
	}
	defer db.Close()

	s := storage.NewUsersPostgresStorage(db)
	ctx := context.TODO()
	var tgid int64
	tgid = 222
	corz, err := s.GetOrderByTgID(ctx, tgid)
	if err != nil {
		fmt.Println(err)
		t.Fatalf("GetOrderByTgID: %v", err)
	}
	t.Logf("cout GetOrderByTgID tgid: %v, corz:%v", tgid, corz)
}
func TestUsersPostgresStorage_GetStatusUserByTgID(t *testing.T) {

	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	if err != nil {
		log.Printf("failed to connect db:%v", err)
		t.Fatalf("coonect to db_ addUser %v", err)
		return
	}
	defer db.Close()

	s := storage.NewUsersPostgresStorage(db)
	ctx := context.TODO()
	var tgid int64
	tgid = 14124
	status, state, err := s.GetStatusUserByTgID(ctx, tgid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.Logf("ok, status:%v, state:%v", status, state)
			err = nil
		} else {
			fmt.Println(err)
			t.Fatalf("GetOrderByTgID: %v", err)
		}
	}
	t.Logf("cout GetOrderByTgID tgid: %v, status:%v, state:%v", tgid, status, state)
}
