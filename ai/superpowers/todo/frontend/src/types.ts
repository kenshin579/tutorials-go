// 본 파일은 backend/todo/todo.go와 수동 동기화한다.
// BE 모델 변경 시 같이 수정 (CLAUDE.md 정책).

export type Priority = 'low' | 'medium' | 'high'

export type Status = 'all' | 'active' | 'completed'

export type SortKey = 'createdAt' | 'dueDate' | 'priority'

export type Order = 'asc' | 'desc'

export interface Todo {
  id: string
  title: string
  completed: boolean
  priority: Priority
  dueDate: string | null // RFC3339 또는 null
  createdAt: string
  updatedAt: string
}

export interface NewTodo {
  title: string
  priority?: Priority
  dueDate?: string | null
}

export interface Patch {
  title?: string
  completed?: boolean
  priority?: Priority
  dueDate?: string | null // null이면 명시적 클리어
}

export interface Query {
  status: Status
  sort: SortKey
  order: Order
}
