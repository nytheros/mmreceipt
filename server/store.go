package main

import (
	"encoding/json"
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

const (
	StatusDelivered = "delivered"
	StatusRead      = "read"
)

type ReceiptRecord struct {
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	UpdatedAt int64  `json:"updated_at"`
}

type ReceiptStore struct{ api plugin.API }

func NewReceiptStore(api plugin.API) *ReceiptStore { return &ReceiptStore{api: api} }

func receiptKey(postID, userID string) string { return fmt.Sprintf("receipt:%s:%s", postID, userID) }
func receiptIndexKey(postID string) string    { return fmt.Sprintf("receipt_index:%s", postID) }

func statusRank(status string) int {
	if status == StatusRead {
		return 2
	}
	if status == StatusDelivered {
		return 1
	}
	return 0
}

func (s *ReceiptStore) Upsert(postID, userID, status string, now int64) (*ReceiptRecord, *model.AppError) {
	key := receiptKey(postID, userID)
	existingBytes, appErr := s.api.KVGet(key)
	if appErr != nil {
		return nil, appErr
	}
	if existingBytes != nil {
		var existing ReceiptRecord
		if json.Unmarshal(existingBytes, &existing) == nil && statusRank(existing.Status) >= statusRank(status) {
			return &existing, nil
		}
	}
	rec := ReceiptRecord{PostID: postID, UserID: userID, Status: status, UpdatedAt: now}
	data, err := json.Marshal(rec)
	if err != nil {
		return nil, model.NewAppError("ReceiptStore.Upsert", "plugin.readreceipt.marshal", nil, err.Error(), 500)
	}
	if appErr = s.api.KVSet(key, data); appErr != nil {
		return nil, appErr
	}
	if appErr = s.addToIndex(postID, userID); appErr != nil {
		return nil, appErr
	}
	return &rec, nil
}

func (s *ReceiptStore) addToIndex(postID, userID string) *model.AppError {
	key := receiptIndexKey(postID)
	b, appErr := s.api.KVGet(key)
	if appErr != nil {
		return appErr
	}
	users := map[string]bool{}
	if b != nil {
		_ = json.Unmarshal(b, &users)
	}
	if users[userID] {
		return nil
	}
	users[userID] = true
	data, err := json.Marshal(users)
	if err != nil {
		return model.NewAppError("ReceiptStore.addToIndex", "plugin.readreceipt.marshal", nil, err.Error(), 500)
	}
	return s.api.KVSet(key, data)
}

func (s *ReceiptStore) GetForPost(postID string) ([]ReceiptRecord, *model.AppError) {
	return s.getReceipts(receiptIndexKey(postID), postID)
}

func (s *ReceiptStore) GetForPosts(postIDs []string) (map[string][]ReceiptRecord, *model.AppError) {
	out := make(map[string][]ReceiptRecord, len(postIDs))
	for _, postID := range postIDs {
		recs, err := s.getReceipts(receiptIndexKey(postID), postID)
		if err != nil {
			return nil, err
		}
		out[postID] = recs
	}
	return out, nil
}

func (s *ReceiptStore) getReceipts(indexKey, postID string) ([]ReceiptRecord, *model.AppError) {
	b, appErr := s.api.KVGet(indexKey)
	if appErr != nil {
		return nil, appErr
	}
	if b == nil {
		return []ReceiptRecord{}, nil
	}
	users := map[string]bool{}
	_ = json.Unmarshal(b, &users)
	out := make([]ReceiptRecord, 0, len(users))
	for userID := range users {
		rb, err := s.api.KVGet(receiptKey(postID, userID))
		if err != nil {
			return nil, err
		}
		if rb == nil {
			continue
		}
		var rec ReceiptRecord
		if json.Unmarshal(rb, &rec) == nil {
			out = append(out, rec)
		}
	}
	return out, nil
}
