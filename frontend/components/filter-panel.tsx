'use client'

import { RotateCcw, SlidersHorizontal } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import {
  countActiveAdvancedFilters,
  defaultFilters,
  type DashboardFilters,
  type TimeRange,
} from '@/lib/dashboard-filters'
import { sopStageConfig, sopStageOrder, trustLevelConfig, type TrustLevel } from '@/lib/sop'
import type { SopStage, SourceType } from '@/lib/types'
import { cn } from '@/lib/utils'

const timeRangeOptions: { value: TimeRange; label: string }[] = [
  { value: 'all', label: '全部时间' },
  { value: 'today', label: '今天' },
  { value: '3d', label: '近 3 天' },
  { value: '7d', label: '近 7 天' },
  { value: 'custom', label: '自定义' },
]

const sourceOptions: { value: 'all' | SourceType; label: string }[] = [
  { value: 'all', label: '全部来源' },
  { value: 'group', label: '群消息' },
  { value: 'private', label: '私聊消息' },
]

const trustOptions = (Object.keys(trustLevelConfig) as TrustLevel[]).map((key) => ({
  value: key,
  label: trustLevelConfig[key].label,
}))

function ChipButton({
  active,
  onClick,
  children,
  className,
}: {
  active: boolean
  onClick: () => void
  children: React.ReactNode
  className?: string
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-pressed={active}
      className={cn(
        'rounded-full border px-2.5 py-1 text-xs transition-colors',
        active
          ? 'border-primary bg-primary text-primary-foreground'
          : 'bg-secondary text-secondary-foreground hover:border-primary/40 hover:bg-accent hover:text-accent-foreground',
        className,
      )}
    >
      {children}
    </button>
  )
}

export function FilterPanel({
  filters,
  onChange,
  keywordOptions,
}: {
  filters: DashboardFilters
  onChange: (next: DashboardFilters) => void
  keywordOptions: string[]
}) {
  const activeCount = countActiveAdvancedFilters(filters)

  const toggleInArray = <T,>(arr: T[], value: T): T[] =>
    arr.includes(value) ? arr.filter((v) => v !== value) : [...arr, value]

  return (
    <Popover>
      <PopoverTrigger
        render={
          <Button variant="outline" size="sm" className="gap-1.5 bg-transparent">
            <SlidersHorizontal className="size-3.5" />
            筛选
            {activeCount > 0 && (
              <Badge className="size-4.5 justify-center rounded-full p-0 text-[10px]">{activeCount}</Badge>
            )}
          </Button>
        }
      />
      <PopoverContent align="start" className="max-h-[70svh] w-80 overflow-y-auto p-4 md:w-96">
        <div className="flex flex-col gap-4">
          <div className="flex items-center justify-between">
            <p className="text-sm font-semibold">高级筛选</p>
            <Button
              variant="ghost"
              size="sm"
              className="h-7 gap-1 text-xs text-muted-foreground"
              onClick={() =>
                onChange({ ...defaultFilters, status: filters.status, platform: filters.platform, sort: filters.sort })
              }
            >
              <RotateCcw className="size-3" />
              重置
            </Button>
          </div>

          <FilterSection title="时间范围">
            <div className="flex flex-wrap gap-1.5">
              {timeRangeOptions.map((opt) => (
                <ChipButton
                  key={opt.value}
                  active={filters.timeRange === opt.value}
                  onClick={() => onChange({ ...filters, timeRange: opt.value })}
                >
                  {opt.label}
                </ChipButton>
              ))}
            </div>
            {filters.timeRange === 'custom' && (
              <div className="mt-2 flex items-center gap-2">
                <Input
                  type="date"
                  value={filters.customFrom}
                  onChange={(e) => onChange({ ...filters, customFrom: e.target.value })}
                  className="h-8 text-xs"
                  aria-label="开始日期"
                />
                <span className="text-xs text-muted-foreground">至</span>
                <Input
                  type="date"
                  value={filters.customTo}
                  onChange={(e) => onChange({ ...filters, customTo: e.target.value })}
                  className="h-8 text-xs"
                  aria-label="结束日期"
                />
              </div>
            )}
          </FilterSection>

          <FilterSection title="消息来源">
            <div className="flex flex-wrap gap-1.5">
              {sourceOptions.map((opt) => (
                <ChipButton
                  key={opt.value}
                  active={filters.source === opt.value}
                  onClick={() => onChange({ ...filters, source: opt.value })}
                >
                  {opt.label}
                </ChipButton>
              ))}
            </div>
          </FilterSection>

          <FilterSection title="可信度（多选）">
            <div className="flex flex-wrap gap-1.5">
              {trustOptions.map((opt) => (
                <ChipButton
                  key={opt.value}
                  active={filters.trustLevels.includes(opt.value)}
                  onClick={() => onChange({ ...filters, trustLevels: toggleInArray(filters.trustLevels, opt.value) })}
                  className={filters.trustLevels.includes(opt.value) ? undefined : trustLevelConfig[opt.value].className}
                >
                  {opt.label}
                </ChipButton>
              ))}
            </div>
          </FilterSection>

          <FilterSection title="流程阶段（多选）">
            <div className="grid grid-cols-2 gap-x-3 gap-y-2">
              {sopStageOrder.map((stage: SopStage) => (
                <Label key={stage} className="flex cursor-pointer items-center gap-2 text-xs font-normal">
                  <Checkbox
                    checked={filters.stages.includes(stage)}
                    onCheckedChange={() => onChange({ ...filters, stages: toggleInArray(filters.stages, stage) })}
                  />
                  <span className="flex items-center gap-1.5">
                    <span className={cn('size-1.5 rounded-full', sopStageConfig[stage].dotClass)} aria-hidden="true" />
                    {sopStageConfig[stage].label}
                  </span>
                </Label>
              ))}
            </div>
          </FilterSection>

          <FilterSection title="关键词标签（多选）">
            <div className="flex max-h-32 flex-wrap gap-1.5 overflow-y-auto">
              {keywordOptions.map((keyword) => (
                <ChipButton
                  key={keyword}
                  active={filters.keywords.includes(keyword)}
                  onClick={() => onChange({ ...filters, keywords: toggleInArray(filters.keywords, keyword) })}
                >
                  {keyword}
                </ChipButton>
              ))}
            </div>
          </FilterSection>
        </div>
      </PopoverContent>
    </Popover>
  )
}

function FilterSection({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div>
      <p className="mb-2 text-xs font-medium text-muted-foreground">{title}</p>
      {children}
    </div>
  )
}
