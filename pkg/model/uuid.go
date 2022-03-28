package model

import "github.com/google/uuid"

// UUID to generate the unique id
func UUID() string {
	return uuid.New().String()
}
