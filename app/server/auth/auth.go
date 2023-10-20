package auth

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/common/models"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/distributed"
	"github.com/IliyaYavorovPetrov/api-gateway/app/gateways/cache/local"
	"log"
	"strings"

	"github.com/google/uuid"
)

var prefixAuthSession = "auth:session:"

var localCache gateways.Cache[models.Session]
var distributedCache gateways.Cache[models.Session]

func Init(ctx context.Context) {
	localCache = local.New[models.Session]("auth-local-cache")
	distributedCache = distributed.New[models.Session]("auth-distributed-cache")

	err := cache.SyncFromTo[models.Session](ctx, localCache, distributedCache)
	if err != nil {
		panic(err)
	}
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

func AddToSessionStore(ctx context.Context, session models.Session) (string, error) {
	sessionID := uuid.New().String()

	err := distributedCache.Add(ctx, createSessionHashKey(sessionID), session)
	if err != nil {
		log.Fatalf("failed to add an auth session %session", err)
		return "", err
	}

	return sessionID, nil
}

func GetSessionFromSessionID(ctx context.Context, sessionID string) (models.Session, error) {
	session, err := distributedCache.Get(ctx, createSessionHashKey(sessionID))
	if err != nil {
		return models.Session{}, err
	}

	return *session, nil
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
