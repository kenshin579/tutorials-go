import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TodoItem } from '../components/TodoItem'
import type { Todo } from '../types'

function makeTodo(overrides: Partial<Todo> = {}): Todo {
  return {
    id: 'id1',
    title: '할일 1',
    completed: false,
    priority: 'high',
    dueDate: null,
    createdAt: '2026-04-30T10:00:00Z',
    updatedAt: '2026-04-30T10:00:00Z',
    ...overrides,
  }
}

describe('TodoItem', () => {
  it('priority 뱃지 표시', () => {
    render(<TodoItem todo={makeTodo({ priority: 'high' })} onUpdate={vi.fn()} onRemove={vi.fn()} />)
    expect(screen.getByText('높음')).toBeDefined()
  })

  it('dueDate 표시 (있을 때만)', () => {
    const { rerender } = render(
      <TodoItem todo={makeTodo({ dueDate: null })} onUpdate={vi.fn()} onRemove={vi.fn()} />,
    )
    expect(screen.queryByText(/마감/)).toBeNull()

    rerender(
      <TodoItem
        todo={makeTodo({ dueDate: '2026-05-15T18:00:00Z' })}
        onUpdate={vi.fn()}
        onRemove={vi.fn()}
      />,
    )
    expect(screen.getByText(/마감/)).toBeDefined()
  })

  it('체크박스 클릭 시 onUpdate 호출', async () => {
    const onUpdate = vi.fn().mockResolvedValue(undefined)
    render(<TodoItem todo={makeTodo()} onUpdate={onUpdate} onRemove={vi.fn()} />)
    const cb = screen.getByRole('checkbox', { name: /완료/ })
    await userEvent.click(cb)
    expect(onUpdate).toHaveBeenCalledWith('id1', { completed: true })
  })

  it('삭제 버튼 클릭 시 onRemove 호출', async () => {
    const onRemove = vi.fn().mockResolvedValue(undefined)
    render(<TodoItem todo={makeTodo()} onUpdate={vi.fn()} onRemove={onRemove} />)
    await userEvent.click(screen.getByRole('button', { name: /삭제/ }))
    expect(onRemove).toHaveBeenCalledWith('id1')
  })

  it('제목 클릭 → 편집 → blur 시 onUpdate 호출', async () => {
    const onUpdate = vi.fn().mockResolvedValue(undefined)
    render(<TodoItem todo={makeTodo()} onUpdate={onUpdate} onRemove={vi.fn()} />)
    const title = screen.getByText('할일 1')
    await userEvent.click(title)
    const input = screen.getByRole('textbox', { name: /제목/ })
    await userEvent.clear(input)
    await userEvent.type(input, '수정됨')
    input.blur()
    expect(onUpdate).toHaveBeenCalledWith('id1', { title: '수정됨' })
  })
})
