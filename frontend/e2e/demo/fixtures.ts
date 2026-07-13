import { expect, type Page } from '@playwright/test'

export async function prepareDemoPage(page: Page, path: string, testId: string) {
  await page.addInitScript(() => {
    localStorage.setItem('theme', 'light')
    document.documentElement.dataset.demoCapture = 'true'
  })
  await page.goto(path, { waitUntil: 'networkidle' })
  await page.evaluate(() => document.fonts.ready)
  await expect(page.getByTestId(testId)).toBeVisible()
}

export async function setTheme(page: Page, theme: 'light' | 'dark') {
  await page.evaluate((nextTheme) => {
    localStorage.setItem('theme', nextTheme)
    document.documentElement.classList.toggle('dark', nextTheme === 'dark')
    document.documentElement.style.colorScheme = nextTheme
  }, theme)
}
