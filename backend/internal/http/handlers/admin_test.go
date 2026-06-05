package handlers

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewAdminHandler(t *testing.T) {
	tests := []struct {
		name string
		want *AdminHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAdminHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAdminHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminHandler_Stats(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AdminHandler{}
			h.Stats(tt.args.c)
		})
	}
}

func TestAdminHandler_Ping(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AdminHandler{}
			h.Ping(tt.args.c)
		})
	}
}
