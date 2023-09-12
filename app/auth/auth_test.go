package auth_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/IliyaYavorovPetrov/api-gateway/app/auth"
	"github.com/IliyaYavorovPetrov/api-gateway/app/common/models"
)

func clearSessionStore(ctx context.Context) {
	auth.ClearSessionStore(ctx)
}

func TestAddAndGetFromSessionStore(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	s1 := &models.Session{
		UserID:        "id",
		Username:      "ivan",
		UserRole:      "User",
		IsBlacklisted: false,
	}

	sessionID, err := auth.AddToSessionStore(ctx, s1)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	s2, err := auth.GetSessionFromSessionID(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetSessionFromSessionID failed: %v", err)
	}

	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("sessions are different")
	}
}

func TestRemovingSessionFromSessionStore(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	s1 := &models.Session{
		UserID:        "id1",
		Username:      "ivan",
		UserRole:      "User",
		IsBlacklisted: false,
	}

	sessionID1, err := auth.AddToSessionStore(ctx, s1)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	s2 := &models.Session{
		UserID:        "id2",
		Username:      "gosho",
		UserRole:      "Admin",
		IsBlacklisted: false,
	}

	sessionID2, err := auth.AddToSessionStore(ctx, s2)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	valuesToCheck := make(map[string]struct{})
	valuesToCheck[sessionID1] = struct{}{}
	valuesToCheck[sessionID2] = struct{}{}

	allSessionIDs, err := auth.GetAllSessionIDs(ctx)
	if err != nil {
		t.Fatalf("GetAllSessionIDs failed: %v", err)
	}

	if len(allSessionIDs) != len(valuesToCheck) {
		t.Errorf("wrong number of session ids, %d expected, %d received", len(allSessionIDs), len(valuesToCheck))
	}

	for _, str := range allSessionIDs {
		if _, found := valuesToCheck[auth.GetSessionIdFromSessionHashKey(str)]; !found {
			t.Errorf("value %s not found in the list", str)
		}
	}
}

func TestChangeBlacklistStatusUser(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	s1 := &models.Session{
		UserID:        "id",
		Username:      "ivan",
		UserRole:      "User",
		IsBlacklisted: false,
	}

	sessionID, err := auth.AddToSessionStore(ctx, s1)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	err = auth.ChangeBlacklistStatusOfSession(ctx, sessionID, true)
	if err != nil {
		t.Fatalf("ChangeBlacklistStatusOfSession failed: %v", err)
	}

	s2, err := auth.GetSessionFromSessionID(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetSessionFromSessionID failed: %v", err)
	}

	if !s2.IsBlacklisted {
		t.Errorf("status is not changed")
	}
}
