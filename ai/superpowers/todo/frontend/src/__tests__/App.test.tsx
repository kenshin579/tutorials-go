import { describe, it, expect } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from '../App'
import { server } from './mocks/server'
import { http, HttpResponse } from 'msw'

describe('App 통합 시나리오', () => {
  it('빈 상태 → 추가 → 토글 → 삭제', async () => {
    render(<App />)
    await waitFor(() => expect(screen.getByText('할 일이 없습니다.')).toBeDefined())

    // 추가
    await userEvent.type(screen.getByRole('textbox', { name: '제목' }), '우유 사기')
    await userEvent.click(screen.getByRole('button', { name: '추가' }))
    await waitFor(() => expect(screen.getByText('우유 사기')).toBeDefined())

    // 완료 토글
    await userEvent.click(screen.getByRole('checkbox', { name: /완료/ }))
    await waitFor(() => {
      const cb = screen.getByRole('checkbox', { name: /완료/ }) as HTMLInputElement
      expect(cb.checked).toBe(true)
    })

    // 삭제
    await userEvent.click(screen.getByRole('button', { name: '삭제' }))
    await waitFor(() => expect(screen.getByText('할 일이 없습니다.')).toBeDefined())
  })

  it('네트워크 실패 시 에러 배너 표시', async () => {
    server.use(http.get('/api/todos', () => HttpResponse.error()))
    render(<App />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeDefined())
  })

  it('헤더에 active/completed 카운트를 표시한다', async () => {
    render(<App />)
    // 빈 상태에서는 카운트 표시 생략
    await waitFor(() => expect(screen.getByText('할 일이 없습니다.')).toBeDefined())
    expect(screen.queryByText(/진행 중/)).toBeNull()

    // 1개 추가
    await userEvent.type(screen.getByRole('textbox', { name: '제목' }), 'A')
    await userEvent.click(screen.getByRole('button', { name: '추가' }))
    await waitFor(() => expect(screen.getByText(/1개 진행 중/)).toBeDefined())

    // 1개 추가 후 완료 토글
    await userEvent.type(screen.getByRole('textbox', { name: '제목' }), 'B')
    await userEvent.click(screen.getByRole('button', { name: '추가' }))
    await waitFor(() => expect(screen.getByText(/2개 진행 중/)).toBeDefined())
    const checkboxes = screen.getAllByRole('checkbox', { name: /완료/ })
    await userEvent.click(checkboxes[0])
    await waitFor(() => expect(screen.getByText(/1개 진행 중 · 1개 완료/)).toBeDefined())
  })
})
