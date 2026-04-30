import { test, expect, type Page, type APIRequestContext } from '@playwright/test'

const API = 'http://localhost:8080/api'

async function clearAllTodos(request: APIRequestContext) {
  const res = await request.get(`${API}/todos`)
  const todos = (await res.json()) as Array<{ id: string }>
  for (const t of todos) {
    await request.delete(`${API}/todos/${t.id}`)
  }
}

async function addTodo(page: Page, title: string, priority: 'low' | 'medium' | 'high' = 'medium') {
  const titleInput = page.getByLabel('제목')
  // Wait for previous submit to complete (input becomes empty and button re-enabled)
  await expect(titleInput).toBeEnabled()
  await titleInput.fill(title)
  await page.getByLabel('우선순위', { exact: true }).selectOption(priority)
  await page.getByRole('button', { name: '추가' }).click()
  // Wait for the item to appear in the list (confirms network round-trip completed)
  await expect(page.getByRole('list').getByText(title)).toBeVisible()
}

test.beforeEach(async ({ page, request }) => {
  await clearAllTodos(request)
  await page.goto('/')
})

test('1. 빈 상태 → 항목 추가 → 목록 등장', async ({ page }) => {
  await expect(page.getByText('할 일이 없습니다.')).toBeVisible()

  await addTodo(page, '우유 사기', 'high')

  await expect(page.getByText('우유 사기')).toBeVisible()
  await expect(page.getByRole('list').getByText('높음')).toBeVisible()
  await expect(page.getByText('할 일이 없습니다.')).not.toBeVisible()
})

test('2. 완료 토글 → 체크박스 상태 + BE 반영', async ({ page, request }) => {
  await addTodo(page, '운동하기')
  const checkbox = page.getByRole('checkbox', { name: /완료/ })
  await expect(checkbox).not.toBeChecked()

  await checkbox.click()
  await expect(checkbox).toBeChecked()

  // BE 검증
  const res = await request.get(`${API}/todos`)
  const todos = (await res.json()) as Array<{ title: string; completed: boolean }>
  expect(todos).toHaveLength(1)
  expect(todos[0].completed).toBe(true)
})

test('3. 필터 (전체/미완료/완료) 전환', async ({ page }) => {
  await addTodo(page, '활성 항목')
  await addTodo(page, '완료 항목')
  // '완료 항목'을 완료 처리 — li 기준으로 체크박스를 찾아서 클릭
  const completedItem = page.getByRole('list').locator('li').filter({ hasText: '완료 항목' })
  const completedCheckbox = completedItem.getByRole('checkbox')
  await completedCheckbox.click()
  await expect(completedCheckbox).toBeChecked()

  // 미완료 필터: "활성 항목"만
  await page.getByRole('radio', { name: '미완료', exact: true }).click()
  await expect(page.getByRole('list').getByText('활성 항목')).toBeVisible()
  await expect(page.getByRole('list').getByText('완료 항목')).not.toBeVisible()

  // 완료 필터: "완료 항목"만
  await page.getByRole('radio', { name: '완료', exact: true }).click()
  await expect(page.getByRole('list').getByText('완료 항목')).toBeVisible()
  await expect(page.getByRole('list').getByText('활성 항목')).not.toBeVisible()

  // 전체로 복귀
  await page.getByRole('radio', { name: '전체', exact: true }).click()
  await expect(page.getByRole('list').getByText('활성 항목')).toBeVisible()
  await expect(page.getByRole('list').getByText('완료 항목')).toBeVisible()
})

test('4. 인라인 편집: 제목 클릭 → 수정 → blur 저장', async ({ page, request }) => {
  await addTodo(page, '원본 제목')
  await page.getByText('원본 제목').click()

  const editor = page.getByLabel('제목 편집')
  await editor.fill('수정된 제목')
  await editor.blur()

  await expect(page.getByText('수정된 제목')).toBeVisible()
  await expect(page.getByText('원본 제목')).not.toBeVisible()

  // BE 검증
  const res = await request.get(`${API}/todos`)
  const todos = (await res.json()) as Array<{ title: string }>
  expect(todos[0].title).toBe('수정된 제목')
})

test('5. 인라인 편집 Escape 시 원복', async ({ page }) => {
  await addTodo(page, '원본')
  await page.getByText('원본').click()
  const editor = page.getByLabel('제목 편집')
  await editor.fill('취소될 값')
  await editor.press('Escape')
  await expect(page.getByText('원본')).toBeVisible()
  await expect(page.getByText('취소될 값')).not.toBeVisible()
})

test('6. 정렬 토글 (asc ↔ desc)', async ({ page }) => {
  await addTodo(page, '첫 번째')
  await addTodo(page, '두 번째')
  // 기본은 createdAt desc → 두 번째가 먼저
  await expect(page.locator('li').first()).toContainText('두 번째')
  await expect(page.locator('li').nth(1)).toContainText('첫 번째')

  // asc로 토글
  await page.getByRole('button', { name: '정렬 방향 토글' }).click()
  // Wait for order to flip before reading
  await expect(page.locator('li').first()).toContainText('첫 번째')
  await expect(page.locator('li').nth(1)).toContainText('두 번째')
})

test('7. 삭제 → 빈 상태 복귀', async ({ page }) => {
  await addTodo(page, '삭제 예정')
  await expect(page.getByText('삭제 예정')).toBeVisible()

  await page.getByRole('button', { name: '삭제' }).click()

  await expect(page.getByText('삭제 예정')).not.toBeVisible()
  await expect(page.getByText('할 일이 없습니다.')).toBeVisible()
})

test('8. API 실패 시 에러 배너 표시', async ({ page }) => {
  // route 가로채기로 BE 응답을 강제 실패
  await page.route('**/api/todos*', (route) => {
    if (route.request().method() === 'GET') {
      return route.abort('failed')
    }
    return route.continue()
  })
  await page.reload()
  await expect(page.getByRole('alert')).toBeVisible()
})

test('9. priority 정렬 (high → low)', async ({ page }) => {
  await addTodo(page, '낮음 항목', 'low')
  await addTodo(page, '높음 항목', 'high')
  await addTodo(page, '보통 항목', 'medium')

  // priority 정렬로 변경 (sort select is inside the toolbar)
  const sortSelect = page.getByRole('toolbar').locator('select')
  await sortSelect.selectOption('priority')
  // 현재 desc이므로 high가 먼저
  await expect(page.locator('li').first()).toContainText('높음 항목')

  // asc로 토글 → low가 먼저
  await page.getByRole('button', { name: '정렬 방향 토글' }).click()
  await expect(page.locator('li').first()).toContainText('낮음 항목')
  await expect(page.locator('li').nth(2)).toContainText('높음 항목')
})
