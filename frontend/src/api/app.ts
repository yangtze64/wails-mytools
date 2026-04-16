import { AppService } from '../../bindings/wails-mytools/internal/services'
import type { AppInfo } from '../types/app'

export async function getAppInfo(): Promise<AppInfo> {
  return AppService.Info()
}

export async function getHealthStatus(): Promise<string> {
  return AppService.Health()
}
