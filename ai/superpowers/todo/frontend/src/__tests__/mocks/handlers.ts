import { http, HttpResponse } from 'msw'
import type { Todo, NewTodo, Patch } from '../../types'

let todos: Todo[] = []

function uid(): string {
  return Math.random().toString(36).slice(2) + Date.now().toString(36)
}

export function resetMockTodos(seed: Todo[] = []) {
  todos = seed.map((t) => ({ ...t }))
}

export const handlers = [
  http.get('/api/health', () => HttpResponse.json({ status: 'ok' })),

  http.get('/api/todos', ({ request }) => {
    const url = new URL(request.url)
    const status = url.searchParams.get('status') ?? 'all'
    let out = [...todos]
    if (status === 'active') out = out.filter((t) => !t.completed)
    if (status === 'completed') out = out.filter((t) => t.completed)
    return HttpResponse.json(out)
  }),

  http.post('/api/todos', async ({ request }) => {
    const body = (await request.json()) as NewTodo
    if (!body.title?.trim()) {
      return HttpResponse.json(
        { error: { code: 'validation_failed', message: 'title is required', details: { field: 'title' } } },
        { status: 400 },
      )
    }
    const now = new Date().toISOString()
    const created: Todo = {
      id: uid(),
      title: body.title.trim(),
      completed: false,
      priority: body.priority ?? 'medium',
      dueDate: body.dueDate ?? null,
      createdAt: now,
      updatedAt: now,
    }
    todos.push(created)
    return HttpResponse.json(created, { status: 201 })
  }),

  http.patch('/api/todos/:id', async ({ params, request }) => {
    const id = params.id as string
    const idx = todos.findIndex((t) => t.id === id)
    if (idx === -1) {
      return HttpResponse.json(
        { error: { code: 'not_found', message: 'todo not found' } },
        { status: 404 },
      )
    }
    const patch = (await request.json()) as Patch
    const t = { ...todos[idx], ...patch, updatedAt: new Date().toISOString() }
    if ('dueDate' in patch && patch.dueDate === null) t.dueDate = null
    todos[idx] = t
    return HttpResponse.json(t)
  }),

  http.delete('/api/todos/:id', ({ params }) => {
    const id = params.id as string
    const before = todos.length
    todos = todos.filter((t) => t.id !== id)
    if (todos.length === before) {
      return HttpResponse.json(
        { error: { code: 'not_found', message: 'todo not found' } },
        { status: 404 },
      )
    }
    return new HttpResponse(null, { status: 204 })
  }),
]
