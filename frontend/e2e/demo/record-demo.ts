import { mkdir, rename, rm } from 'node:fs/promises'
import path from 'node:path'
import { chromium, expect } from '@playwright/test'

const baseURL = process.env.DEMO_BASE_URL || 'http://127.0.0.1:3100'
const outputDir = path.resolve(process.cwd(), '../docs/assets/demo/raw')
const outputPath = path.join(outputDir, 'web-demo.webm')
const pause = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

async function main() {
  await mkdir(outputDir, { recursive: true })
  await rm(outputPath, { force: true })
  const browser = await chromium.launch({ channel: 'chrome', headless: true })
  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
    recordVideo: { dir: outputDir, size: { width: 1920, height: 1080 } },
    colorScheme: 'light', locale: 'zh-CN', timezoneId: 'Asia/Shanghai', reducedMotion: 'reduce',
  })
  const page = await context.newPage()
  await page.addInitScript(() => { localStorage.setItem('theme', 'light'); document.documentElement.dataset.demoCapture = 'true' })
  const video = page.video()
  await page.goto(`${baseURL}/`, { waitUntil: 'networkidle' }); await expect(page.getByTestId('hero')).toBeVisible(); await pause(4000)
  await page.getByTestId('start-experience').click(); await expect(page.getByTestId('demo-dashboard')).toBeVisible(); await pause(3000)
  await page.getByTestId('open-attention').click(); await expect(page.getByTestId('demo-opportunity-detail')).toBeVisible(); await pause(3500)
  await page.getByTestId('generate-demo-draft').click(); await expect(page.getByTestId('demo-draft')).not.toHaveValue(''); await pause(3500)
  await page.goto(`${baseURL}/demo`, { waitUntil: 'networkidle' }); await page.getByTestId('filter-telegram').click(); await page.getByTestId('advanced-filter').click(); await expect(page.getByTestId('advanced-filter-panel')).toBeVisible(); await pause(3500)
  await page.goto(`${baseURL}/demo/settings/subscription`, { waitUntil: 'networkidle' }); await expect(page.getByTestId('demo-settings-subscription')).toBeVisible(); await pause(2500)
  await page.goto(`${baseURL}/demo/settings/telegram`, { waitUntil: 'networkidle' }); await expect(page.getByTestId('demo-settings-telegram')).toBeVisible(); await pause(2500)
  await page.goto(`${baseURL}/demo/settings/working-hours`, { waitUntil: 'networkidle' }); await expect(page.getByTestId('demo-settings-working-hours')).toBeVisible(); await pause(2500)
  await page.goto(`${baseURL}/#apps`, { waitUntil: 'networkidle' }); await expect(page.getByTestId('multi-platform')).toBeVisible(); await page.getByTestId('multi-platform').scrollIntoViewIfNeeded(); await pause(3500)
  await context.close()
  if (!video) throw new Error('Playwright did not create a video recorder')
  const generated = await video.path()
  await browser.close()
  await rename(generated, outputPath)
  process.stdout.write(`${outputPath}\n`)
}

main().catch((error) => { console.error(error); process.exit(1) })
