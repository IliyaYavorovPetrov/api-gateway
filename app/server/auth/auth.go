package auth

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/distributed"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/local"
	"github.com/IliyaYavorovPetrov/api-gateway/app/server/models"
	"log"
	"strings"

	"github.com/google/uuid"
)

var prefixAuthSession = "auth:session:"

var localCache gateways.Cache
var distributedCache gateways.Cache

func init() {
	localCache = local.CreateInstance("auth-local-cache")
	distributedCache = distributed.CreateInstance("auth-distributed-cache")
}

func createSessionHashKey(sessionID string) string {
	return prefixAuthSession + sessionID
}

func ExtractSessionIDFromSessionHashKey(s string) (string, error) {
	if strings.HasPrefix(s, prefixAuthSession) {
		return s[len(prefixAuthSession):], nil
	}

	return s, ErrNotValidSessionHashKey
}

func AddToSessionStore(ctx context.Context, s *models.Session) (string, error) {
	sessionID := uuid.New().String()

	err := distributedCache.Add(ctx, createSessionHashKey(sessionID), s)
	if err != nil {
		log.Fatalf("failed to add an auth session %s", err)
		return "", err
	}

	return sessionID, nil
}

func GetSessionFromSessionID(ctx context.Context, sessionID string) (models.Session, error) {
	session, err := distributedCache.Get(ctx, createSessionHashKey(sessionID))
	if err != nil {
		return models.Session{}, err
	}

	return session.(models.Session), nil
}

func GetAllSessionIDs(ctx context.Context) ([]string, error) {
	sessions, err := distributedCache.GetAllKeysByPrefix(ctx, prefixAuthSession)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func RemoveSessionFromSessionStore(ctx context.Context, sessionID string) error {
	err := distributedCache.Delete(ctx, createSessionHashKey(sessionID))
	if err != nil {
		return err
	}

	return nil
}

func ClearSessionStore(ctx context.Context) error {
	sessionIDs, err := GetAllSessionIDs(ctx)
	if err != nil {
		return err
	}

	for _, sessionID := range sessionIDs {
		err = distributedCache.Delete(ctx, sessionID)

		if err != nil {
			log.Fatalf("failed to delete a session %s", err)
		}
	}

	return nil
}
