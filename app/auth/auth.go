package auth

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var prefixAuthSession string = "auth:session:"

func createSessionHashKey(sessionID string) string {
	return prefixAuthSession + sessionID
}

func GetSessionIDFromSessionHashKey(s string) (string, error) {
	if strings.HasPrefix(s, prefixAuthSession) {
		return s[len(prefixAuthSession):], nil
	}

	return s, ErrNotValidSessionHashKey
}

func AddToSessionStore(ctx context.Context, s *Session) (string, error) {
	sessionID := uuid.New().String()

	if _, err := rdb.HSet(ctx, createSessionHashKey(sessionID), s).Result(); err != nil {
		log.Fatalf("failed to add an auth session %s", err)
		return "", err
	}

	return sessionID, nil
}

func GetSessionFromSessionID(ctx context.Context, sessionID string) (*Session, error) {
	var s Session
	err := rdb.HGetAll(ctx, createSessionHashKey(sessionID)).Scan(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func GetAllSessionIDs(ctx context.Context) ([]string, error) {
	var sessions []string
	iter := rdb.Scan(ctx, 0, prefixAuthSession+"*", 0).Iterator()
	for iter.Next(ctx) {
		sessions = append(sessions, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func ChangeBlacklistStatusOfSession(ctx context.Context, sessionID string, blackListStatus bool) error {
	if _, err := rdb.HSet(ctx, createSessionHashKey(sessionID), "isBlacklisted", blackListStatus).Result(); err != nil {
		log.Fatalf("failed to create an auth session %s", err)
		return err
	}

	return nil
}

func RemoveSessionFromSessionStore(ctx context.Context, sessionID string) error {
	err := rdb.HDel(ctx, createSessionHashKey(sessionID)).Err()
	if err != nil {
		return err
	}

	return nil
}

func ClearSessionStore(ctx context.Context) {
	// TODO: Delete by pattern
	rdb.FlushAll(ctx)
}
