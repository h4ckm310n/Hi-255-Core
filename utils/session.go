package utils

import "github.com/google/uuid"

var sessions map[string]SessionObject

type SessionObject interface{}

func init() {
	sessions = make(map[string]SessionObject)
}

func GenerateSession() []byte {
	session, _ := uuid.New().MarshalText()
	return session
}

func SetSession(session string, object SessionObject) {
	sessions[session] = object
}

func GetSessionValue(session string) SessionObject {
	return sessions[session]
}

func DeleteSession(session string) {
	delete(sessions, session)
}
