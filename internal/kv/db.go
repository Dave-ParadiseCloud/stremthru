package kv

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/MunifTanjim/stremthru/internal/db"
)

const TableName = "kv"

type KV struct {
	Type      string
	Key       string
	Value     string // JSON Encoded Value
	CreatedAt db.Timestamp
	UpdatedAt db.Timestamp
}

type ParsedKV[V any] struct {
	Key   string
	Value V
}

type KVStoreConfig struct {
	Type   string
	GetKey func(key string) string
}

type KVStore[V any] interface {
	Get(key string, value *V) error
	List() ([]ParsedKV[V], error)
	Set(key string, value V) error
	Del(key string) error

	WithScope(scope string) KVStore[V]
}

type SQLKVStore[V any] struct {
	t      string
	getKey func(key string) string
}

func (kv *SQLKVStore[V]) Get(key string, value *V) error {
	var val string
	query := "SELECT v FROM " + TableName + " WHERE t = ? AND k = ?"
	row := db.QueryRow(query, kv.t, kv.getKey(key))
	if err := row.Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return json.Unmarshal([]byte(val), &value)
}

func (kv *SQLKVStore[V]) List() ([]ParsedKV[V], error) {
	if kv.t == "" {
		return nil, errors.New("missing kv type value")
	}
	query := "SELECT k, v FROM " + TableName + " WHERE t = ?"
	rows, err := db.Query(query, kv.t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vs := []ParsedKV[V]{}
	for rows.Next() {
		var k string
		var v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		var val V
		if err := json.Unmarshal([]byte(v), &val); err != nil {
			return nil, err
		}
		vs = append(vs, ParsedKV[V]{
			Key:   k,
			Value: val,
		})
	}

	return vs, nil
}

func (kv *SQLKVStore[V]) Set(key string, value V) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	query := "INSERT INTO " + TableName + " (t,k,v) VALUES (?,?,?) ON CONFLICT (t,k) DO UPDATE SET v = EXCLUDED.v, uat = " + db.CurrentTimestamp
	_, err = db.Exec(query, kv.t, kv.getKey(key), val)
	return err
}

func (kv *SQLKVStore[V]) Del(key string) error {
	query := "DELETE FROM " + TableName + " WHERE t = ? AND k = ?"
	_, err := db.Exec(query, kv.t, kv.getKey(key))
	return err
}

func (kv *SQLKVStore[V]) WithScope(scope string) KVStore[V] {
	if scope == "" {
		return kv
	}
	skv := *kv
	skv.t = skv.t + ":" + scope
	return &skv
}

func NewKVStore[V any](config *KVStoreConfig) KVStore[V] {
	if config.Type != "" && config.GetKey == nil {
		config.GetKey = func(key string) string {
			return key
		}
	}

	inputKey := "key"
	outputKey := config.GetKey(inputKey)
	if config.Type == "" && outputKey == inputKey {
		panic("GetKey ouput is same as input, when type is missing")
	}
	if !strings.Contains(outputKey, inputKey) {
		panic("GetKey output does not contain input")
	}
	return &SQLKVStore[V]{
		t:      strings.ToLower(config.Type),
		getKey: config.GetKey,
	}
}
