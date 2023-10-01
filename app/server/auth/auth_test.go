package auth_test

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/auth"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/models"
	"log"
	"reflect"
	"testing"
)

func clearSessionStore(ctx context.Context) {
	err := auth.ClearSessionStore(ctx)
	if err != nil {
		log.Fatal("could not clear session store")
	}
}

func TestAddAndGetFromSessionStore(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	session1 := models.Session{
		UserID:        "id1",
		Username:      "ivan",
		UserRole:      "User",
		IsBlacklisted: false,
	}

	sessionID, err := auth.AddToSessionStore(ctx, session1)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	session2, err := auth.GetSessionFromSessionID(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetSessionFromSessionID failed: %v", err)
	}

	if !reflect.DeepEqual(session1, session2) {
		t.Errorf("sessions are different")
	}
}

func TestRemovingSessionFromSessionStore(t *testing.T) {
	ctx := context.Background()
	clearSessionStore(ctx)

	s1 := models.Session{
		UserID:        "id1",
		Username:      "ivan",
		UserRole:      "User",
		IsBlacklisted: false,
	}

	sessionID1, err := auth.AddToSessionStore(ctx, s1)
	if err != nil {
		t.Fatalf("AddToSessionStore failed: %v", err)
	}

	s2 := models.Session{
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
		sID, err := auth.ExtractSessionIDFromSessionHashKey(str)
		if err != nil {
			t.Errorf("%v", err)
		}
		if _, found := valuesToCheck[sID]; !found {
			t.Errorf("value %s not found in the list", str)
		}
	}
}
