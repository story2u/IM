import { mkdir } from 'node:fs/promises'
import path from 'node:path'
import { expect, test } from '@playwright/test'
import { prepareDemoPage, setTheme } from './fixtures'

const output = path.resolve(process.cwd(), '../docs/assets/screenshots/web')
const demoOutput = path.resolve(process.cwd(), '../docs/assets/demo')

test.beforeAll(async () => { await mkdir(output, { recursive: true }); await mkdir(demoOutput, { recursive: true }) })

test('desktop screenshots', async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 1000 })
  await prepareDemoPage(page, '/', 'product-home')
  await page.screenshot({ path: path.join(output, 'web-home-desktop.png'), fullPage: true })

  await prepareDemoPage(page, '/demo', 'demo-dashboard')
  await page.screenshot({ path: path.join(output, 'dashboard-desktop.png'), fullPage: true })
  await page.getByTestId('attention-alert').screenshot({ path: path.join(output, 'dashboard-attention.png') })
  await page.getByTestId('advanced-filter').click()
  await expect(page.getByTestId('advanced-filter-panel')).toBeVisible()
  await page.screenshot({ path: path.join(output, 'dashboard-filter.png'), fullPage: false })

  await prepareDemoPage(page, '/demo/opportunity/demo-procurement-50', 'demo-opportunity-detail')
  await page.screenshot({ path: path.join(output, 'opportunity-detail.png'), fullPage: true })
  await prepareDemoPage(page, '/demo/settings', 'demo-settings-overview')
  await page.screenshot({ path: path.join(output, 'settings-center.png'), fullPage: true })
  await prepareDemoPage(page, '/demo/settings/telegram', 'demo-settings-telegram')
  await page.screenshot({ path: path.join(output, 'telegram-connections.png'), fullPage: true })
  await prepareDemoPage(page, '/demo/settings/subscription', 'demo-settings-subscription')
  await page.screenshot({ path: path.join(output, 'subscription.png'), fullPage: true })
})

test('mobile and dark screenshots', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 })
  await prepareDemoPage(page, '/', 'product-home')
  await page.screenshot({ path: path.join(output, 'web-home-mobile.png'), fullPage: true })
  await prepareDemoPage(page, '/demo', 'demo-dashboard')
  await page.screenshot({ path: path.join(output, 'dashboard-mobile.png'), fullPage: true })
  await prepareDemoPage(page, '/demo/opportunity/demo-procurement-50', 'demo-opportunity-detail')
  await page.screenshot({ path: path.join(output, 'opportunity-detail-mobile.png'), fullPage: true })
  await prepareDemoPage(page, '/demo/settings', 'demo-settings-overview')
  await page.screenshot({ path: path.join(output, 'settings-mobile.png'), fullPage: true })

  await page.setViewportSize({ width: 1440, height: 1000 })
  await prepareDemoPage(page, '/', 'product-home'); await setTheme(page, 'dark')
  await page.screenshot({ path: path.join(output, 'web-home-dark.png'), fullPage: true })
  await prepareDemoPage(page, '/demo', 'demo-dashboard'); await setTheme(page, 'dark')
  await page.screenshot({ path: path.join(output, 'dashboard-dark.png'), fullPage: true })
})

test('social preview and video cover sources', async ({ page }) => {
  await page.setViewportSize({ width: 1280, height: 640 })
  await prepareDemoPage(page, '/', 'product-home')
  await page.screenshot({ path: path.join(output, 'github-social-preview.png') })
  await page.setViewportSize({ width: 1920, height: 1080 })
  await page.goto('/demo/cover', { waitUntil: 'networkidle' })
  await expect(page.getByTestId('demo-cover')).toBeVisible()
  await page.screenshot({ path: path.join(output, 'demo-video-cover.png') })
})
