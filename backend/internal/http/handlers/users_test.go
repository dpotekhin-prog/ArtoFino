package handlers

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewUsersHandler(t *testing.T) {
	tests := []struct {
		name string
		want *UsersHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUsersHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUsersHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsersHandler_Me(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *UsersHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UsersHandler{}
			h.Me(tt.args.c)
		})
	}
}
