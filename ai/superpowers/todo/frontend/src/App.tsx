import { useState } from 'react'
import { useTodos } from './hooks/useTodos'
import { TodoForm } from './components/TodoForm'
import { FilterBar } from './components/FilterBar'
import { TodoList } from './components/TodoList'
import type { Query } from './types'

const defaultQuery: Query = { status: 'all', sort: 'createdAt', order: 'desc' }

export default function App() {
  const [query, setQuery] = useState<Query>(defaultQuery)
  const { todos, error, create, update, remove } = useTodos(query)

  return (
    <main>
      <h1>Superpowers Todo</h1>
      {error && <div role="alert">에러: {error}</div>}
      <TodoForm onCreate={create} />
      <FilterBar query={query} onChange={setQuery} />
      <TodoList todos={todos} onUpdate={update} onRemove={remove} />
    </main>
  )
}
