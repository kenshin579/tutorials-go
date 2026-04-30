import { useCallback, useEffect, useReducer, useRef } from 'react'
import { api, ApiError } from '../api'
import type { NewTodo, Patch, Query, Todo } from '../types'

type State = {
  todos: Todo[]
  loading: boolean
  error: string | null
}

type Action =
  | { type: 'fetch_start' }
  | { type: 'fetch_success'; todos: Todo[] }
  | { type: 'fetch_error'; error: string }

const initial: State = { todos: [], loading: false, error: null }

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case 'fetch_start':
      return { ...state, loading: true, error: null }
    case 'fetch_success':
      return { todos: action.todos, loading: false, error: null }
    case 'fetch_error':
      return { ...state, loading: false, error: action.error }
  }
}

export function useTodos(query: Query) {
  const [state, dispatch] = useReducer(reducer, initial)
  const queryRef = useRef(query)
  queryRef.current = query

  const refetch = useCallback(async () => {
    dispatch({ type: 'fetch_start' })
    try {
      const todos = await api.list(queryRef.current)
      dispatch({ type: 'fetch_success', todos })
    } catch (e) {
      const msg =
        e instanceof ApiError
          ? e.message
          : e instanceof Error
            ? e.message
            : '서버에 연결할 수 없습니다'
      dispatch({ type: 'fetch_error', error: msg })
    }
  }, [])

  useEffect(() => {
    void refetch()
  }, [query.status, query.sort, query.order, refetch])

  const create = useCallback(
    async (input: NewTodo) => {
      await api.create(input)
      await refetch()
    },
    [refetch],
  )

  const update = useCallback(
    async (id: string, patch: Patch) => {
      await api.update(id, patch)
      await refetch()
    },
    [refetch],
  )

  const remove = useCallback(
    async (id: string) => {
      await api.remove(id)
      await refetch()
    },
    [refetch],
  )

  return { ...state, create, update, remove, refetch }
}
