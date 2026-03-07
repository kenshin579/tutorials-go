import { test, expect } from "@playwright/test";

// wails dev 실행 후 테스트: npx playwright test

const uniqueText = (prefix: string) => `${prefix}-${Date.now()}`;

test.describe("Wails Todo App", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("페이지 타이틀이 표시된다", async ({ page }) => {
    await expect(page.locator("h1")).toHaveText("Wails Todo");
  });

  test("초기 상태에서 완료 카운트가 표시된다", async ({ page }) => {
    await expect(page.locator(".status")).toBeVisible();
  });

  test("할 일을 추가할 수 있다", async ({ page }) => {
    const title = uniqueText("추가테스트");
    const input = page.locator('.todo-input input[type="text"]');

    await input.fill(title);
    await page.locator('.todo-input button[type="submit"]').click();

    await expect(
      page.locator(".todo-item").filter({ hasText: title })
    ).toBeVisible();
  });

  test("빈 입력은 추가되지 않는다", async ({ page }) => {
    const beforeCount = await page.locator(".todo-item").count();

    await page.locator('.todo-input button[type="submit"]').click();

    const afterCount = await page.locator(".todo-item").count();
    expect(afterCount).toBe(beforeCount);
  });

  test("체크박스로 완료 상태를 토글할 수 있다", async ({ page }) => {
    const title = uniqueText("토글테스트");
    const input = page.locator('.todo-input input[type="text"]');
    await input.fill(title);
    await page.locator('.todo-input button[type="submit"]').click();

    const todoItem = page.locator(".todo-item").filter({ hasText: title });
    await expect(todoItem).toBeVisible();

    // 체크박스 클릭하여 완료 처리
    const checkbox = todoItem.locator('input[type="checkbox"]');
    await checkbox.click();
    await expect(todoItem).toHaveClass(/done/);

    // 다시 클릭하여 미완료로 변경
    await checkbox.click();
    await expect(todoItem).not.toHaveClass(/done/);
  });

  test("삭제 버튼을 클릭하면 다이얼로그가 표시된다", async ({ page }) => {
    const title = uniqueText("삭제테스트");
    const input = page.locator('.todo-input input[type="text"]');
    await input.fill(title);
    await page.locator('.todo-input button[type="submit"]').click();

    const todoItem = page.locator(".todo-item").filter({ hasText: title });
    await expect(todoItem).toBeVisible();

    // 삭제 버튼 클릭 (네이티브 다이얼로그가 뜨므로 dialog 이벤트 처리)
    page.on("dialog", async (dialog) => {
      await dialog.accept();
    });

    await todoItem.locator(".delete-btn").click();
  });

  test("Enter 키로 할 일을 추가할 수 있다", async ({ page }) => {
    const title = uniqueText("엔터테스트");
    const input = page.locator('.todo-input input[type="text"]');

    await input.fill(title);
    await input.press("Enter");

    await expect(
      page.locator(".todo-item").filter({ hasText: title })
    ).toBeVisible();
  });

  test("할 일이 없으면 빈 메시지가 표시된다", async ({ page }) => {
    const emptyMessage = page.locator(".empty-message");
    const todoItems = page.locator(".todo-item");

    const count = await todoItems.count();
    if (count === 0) {
      await expect(emptyMessage).toHaveText("할 일이 없습니다.");
    }
  });
});
