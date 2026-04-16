package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/browser"
)

type BookmarkService struct {
	mu       sync.Mutex
	filePath string
}

type Bookmark struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Tags      []string `json:"tags"`
	Note      string   `json:"note"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

type BookmarkInput struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags"`
	Note  string   `json:"note"`
}

type bookmarkStore struct {
	Bookmarks []Bookmark `json:"bookmarks"`
}

func NewBookmarkService() *BookmarkService {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}

	return &BookmarkService{
		filePath: filepath.Join(configDir, "MyTools", "bookmarks.json"),
	}
}

func (s *BookmarkService) List() ([]Bookmark, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.readStore()
	if err != nil {
		return nil, err
	}

	sortBookmarks(store.Bookmarks)
	return store.Bookmarks, nil
}

func (s *BookmarkService) Create(input BookmarkInput) (Bookmark, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedInput, err := normalizeBookmarkInput(input)
	if err != nil {
		return Bookmark{}, err
	}

	store, err := s.readStore()
	if err != nil {
		return Bookmark{}, err
	}

	now := time.Now().Format(time.RFC3339)
	bookmark := Bookmark{
		ID:        uuid.NewString(),
		Title:     normalizedInput.Title,
		URL:       normalizedInput.URL,
		Tags:      normalizedInput.Tags,
		Note:      normalizedInput.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}

	store.Bookmarks = append(store.Bookmarks, bookmark)
	sortBookmarks(store.Bookmarks)

	if err := s.writeStore(store); err != nil {
		return Bookmark{}, err
	}

	return bookmark, nil
}

func (s *BookmarkService) Update(id string, input BookmarkInput) (Bookmark, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedInput, err := normalizeBookmarkInput(input)
	if err != nil {
		return Bookmark{}, err
	}

	store, err := s.readStore()
	if err != nil {
		return Bookmark{}, err
	}

	for index := range store.Bookmarks {
		if store.Bookmarks[index].ID == id {
			store.Bookmarks[index].Title = normalizedInput.Title
			store.Bookmarks[index].URL = normalizedInput.URL
			store.Bookmarks[index].Tags = normalizedInput.Tags
			store.Bookmarks[index].Note = normalizedInput.Note
			store.Bookmarks[index].UpdatedAt = time.Now().Format(time.RFC3339)

			bookmark := store.Bookmarks[index]
			sortBookmarks(store.Bookmarks)
			if err := s.writeStore(store); err != nil {
				return Bookmark{}, err
			}

			return bookmark, nil
		}
	}

	return Bookmark{}, errors.New("书签不存在")
}

func (s *BookmarkService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.readStore()
	if err != nil {
		return err
	}

	nextBookmarks := store.Bookmarks[:0]
	deleted := false
	for _, bookmark := range store.Bookmarks {
		if bookmark.ID == id {
			deleted = true
			continue
		}
		nextBookmarks = append(nextBookmarks, bookmark)
	}

	if !deleted {
		return errors.New("书签不存在")
	}

	store.Bookmarks = nextBookmarks
	return s.writeStore(store)
}

func (s *BookmarkService) Open(id string) error {
	s.mu.Lock()
	store, err := s.readStore()
	if err != nil {
		s.mu.Unlock()
		return err
	}

	var bookmarkURL string
	for _, bookmark := range store.Bookmarks {
		if bookmark.ID == id {
			bookmarkURL = bookmark.URL
			break
		}
	}
	s.mu.Unlock()

	if bookmarkURL == "" {
		return errors.New("书签不存在")
	}

	return browser.OpenURL(bookmarkURL)
}

func (s *BookmarkService) StorePath() string {
	return s.filePath
}

func (s *BookmarkService) readStore() (bookmarkStore, error) {
	content, err := os.ReadFile(s.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return bookmarkStore{Bookmarks: []Bookmark{}}, nil
	}
	if err != nil {
		return bookmarkStore{}, fmt.Errorf("读取书签失败: %w", err)
	}
	if len(strings.TrimSpace(string(content))) == 0 {
		return bookmarkStore{Bookmarks: []Bookmark{}}, nil
	}

	var store bookmarkStore
	if err := json.Unmarshal(content, &store); err != nil {
		return bookmarkStore{}, fmt.Errorf("解析书签失败: %w", err)
	}
	if store.Bookmarks == nil {
		store.Bookmarks = []Bookmark{}
	}

	return store, nil
}

func (s *BookmarkService) writeStore(store bookmarkStore) error {
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o755); err != nil {
		return fmt.Errorf("创建书签目录失败: %w", err)
	}

	content, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化书签失败: %w", err)
	}

	tempPath := s.filePath + ".tmp"
	if err := os.WriteFile(tempPath, content, 0o644); err != nil {
		return fmt.Errorf("写入书签失败: %w", err)
	}

	if err := os.Rename(tempPath, s.filePath); err != nil {
		return fmt.Errorf("保存书签失败: %w", err)
	}

	return nil
}

func normalizeBookmarkInput(input BookmarkInput) (BookmarkInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.URL = strings.TrimSpace(input.URL)
	input.Note = strings.TrimSpace(input.Note)
	input.Tags = normalizeTags(input.Tags)

	if input.URL == "" {
		return BookmarkInput{}, errors.New("请输入网址")
	}
	if !strings.Contains(input.URL, "://") {
		input.URL = "https://" + input.URL
	}

	parsedURL, err := url.ParseRequestURI(input.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return BookmarkInput{}, errors.New("请输入有效网址")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return BookmarkInput{}, errors.New("仅支持 http 和 https 网址")
	}

	if input.Title == "" {
		input.Title = parsedURL.Host
	}

	return input, nil
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		for _, part := range strings.FieldsFunc(tag, func(r rune) bool {
			return r == ',' || r == '，' || r == ';' || r == '；'
		}) {
			normalized := strings.TrimSpace(part)
			if normalized == "" || seen[normalized] {
				continue
			}
			seen[normalized] = true
			result = append(result, normalized)
		}
	}

	return result
}

func sortBookmarks(bookmarks []Bookmark) {
	sort.SliceStable(bookmarks, func(i, j int) bool {
		return bookmarks[i].UpdatedAt > bookmarks[j].UpdatedAt
	})
}
