package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Verification struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type VerificationRepository interface {
	Create(email, hash string) error
	FindByHash(hash string) (*Verification, error)
	Delete(hash string) error
	LoadAll() ([]Verification, error)
}

type FileVerificationRepo struct {
	filePath string
	mu       sync.Mutex // чтобы избежать гонок при записи
}

func NewFileVerificationRepo(path string) *FileVerificationRepo {
	return &FileVerificationRepo{filePath: path}
}

func (r *FileVerificationRepo) getFilePath() string {
	if r.filePath == "" {
		return filepath.Join(".", "verifications.json")
	}
	return r.filePath
}

func (r *FileVerificationRepo) LoadAll() ([]Verification, error) {
	path := r.getFilePath()
	var data []Verification

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil // если файла нет — возвращаем пустой список
		}
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&data); err != nil && err.Error() != "EOF" {
		return nil, err
	}
	return data, nil
}

func (r *FileVerificationRepo) Create(email, hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Загружаем всё
	all, err := r.LoadAll()
	if err != nil {
		return err
	}

	updated := false
	for i, v := range all {
		if v.Email == email {
			all[i].Hash = hash
			updated = true
			break
		}
	}

	if !updated {
		all = append(all, Verification{Email: email, Hash: hash})
	}

	f, err := os.Create(r.getFilePath())
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(all)
}

func (r *FileVerificationRepo) FindByHash(hash string) (*Verification, error) {
	all, err := r.LoadAll()
	if err != nil {
		return nil, err
	}

	for _, v := range all {
		if v.Hash == hash {
			return &v, nil
		}
	}
	return nil, nil
}

func (r *FileVerificationRepo) Delete(hash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	all, err := r.LoadAll()
	if err != nil {
		return err
	}

	newList := make([]Verification, 0, len(all))
	for _, v := range all {
		if v.Hash != hash {
			newList = append(newList, v)
		}
	}

	f, err := os.Create(r.getFilePath())
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(newList)
}
