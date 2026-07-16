package config

import (
	"errors"
	"strings"
	"testing"
)

func TestLoadStoreDBPath(t *testing.T) {
	tests := []struct {
		name string
		args []string
		env  map[string]string
		want string
	}{
		{
			name: "default when neither flag nor env is set",
			want: defaultDBPath,
		},
		{
			name: "environment variable overrides the default",
			env:  map[string]string{"BOLTE_BRIDGE_DB_PATH": "/var/lib/bridge/env.db"},
			want: "/var/lib/bridge/env.db",
		},
		{
			name: "long flag overrides both the environment and the default",
			args: []string{"--db-path", "/tmp/flag.db"},
			env:  map[string]string{"BOLTE_BRIDGE_DB_PATH": "/var/lib/bridge/env.db"},
			want: "/tmp/flag.db",
		},
		{
			name: "short flag overrides both the environment and the default",
			args: []string{"-d", "/tmp/flag.db"},
			env:  map[string]string{"BOLTE_BRIDGE_DB_PATH": "/var/lib/bridge/env.db"},
			want: "/tmp/flag.db",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			cfg, err := Load(tc.args, DefaultSections...)
			if err != nil {
				t.Fatalf("Load returned error: %v", err)
			}
			if got := cfg.Store.SQLite.Path; got != tc.want {
				t.Errorf("Store.SQLite.Path = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestLoadUnknownFlagErrors(t *testing.T) {
	if _, err := Load([]string{"--nope"}, DefaultSections...); err == nil {
		t.Fatal("Load with an unknown flag returned nil error, want non-nil")
	}
}

func TestLoadValidationErrorAborts(t *testing.T) {
	// An empty db.path is rejected by storeSection's ApplyFunc.
	_, err := Load([]string{"--db-path", ""}, DefaultSections...)
	if err == nil {
		t.Fatal("Load with empty db.path returned nil error, want non-nil")
	}
	if !strings.Contains(err.Error(), "db.path") {
		t.Errorf("error %q does not mention the offending key db.path", err)
	}
}

func TestLoadBindingErrorAborts(t *testing.T) {
	failing := func(_ *Binder) (ApplyFunc, error) {
		return nil, errors.New("boom")
	}
	if _, err := Load(nil, failing); err == nil {
		t.Fatal("Load with a failing section returned nil error, want non-nil")
	}
}
