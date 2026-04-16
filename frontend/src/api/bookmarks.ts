import { BookmarkService } from '../../bindings/wails-mytools/internal/services'
import type { Bookmark, BookmarkInput } from '../../bindings/wails-mytools/internal/services'

export type { Bookmark, BookmarkInput }

export function listBookmarks() {
  return BookmarkService.List()
}

export function createBookmark(input: BookmarkInput) {
  return BookmarkService.Create(input)
}

export function updateBookmark(id: string, input: BookmarkInput) {
  return BookmarkService.Update(id, input)
}

export function deleteBookmark(id: string) {
  return BookmarkService.Delete(id)
}

export function openBookmark(id: string) {
  return BookmarkService.Open(id)
}

export function getBookmarkStorePath() {
  return BookmarkService.StorePath()
}
