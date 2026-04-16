import { DevToolsService } from '../../bindings/wails-mytools/internal/services'
import type { DevEnvironment } from '../../bindings/wails-mytools/internal/services'

export type { DevEnvironment }

export function openDevTools() {
  return DevToolsService.OpenDevTools()
}

export function getDevEnvironment() {
  return DevToolsService.Environment()
}
