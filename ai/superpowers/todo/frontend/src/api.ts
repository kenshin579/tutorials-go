import type { NewTodo, Patch, Query, Todo } from './types'

const BASE = '/api'

export class ApiError extends Error {
  constructor(public code: string, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

async function call<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
  })
  if (!res.ok) {
    let code = 'unknown'
    let message = res.statusText
    try {
      const body = (await res.json()) as { error?: { code?: string; message?: string } }
      code = body.error?.code ?? code
      message = body.error?.message ?? message
    } catch {
      // body가 JSON이 아닐 수 있음. statusText 그대로 사용.
    }
    throw new ApiError(code, message)
  }
  if (res.status === 204) return undefined as T
  return (await res.json()) as T
}

export const api = {
  list: (q: Query) =>
    call<Todo[]>(`/todos?${new URLSearchParams(q as unknown as Record<string, string>)}`),
  create: (input: NewTodo) =>
    call<Todo>('/todos', { method: 'POST', body: JSON.stringify(input) }),
  update: (id: string, patch: Patch) =>
    call<Todo>(`/todos/${id}`, { method: 'PATCH', body: JSON.stringify(patch) }),
  remove: (id: string) =>
    call<void>(`/todos/${id}`, { method: 'DELETE' }),
}
