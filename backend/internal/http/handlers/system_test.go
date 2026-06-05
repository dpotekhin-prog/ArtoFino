package handlers

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewSystemHandler(t *testing.T) {
	tests := []struct {
		name string
		want *SystemHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSystemHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystemHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemHandler_SetReady(t *testing.T) {
	type fields struct {
		ready bool
	}
	type args struct {
		v bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SystemHandler{
				ready: tt.fields.ready,
			}
			h.SetReady(tt.args.v)
		})
	}
}

func TestSystemHandler_HealthCheck(t *testing.T) {
	type fields struct {
		ready bool
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SystemHandler{
				ready: tt.fields.ready,
			}
			h.HealthCheck(tt.args.w, tt.args.r)
		})
	}
}

func TestSystemHandler_LivenessCheck(t *testing.T) {
	type fields struct {
		ready bool
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SystemHandler{
				ready: tt.fields.ready,
			}
			h.LivenessCheck(tt.args.w, tt.args.r)
		})
	}
}

func TestSystemHandler_ReadinessCheck(t *testing.T) {
	type fields struct {
		ready bool
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SystemHandler{
				ready: tt.fields.ready,
			}
			h.ReadinessCheck(tt.args.w, tt.args.r)
		})
	}
}
