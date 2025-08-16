import type { PermissionLevel } from '@/genapi/api/types'

// Permission level hierarchy mapping based on protobuf values
const PERMISSION_LEVEL_VALUES: Record<PermissionLevel, number> = {
  'PERMISSION_LEVEL_UNSPECIFIED': 0,
  'PERMISSION_LEVEL_PUBLIC': 1,
  'PERMISSION_LEVEL_READ': 100,
  'PERMISSION_LEVEL_WRITE': 200,
  'PERMISSION_LEVEL_ADMIN': 300,
}

/**
 * Get the numerical value of a permission level for comparison
 */
export function getPermissionValue(level: PermissionLevel | undefined): number {
  if (!level) return 0
  return PERMISSION_LEVEL_VALUES[level] || 0
}

/**
 * Check if user has at least the required permission level
 * @param userLevel - The user's current permission level
 * @param requiredLevel - The minimum required permission level
 * @returns true if user has sufficient permissions
 */
export function hasPermissionLevel(
  userLevel: PermissionLevel | undefined, 
  requiredLevel: PermissionLevel
): boolean {
  return getPermissionValue(userLevel) >= getPermissionValue(requiredLevel)
}

/**
 * Check if user has write permissions (WRITE or ADMIN)
 */
export function hasWritePermission(level: PermissionLevel | undefined): boolean {
  return hasPermissionLevel(level, 'PERMISSION_LEVEL_WRITE')
}

/**
 * Check if user has write permissions (WRITE but not ADMIN)
 */
export function hasWriteOnlyPermission(level: PermissionLevel | undefined): boolean {
  return hasPermissionLevel(level, 'PERMISSION_LEVEL_WRITE') && 
         !hasPermissionLevel(level, 'PERMISSION_LEVEL_ADMIN')
}

/**
 * Check if user has read permissions (READ, WRITE, or ADMIN)
 */
export function hasReadPermission(level: PermissionLevel | undefined): boolean {
  return hasPermissionLevel(level, 'PERMISSION_LEVEL_READ')
}

/**
 * Check if user has admin permissions
 */
export function hasAdminPermission(level: PermissionLevel | undefined): boolean {
  return hasPermissionLevel(level, 'PERMISSION_LEVEL_ADMIN')
} 