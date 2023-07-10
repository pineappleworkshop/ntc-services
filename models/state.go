package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type State struct {
	Path               string `json:"path"`
	ScannerBlockHeight int64  `json:"scanner_block_height"`
}

func NewState(path string) *State {
	return &State{
		Path:               path,
		ScannerBlockHeight: -1,
	}
}

func (s *State) Read() error {
	data, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return nil
}

func (s *State) Write() error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(s.Path, data, 0644); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	return nil
}
