package engine

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const saveFileName = "save.json"
const appDirName = "poker_tui"

// SaveDir returns the platform-appropriate config directory for the app.
func SaveDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		// Fall back to home dir.
		base, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(base, appDirName), nil
}

// SavePath returns the full path to save.json.
func SavePath() (string, error) {
	dir, err := SaveDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, saveFileName), nil
}

// SaveGame writes the current game state atomically to disk.
func SaveGame(gs *GameState) error {
	dir, err := SaveDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}
	savePath := filepath.Join(dir, saveFileName)
	tmpPath := savePath + ".tmp"

	data, err := json.MarshalIndent(gs, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	// Atomic rename.
	return os.Rename(tmpPath, savePath)
}

// LoadGame reads save.json and returns a populated GameState.
// If no save file exists, returns a fresh NewGameState.
func LoadGame() (*GameState, error) {
	savePath, err := SavePath()
	if err != nil {
		return NewGameState(), nil
	}

	data, err := os.ReadFile(savePath)
	if os.IsNotExist(err) {
		return NewGameState(), nil
	}
	if err != nil {
		return NewGameState(), err
	}

	gs := NewGameState()
	if err := json.Unmarshal(data, gs); err != nil {
		// Corrupt save — start fresh.
		return NewGameState(), nil
	}

	// Migration: if save version is older, handle gracefully.
	if gs.SaveVersion != SaveVersion {
		gs.SaveVersion = SaveVersion
	}

	// Ensure Gamble.MaxStages is set.
	if gs.Gamble.MaxStages == 0 {
		gs.Gamble.MaxStages = MaxGambleStages
	}
	return gs, nil
}

// DeleteSave removes the save file (used by Reset Progress).
func DeleteSave() error {
	savePath, err := SavePath()
	if err != nil {
		return err
	}
	err = os.Remove(savePath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
