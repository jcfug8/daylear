// Centralized API base URL for backend requests
//
// For local development, set VITE_API_BASE_URL in .env to http://localhost:8080/
// For production, set VITE_API_BASE_URL to https://daylear.com/
//
// Usage: import { API_BASE_URL } from '@/constants/api'

const raw = import.meta.env.VITE_API_BASE_URL as string | undefined
if (!raw) throw new Error('VITE_API_BASE_URL is not set')
export const API_BASE_URL = raw.endsWith('/') ? raw : raw + '/' 