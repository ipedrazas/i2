package dckr

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDockerClient(t *testing.T) {
	client, err := NewDockerClient()
	require.NoError(t, err, "NewDockerClient should not return an error")
	assert.NotNil(t, client, "NewDockerClient should return a non-nil client")

	// Clean up
	err = client.Close()
	assert.NoError(t, err, "Closing the client should not return an error")
}

func TestListDockerContexts(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name: "test",
			want: []string{"test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListDockerContexts()

			require.NoError(t, err, "ListDockerContexts should not return an error")
			for _, ctx := range got {
				data := ctx.Endpoints
				ep := data["docker"].(map[string]any)
				host := GetHostFromEndpoint(ep)
				require.NoError(t, err, "ParseEndpoint should not return an error")
				fmt.Println(host)
			}
		})
	}
}

func TestGetActiveDockerContext(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name: "test",
			want: "unix:///Users/ivan/.docker/run/docker.sock",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetActiveDockerContext()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveDockerContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveDockerContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
