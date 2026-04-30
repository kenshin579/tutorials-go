import { useMemo, useState } from 'react'
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

  return (
    <main className="app">
      <header className="app-header">
        <h1 className="app-title">Todo</h1>
        {statusText && <span className="app-status">{statusText}</span>}
      </header>
      {error && <div role="alert">에러: {error}</div>}
      <TodoForm onCreate={create} />
      <FilterBar query={query} onChange={setQuery} />
      <TodoList todos={todos} onUpdate={update} onRemove={remove} />
    </main>
  )
}
