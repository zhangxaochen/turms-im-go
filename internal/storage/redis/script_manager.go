package redis

	import (
		"context"
		"fmt"
		"os"
		"path/filepath"
		"strings"

		"github.com/redis/go-redis/v9"
	)

	type ScriptManager struct {
		client  *Client
		scripts map[string]*redis.Script
	}

	func NewScriptManager(client *Client, scriptDir string, placeholders map[string]string) (*ScriptManager, error) {
		scripts := make(map[string]*redis.Script)

		entries, err := os.ReadDir(scriptDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read script directory: %w", err)
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".lua") {
				continue
			}

			path := filepath.Join(scriptDir, entry.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read script %s: %w", entry.Name(), err)
			}

			scriptContent := string(content)
			for k, v := range placeholders {
				scriptContent = strings.ReplaceAll(scriptContent, k, v)
			}

			name := strings.TrimSuffix(entry.Name(), ".lua")
			scripts[name] = redis.NewScript(scriptContent)
		}

		return &ScriptManager{
			client:  client,
			scripts: scripts,
		}, nil
	}

	func (m *ScriptManager) Run(ctx context.Context, name string, keys []string, args ...interface{}) *redis.Cmd {
		script, ok := m.scripts[name]
		if !ok {
			return redis.NewCmd(ctx, fmt.Errorf("script %s not found", name))
		}
		return script.Run(ctx, m.client.RDB, keys, args...)
	}
	
