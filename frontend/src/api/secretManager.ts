import { SecretManagerService } from '../../bindings/wails-mytools/internal/services'
import type { SecretInput, SecretItem } from '../../bindings/wails-mytools/internal/services'

export type { SecretInput, SecretItem }

export function listSecrets() {
  return SecretManagerService.List()
}

export function createSecret(input: SecretInput) {
  return SecretManagerService.Create(input)
}

export function updateSecret(id: string, input: SecretInput) {
  return SecretManagerService.Update(id, input)
}

export function deleteSecret(id: string) {
  return SecretManagerService.Delete(id)
}

export function revealSecret(id: string) {
  return SecretManagerService.RevealSecret(id)
}

export function revealSecretExtra(id: string) {
  return SecretManagerService.RevealSecretExtra(id)
}

export function getSecretStorePath() {
  return SecretManagerService.StorePath()
}
