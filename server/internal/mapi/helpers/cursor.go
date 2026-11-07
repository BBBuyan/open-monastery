package helpers

import (
    "encoding/base64"
    "encoding/json"
)

type Cursor struct {
    ID int64 `json:"id"`
}

func EncodeCursor(cursor Cursor) (string, error) {
    data, err := json.Marshal(cursor)
    if err != nil {
        return "", err
    }

    return base64.RawURLEncoding.EncodeToString(data), nil
}

func DecodeCursor(encoded string) (Cursor, error) {
    var cursor Cursor

    data, err := base64.RawURLEncoding.DecodeString(encoded)
    if err != nil {
        return cursor, err
    }

    err = json.Unmarshal(data, &cursor)
    return cursor, err
}
