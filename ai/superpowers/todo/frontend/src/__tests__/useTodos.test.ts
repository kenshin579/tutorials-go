import { describe, it, expect } from 'vitest'
import { renderHook, waitFor, act } from '@testing-library/react'
import { useTodos } from '../hooks/useTodos'
import { resetMockTodos } from './mocks/handlers'
import type { Query } from '../types'

const defaultQuery: Query = { status: 'all', sort: 'createdAt', order: 'desc' }

describe('useTodos', () => {
  it('초기 fetch 후 todos 채워짐', async () => {
    resetMockTodos([
      {
        id: '1',
        title: 'seeded',
        completed: false,
        priority: 'medium',
        dueDate: null,
        createdAt: '2026-04-01T00:00:00Z',
        updatedAt: '2026-04-01T00:00:00Z',
      },
    ])
    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.todos).toHaveLength(1)
    expect(result.current.todos[0].title).toBe('seeded')
  })

  it('create 후 list refetch', async () => {
    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.todos).toHaveLength(0)

    await act(async () => {
      await result.current.create({ title: '새 할일' })
    })
    await waitFor(() => expect(result.current.todos).toHaveLength(1))
    expect(result.current.todos[0].title).toBe('새 할일')
  })

  it('update 후 list refetch', async () => {
    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.loading).toBe(false))
    await act(async () => {
      await result.current.create({ title: 'x' })
    })
    await waitFor(() => expect(result.current.todos).toHaveLength(1))
    const id = result.current.todos[0].id

    await act(async () => {
      await result.current.update(id, { completed: true })
    })
    await waitFor(() => expect(result.current.todos[0].completed).toBe(true))
  })

  it('remove 후 목록에서 제거', async () => {
    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.loading).toBe(false))
    await act(async () => {
      await result.current.create({ title: 'x' })
    })
    await waitFor(() => expect(result.current.todos).toHaveLength(1))
    const id = result.current.todos[0].id

    await act(async () => {
      await result.current.remove(id)
    })
    await waitFor(() => expect(result.current.todos).toHaveLength(0))
  })

  it('네트워크 실패 시 error 상태', async () => {
    const { server } = await import('./mocks/server')
    const { http, HttpResponse } = await import('msw')
    server.use(http.get('/api/todos', () => HttpResponse.error()))

    const { result } = renderHook(() => useTodos(defaultQuery))
    await waitFor(() => expect(result.current.error).not.toBeNull())
    expect(result.current.loading).toBe(false)
  })
})
