package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SecretManagerService struct {
	mu       sync.Mutex
	filePath string
	keyPath  string
}

type SecretItem struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Type      string   `json:"type"`
	Account   string   `json:"account"`
	URL       string   `json:"url"`
	Tags      []string `json:"tags"`
	Note      string   `json:"note"`
	HasSecret bool     `json:"hasSecret"`
	HasExtra  bool     `json:"hasExtra"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

type SecretInput struct {
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Account     string   `json:"account"`
	Secret      string   `json:"secret"`
	SecretExtra string   `json:"secretExtra"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
	Note        string   `json:"note"`
}

type secretStore struct {
	Items []storedSecretItem `json:"items"`
}

type storedSecretItem struct {
	ID                    string   `json:"id"`
	Title                 string   `json:"title"`
	Type                  string   `json:"type"`
	Account               string   `json:"account"`
	URL                   string   `json:"url"`
	Tags                  []string `json:"tags"`
	Note                  string   `json:"note"`
	SecretCiphertext      string   `json:"secretCiphertext"`
	SecretExtraCiphertext string   `json:"secretExtraCiphertext"`
	CreatedAt             string   `json:"createdAt"`
	UpdatedAt             string   `json:"updatedAt"`
}

func NewSecretManagerService() *SecretManagerService {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}

	dataDir := filepath.Join(configDir, "MyTools")
	return &SecretManagerService{
		filePath: filepath.Join(dataDir, "secrets.json"),
		keyPath:  filepath.Join(dataDir, "secrets.key"),
	}
}

func (s *SecretManagerService) List() ([]SecretItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.readStore()
	if err != nil {
		return nil, err
	}

	items := make([]SecretItem, 0, len(store.Items))
	for _, item := range store.Items {
		items = append(items, publicSecretItem(item))
	}
	sortSecretItems(items)
	return items, nil
}

func (s *SecretManagerService) Create(input SecretInput) (SecretItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedInput, err := normalizeSecretInput(input)
	if err != nil {
		return SecretItem{}, err
	}

	encryptedSecret, err := s.encryptSecret(normalizedInput.Secret)
	if err != nil {
		return SecretItem{}, err
	}
	encryptedExtraSecret, err := s.encryptOptionalSecret(normalizedInput.SecretExtra)
	if err != nil {
		return SecretItem{}, err
	}

	store, err := s.readStore()
	if err != nil {
		return SecretItem{}, err
	}

	now := time.Now().Format(time.RFC3339)
	item := storedSecretItem{
		ID:                    uuid.NewString(),
		Title:                 normalizedInput.Title,
		Type:                  normalizedInput.Type,
		Account:               normalizedInput.Account,
		URL:                   normalizedInput.URL,
		Tags:                  normalizedInput.Tags,
		Note:                  normalizedInput.Note,
		SecretCiphertext:      encryptedSecret,
		SecretExtraCiphertext: encryptedExtraSecret,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	store.Items = append(store.Items, item)
	sortStoredSecretItems(store.Items)
	if err := s.writeStore(store); err != nil {
		return SecretItem{}, err
	}

	return publicSecretItem(item), nil
}

func (s *SecretManagerService) Update(id string, input SecretInput) (SecretItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedInput, err := normalizeSecretInput(input)
	if err != nil {
		return SecretItem{}, err
	}

	encryptedSecret, err := s.encryptSecret(normalizedInput.Secret)
	if err != nil {
		return SecretItem{}, err
	}
	encryptedExtraSecret, err := s.encryptOptionalSecret(normalizedInput.SecretExtra)
	if err != nil {
		return SecretItem{}, err
	}

	store, err := s.readStore()
	if err != nil {
		return SecretItem{}, err
	}

	for index := range store.Items {
		if store.Items[index].ID == id {
			store.Items[index].Title = normalizedInput.Title
			store.Items[index].Type = normalizedInput.Type
			store.Items[index].Account = normalizedInput.Account
			store.Items[index].URL = normalizedInput.URL
			store.Items[index].Tags = normalizedInput.Tags
			store.Items[index].Note = normalizedInput.Note
			store.Items[index].SecretCiphertext = encryptedSecret
			store.Items[index].SecretExtraCiphertext = encryptedExtraSecret
			store.Items[index].UpdatedAt = time.Now().Format(time.RFC3339)

			item := store.Items[index]
			sortStoredSecretItems(store.Items)
			if err := s.writeStore(store); err != nil {
				return SecretItem{}, err
			}

			return publicSecretItem(item), nil
		}
	}

	return SecretItem{}, errors.New("账号密钥不存在")
}

func (s *SecretManagerService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.readStore()
	if err != nil {
		return err
	}

	nextItems := store.Items[:0]
	deleted := false
	for _, item := range store.Items {
		if item.ID == id {
			deleted = true
			continue
		}
		nextItems = append(nextItems, item)
	}
	if !deleted {
		return errors.New("账号密钥不存在")
	}

	store.Items = nextItems
	return s.writeStore(store)
}

func (s *SecretManagerService) RevealSecret(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.readStore()
	if err != nil {
		return "", err
	}

	for _, item := range store.Items {
		if item.ID == id {
			return s.decryptSecret(item.SecretCiphertext)
		}
	}

	return "", errors.New("账号密钥不存在")
}

func (s *SecretManagerService) RevealSecretExtra(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.readStore()
	if err != nil {
		return "", err
	}

	for _, item := range store.Items {
		if item.ID == id {
			return s.decryptSecret(item.SecretExtraCiphertext)
		}
	}

	return "", errors.New("账号密钥不存在")
}

func (s *SecretManagerService) StorePath() string {
	return s.filePath
}

func (s *SecretManagerService) readStore() (secretStore, error) {
	content, err := os.ReadFile(s.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return secretStore{Items: []storedSecretItem{}}, nil
	}
	if err != nil {
		return secretStore{}, fmt.Errorf("读取账号密钥失败: %w", err)
	}
	if len(strings.TrimSpace(string(content))) == 0 {
		return secretStore{Items: []storedSecretItem{}}, nil
	}

	var store secretStore
	if err := json.Unmarshal(content, &store); err != nil {
		return secretStore{}, fmt.Errorf("解析账号密钥失败: %w", err)
	}
	if store.Items == nil {
		store.Items = []storedSecretItem{}
	}

	return store, nil
}

func (s *SecretManagerService) writeStore(store secretStore) error {
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return fmt.Errorf("创建账号密钥目录失败: %w", err)
	}

	content, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化账号密钥失败: %w", err)
	}

	tempPath := s.filePath + ".tmp"
	if err := os.WriteFile(tempPath, content, 0o600); err != nil {
		return fmt.Errorf("写入账号密钥失败: %w", err)
	}
	if err := os.Rename(tempPath, s.filePath); err != nil {
		return fmt.Errorf("保存账号密钥失败: %w", err)
	}

	return nil
}

func (s *SecretManagerService) encryptSecret(secret string) (string, error) {
	key, err := s.loadOrCreateKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建加密器失败: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建加密模式失败: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成随机数失败: %w", err)
	}

	ciphertext := aead.Seal(nonce, nonce, []byte(secret), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *SecretManagerService) encryptOptionalSecret(secret string) (string, error) {
	if secret == "" {
		return "", nil
	}
	return s.encryptSecret(secret)
}

func (s *SecretManagerService) decryptSecret(encodedSecret string) (string, error) {
	if encodedSecret == "" {
		return "", nil
	}

	key, err := s.loadOrCreateKey()
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encodedSecret)
	if err != nil {
		return "", errors.New("密钥数据已损坏")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建解密器失败: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建解密模式失败: %w", err)
	}
	if len(ciphertext) < aead.NonceSize() {
		return "", errors.New("密钥数据不完整")
	}

	nonce := ciphertext[:aead.NonceSize()]
	encryptedValue := ciphertext[aead.NonceSize():]
	plaintext, err := aead.Open(nil, nonce, encryptedValue, nil)
	if err != nil {
		return "", errors.New("密钥解密失败")
	}

	return string(plaintext), nil
}

func (s *SecretManagerService) loadOrCreateKey() ([]byte, error) {
	key, err := os.ReadFile(s.keyPath)
	if err == nil {
		decodedKey, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(key)))
		if err != nil || len(decodedKey) != 32 {
			return nil, errors.New("本地加密密钥无效")
		}
		return decodedKey, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("读取本地加密密钥失败: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(s.keyPath), 0o755); err != nil {
		return nil, fmt.Errorf("创建密钥目录失败: %w", err)
	}

	newKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, newKey); err != nil {
		return nil, fmt.Errorf("生成本地加密密钥失败: %w", err)
	}

	if err := os.WriteFile(s.keyPath, []byte(base64.StdEncoding.EncodeToString(newKey)), 0o600); err != nil {
		return nil, fmt.Errorf("保存本地加密密钥失败: %w", err)
	}

	return newKey, nil
}

func normalizeSecretInput(input SecretInput) (SecretInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Type = strings.TrimSpace(input.Type)
	input.Account = strings.TrimSpace(input.Account)
	input.Secret = strings.TrimSpace(input.Secret)
	input.SecretExtra = strings.TrimSpace(input.SecretExtra)
	input.URL = strings.TrimSpace(input.URL)
	input.Note = strings.TrimSpace(input.Note)
	input.Tags = normalizeTags(input.Tags)

	if input.Type == "" || input.Type == "token" || input.Type == "api-key" {
		input.Type = "password"
	}
	if input.Type != "password" && input.Type != "ak-sk" {
		return SecretInput{}, errors.New("不支持的密钥类型")
	}
	if input.Type == "ak-sk" {
		input.Account = ""
	}
	if input.Title == "" {
		return SecretInput{}, errors.New("请输入名称")
	}
	if input.Secret == "" {
		return SecretInput{}, errors.New("请输入密钥或密码")
	}
	if input.Type == "ak-sk" && input.SecretExtra == "" {
		return SecretInput{}, errors.New("请输入 Secret Key")
	}

	return input, nil
}

func publicSecretItem(item storedSecretItem) SecretItem {
	return SecretItem{
		ID:        item.ID,
		Title:     item.Title,
		Type:      normalizeSecretType(item.Type),
		Account:   item.Account,
		URL:       item.URL,
		Tags:      item.Tags,
		Note:      item.Note,
		HasSecret: item.SecretCiphertext != "",
		HasExtra:  item.SecretExtraCiphertext != "",
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func normalizeSecretType(secretType string) string {
	if secretType == "" || secretType == "token" || secretType == "api-key" {
		return "password"
	}
	return secretType
}

func sortSecretItems(items []SecretItem) {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].UpdatedAt > items[j].UpdatedAt
	})
}

func sortStoredSecretItems(items []storedSecretItem) {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].UpdatedAt > items[j].UpdatedAt
	})
}
