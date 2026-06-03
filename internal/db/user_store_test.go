package db

import (
	//"context"
	"testing"
	"time"

	"github.com/MdRasB/LogLine/internal/model"
	"github.com/google/uuid"
)

func TestCreateAndGetUser(t *testing.T) {
	pool := setupTestDB(t)

	store := NewUserStore(pool)

	user := model.User{
		Id:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed-password",
		CreatedAt:    time.Now(),
	}

	err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	foundByEmail, err := store.GetUserByEmail(user.Email)
	if err != nil {
		t.Fatalf("get user by email: %v", err)
	}

	if foundByEmail.Email != user.Email {
		t.Fatalf(
			"expected email %s got %s",
			user.Email,
			foundByEmail.Email,
		)
	}

	foundByID, err := store.GetUserByID(user.Id)
	if err != nil {
		t.Fatalf("get user by id: %v", err)
	}

	if foundByID.Id != user.Id {
		t.Fatalf(
			"expected id %v got %v",
			user.Id,
			foundByID.Id,
		)
	}
}