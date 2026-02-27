package cli_test

import (
	"testing"

	"github.com/aaronshifman/lsm-cause-i-can/pkg/cli"
	"github.com/stretchr/testify/require"
)

func TestCli_Execute_GET(t *testing.T) { // nolint: funlen
	tests := []struct {
		name    string
		setup   []string
		command string
		wantVal string
		wantErr bool
	}{
		{
			name:    "get existing key",
			setup:   []string{"PUT foo bar"},
			command: "GET foo",
			wantVal: "bar",
			wantErr: false,
		},
		{
			name:    "get missing key",
			setup:   []string{},
			command: "GET foo",
			wantVal: "",
			wantErr: true,
		},
		{
			name:    "get with no key argument",
			setup:   []string{},
			command: "GET",
			wantVal: "",
			wantErr: true,
		},
		{
			name:    "get after overwrite returns latest value",
			setup:   []string{"PUT foo first", "PUT foo second"},
			command: "GET foo",
			wantVal: "second",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cli.NewCli()

			for _, cmd := range tt.setup {
				inter, err := c.Parse(cmd)
				require.NoError(t, err)
				_, err = c.Execute(inter)
				require.NoError(t, err)
			}

			inter, err := c.Parse(tt.command)
			require.NoError(t, err)

			got, err := c.Execute(inter)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantVal, got)
		})
	}
}

func TestCli_Execute_PUT(t *testing.T) {
	tests := []struct {
		name    string
		command string
		wantOut string
		wantErr bool
	}{
		{
			name:    "put new key",
			command: "PUT foo bar",
			wantOut: "OK",
			wantErr: false,
		},
		{
			name:    "put value with spaces",
			command: "PUT greeting hello world",
			wantOut: "OK",
			wantErr: false,
		},
		{
			name:    "put missing value",
			command: "PUT foo",
			wantOut: "",
			wantErr: true,
		},
		{
			name:    "put missing key and value",
			command: "PUT",
			wantOut: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cli.NewCli()

			inter, err := c.Parse(tt.command)
			require.NoError(t, err)

			got, err := c.Execute(inter)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantOut, got)
		})
	}
}

func TestCli_Execute_QUIT(t *testing.T) {
	c := cli.NewCli()

	inter, err := c.Parse("QUIT")
	require.NoError(t, err)

	_, err = c.Execute(inter)
	require.ErrorIs(t, err, cli.ErrQuit)
}
