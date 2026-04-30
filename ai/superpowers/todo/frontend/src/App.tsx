import { useEffect, useMemo, useState } from 'react'
import { useTodos } from './hooks/useTodos'
import { TodoForm } from './components/TodoForm'
import { FilterBar } from './components/FilterBar'
import { TodoList } from './components/TodoList'
import type { Query } from './types'

const defaultQuery: Query = { status: 'all', sort: 'createdAt', order: 'desc' }

export default function App() {
  const [query, setQuery] = useState<Query>(defaultQuery)
  const { todos, error, create, update, remove } = useTodos(query)

  const counts = useMemo(() => {
    const active = todos.filter((t) => !t.completed).length
    const completed = todos.length - active
    return { all: todos.length, active, completed }
  }, [todos])

  const statusText = (() => {
    if (counts.all === 0) return null
    if (counts.completed === 0) return `${counts.active}개 진행 중`
    if (counts.active === 0) return `${counts.completed}개 완료`
    return `${counts.active}개 진행 중 · ${counts.completed}개 완료`
  })()

  const [dismissedError, setDismissedError] = useState(false)
  useEffect(() => {
    setDismissedError(false)
  }, [error])
  const showError = error !== null && !dismissedError

  return (
    <main className="app">
      <header className="app-header">
        <h1 className="app-title">Todo</h1>
        {statusText && <span className="app-status">{statusText}</span>}
      </header>
      {showError && (
        <div className="error-banner" role="alert">
          <span className="error-banner__label">에러:</span>
          <span className="error-banner__message">{error}</span>
          <button
            type="button"
            className="error-banner__close"
            aria-label="에러 닫기"
            onClick={() => setDismissedError(true)}
          >
            ×
          </button>
        </div>
      )}
      <TodoForm onCreate={create} />
      <FilterBar query={query} onChange={setQuery} />
      <TodoList todos={todos} onUpdate={update} onRemove={remove} />
    </main>
  )
}
