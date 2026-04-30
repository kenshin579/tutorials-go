import type { Patch, Todo } from '../types'
import { TodoItem } from './TodoItem'

interface Props {
  todos: Todo[]
  onUpdate: (id: string, patch: Patch) => Promise<void>
  onRemove: (id: string) => Promise<void>
}

export function TodoList({ todos, onUpdate, onRemove }: Props) {
  if (todos.length === 0) {
    return <p className="todo-list todo-list--empty">할 일이 없습니다.</p>
  }
  return (
    <ul className="todo-list" aria-label="할 일 목록">
      {todos.map((t) => (
        <TodoItem key={t.id} todo={t} onUpdate={onUpdate} onRemove={onRemove} />
      ))}
    </ul>
  )
}
