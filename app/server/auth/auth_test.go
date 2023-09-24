package auth_test

import (
	"context"
	auth2 "github.com/IliyaYavorovPetrov/api-gateway/app/server/auth"
	"log"
	"reflect"
	"testing"
)

func clearSessionStore(ctx context.Context) {
	err := auth2.ClearSessionStore(ctx)
	if err != nil {
		log.Fatal("could not clear session store")
	}
}

func TestAddAndGetFromSessionStore(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	s1 := &auth2.Session{
		UserID:        "id1",
		Username:      "ivan",
		UserRole:      "User",
		IsBlacklisted: false,
	}

	sessionID, err := auth2.AddToSessionStore(ctx, s1)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	s2, err := auth2.GetSessionFromSessionID(ctx, sessionID)
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

	s1 := &auth2.Session{
		UserID:        "id1",
		Username:      "ivan",
		UserRole:      "User",
		IsBlacklisted: false,
	}

	sessionID1, err := auth2.AddToSessionStore(ctx, s1)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	s2 := &auth2.Session{
		UserID:        "id2",
		Username:      "gosho",
		UserRole:      "Admin",
		IsBlacklisted: false,
	}

	sessionID2, err := auth2.AddToSessionStore(ctx, s2)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	valuesToCheck := make(map[string]struct{})
	valuesToCheck[sessionID1] = struct{}{}
	valuesToCheck[sessionID2] = struct{}{}

	allSessionIDs, err := auth2.GetAllSessionIDs(ctx)
	if err != nil {
		t.Fatalf("GetAllSessionIDs failed: %v", err)
	}

	if len(allSessionIDs) != len(valuesToCheck) {
		t.Errorf("wrong number of session ids, %d expected, %d received", len(allSessionIDs), len(valuesToCheck))
	}

	for _, str := range allSessionIDs {
		sID, err := auth2.ExtractSessionIDFromSessionHashKey(str)
		if err != nil {
			t.Errorf("%v", err)
		}
		if _, found := valuesToCheck[sID]; !found {
			t.Errorf("value %s not found in the list", str)
		}
	}
}

func TestChangeBlacklistStatusUser(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	s1 := &auth2.Session{
		UserID:        "id",
		Username:      "ivan",
		UserRole:      "User",
		IsBlacklisted: false,
	}

	sessionID, err := auth2.AddToSessionStore(ctx, s1)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	err = auth2.ChangeBlacklistStatusOfSession(ctx, sessionID, true)
	if err != nil {
		t.Fatalf("ChangeBlacklistStatusOfSession failed: %v", err)
	}

	s2, err := auth2.GetSessionFromSessionID(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetSessionFromSessionID failed: %v", err)
	}

	if !s2.IsBlacklisted {
		t.Errorf("status is not changed")
	}
}
