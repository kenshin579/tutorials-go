import { useState, type FormEvent } from 'react'
import type { NewTodo, Priority } from '../types'

interface Props {
  onCreate: (input: NewTodo) => Promise<void>
}

export function TodoForm({ onCreate }: Props) {
  const [title, setTitle] = useState('')
  const [priority, setPriority] = useState<Priority>('medium')
  const [dueDate, setDueDate] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    if (!title.trim() || submitting) return
    setSubmitting(true)
    try {
      await onCreate({
        title: title.trim(),
        priority,
        dueDate: dueDate ? new Date(dueDate).toISOString() : undefined,
      })
      setTitle('')
      setDueDate('')
      setPriority('medium')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form className="todo-form" onSubmit={submit} aria-label="새 할일 추가">
      <input
        className="todo-form__input"
        aria-label="제목"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="할 일을 입력하세요…"
        maxLength={200}
      />
      <select
        className="todo-form__priority"
        aria-label="우선순위"
        value={priority}
        onChange={(e) => setPriority(e.target.value as Priority)}
      >
        <option value="low">낮음</option>
        <option value="medium">보통</option>
        <option value="high">높음</option>
      </select>
      <input
        className="todo-form__due"
        aria-label="마감일"
        type="datetime-local"
        value={dueDate}
        onChange={(e) => setDueDate(e.target.value)}
      />
      <button
        className="todo-form__submit"
        type="submit"
        disabled={!title.trim() || submitting}
      >
        추가
      </button>
    </form>
  )
}
