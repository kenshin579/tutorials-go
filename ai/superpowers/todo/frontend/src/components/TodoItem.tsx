import { useState, type KeyboardEvent } from 'react'
import type { Patch, Priority, Todo } from '../types'

interface Props {
  todo: Todo
  onUpdate: (id: string, patch: Patch) => Promise<void>
  onRemove: (id: string) => Promise<void>
}

const priorityLabel: Record<Priority, string> = { low: '낮음', medium: '보통', high: '높음' }

export function TodoItem({ todo, onUpdate, onRemove }: Props) {
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState(todo.title)

  const commit = () => {
    setEditing(false)
    const next = draft.trim()
    if (next && next !== todo.title) {
      void onUpdate(todo.id, { title: next })
    } else {
      setDraft(todo.title)
    }
  }

  const onKey = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.currentTarget.blur()
    } else if (e.key === 'Escape') {
      setDraft(todo.title)
      setEditing(false)
    }
  }

  return (
    <li>
      <input
        type="checkbox"
        aria-label="완료"
        checked={todo.completed}
        onChange={(e) => onUpdate(todo.id, { completed: e.target.checked })}
      />
      {editing ? (
        <input
          aria-label="제목 편집"
          value={draft}
          autoFocus
          onChange={(e) => setDraft(e.target.value)}
          onBlur={commit}
          onKeyDown={onKey}
        />
      ) : (
        <span onClick={() => setEditing(true)}>{todo.title}</span>
      )}
      <span data-priority={todo.priority}>{priorityLabel[todo.priority]}</span>
      {todo.dueDate && <span>마감: {new Date(todo.dueDate).toLocaleString()}</span>}
      <button type="button" aria-label="삭제" onClick={() => onRemove(todo.id)}>×</button>
    </li>
  )
}
