'use client'

import { AlertTriangle, Check, Loader2, Lock, Minus } from 'lucide-react'
import { useEffect, useState } from 'react'
import {
  StepChat,
  StepContacts,
  StepDiscovery,
  StepFriendRequest,
  StepLinkVerification,
} from '@/components/sop-step-panels'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Card } from '@/components/ui/card'
import { formatDateTime } from '@/lib/sop'
import { deriveSteps, type SopStep, type StepState } from '@/lib/sop-steps'
import type { Opportunity } from '@/lib/types'
import { cn } from '@/lib/utils'

const stateVisual: Record<
  StepState,
  { labelClass: string; circleClass: string }
> = {
  done: { labelClass: 'text-success', circleClass: 'border-success bg-success text-success-foreground text-white' },
  active: { labelClass: 'text-primary', circleClass: 'border-primary bg-primary/10 text-primary' },
  in_progress: { labelClass: 'text-primary', circleClass: 'border-primary bg-primary/10 text-primary' },
  blocked: { labelClass: 'text-muted-foreground', circleClass: 'border-border bg-muted text-muted-foreground' },
  failed: { labelClass: 'text-destructive', circleClass: 'border-destructive bg-destructive/10 text-destructive' },
  skipped: { labelClass: 'text-muted-foreground', circleClass: 'border-border bg-muted text-muted-foreground' },
  locked: { labelClass: 'text-muted-foreground', circleClass: 'border-border bg-muted text-muted-foreground' },
}

function StepCircle({ step }: { step: SopStep }) {
  const visual = stateVisual[step.state]
  return (
    <span
      className={cn(
        'flex size-7 shrink-0 items-center justify-center rounded-full border-2 text-xs font-semibold',
        visual.circleClass,
      )}
      aria-hidden="true"
    >
      {step.state === 'done' ? (
        <Check className="size-3.5" />
      ) : step.state === 'in_progress' ? (
        <Loader2 className="size-3.5 animate-spin" />
      ) : step.state === 'failed' ? (
        <AlertTriangle className="size-3.5" />
      ) : step.state === 'skipped' ? (
        <Minus className="size-3.5" />
      ) : step.state === 'locked' || step.state === 'blocked' ? (
        <Lock className="size-3" />
      ) : (
        step.index
      )}
    </span>
  )
}

function StepPanel({ opportunity, step }: { opportunity: Opportunity; step: SopStep }) {
  switch (step.key) {
    case 'discovery':
      return <StepDiscovery opportunity={opportunity} />
    case 'link_verification':
      return <StepLinkVerification opportunity={opportunity} step={step} />
    case 'contacts':
      return <StepContacts opportunity={opportunity} step={step} />
    case 'friend_request':
      return <StepFriendRequest opportunity={opportunity} step={step} />
    case 'chat':
      return <StepChat opportunity={opportunity} step={step} />
  }
}

function currentStepKey(steps: SopStep[]) {
  const focus = steps.find((s) => s.state === 'active' || s.state === 'in_progress' || s.state === 'failed')
  return (focus ?? steps[steps.length - 1]).key
}

export function SopStepper({ opportunity }: { opportunity: Opportunity }) {
  const steps = deriveSteps(opportunity)
  const autoKey = currentStepKey(steps)
  const [selectedKey, setSelectedKey] = useState<SopStep['key']>(autoKey)
  const [lastAutoKey, setLastAutoKey] = useState(autoKey)

  // 流程推进时自动跟随到新的当前步骤
  useEffect(() => {
    if (autoKey !== lastAutoKey) {
      setLastAutoKey(autoKey)
      setSelectedKey(autoKey)
    }
  }, [autoKey, lastAutoKey])

  const selectedStep = steps.find((s) => s.key === selectedKey) ?? steps[0]

  return (
    <>
      {/* 桌面端：左侧时间线 + 右侧操作区 */}
      <div className="hidden gap-4 lg:flex">
        <Card className="h-fit w-72 shrink-0 gap-0 rounded-xl p-2 shadow-sm">
          <p className="px-3 pb-1 pt-2 text-xs font-medium text-muted-foreground">处理流程</p>
          <nav aria-label="商机处理流程步骤">
            <ol className="flex flex-col">
              {steps.map((step, i) => {
                const visual = stateVisual[step.state]
                const isSelected = step.key === selectedKey
                return (
                  <li key={step.key} className="relative">
                    {i < steps.length - 1 && (
                      <span
                        className={cn(
                          'absolute left-[25px] top-11 h-[calc(100%-2.25rem)] w-0.5',
                          step.state === 'done' ? 'bg-success/40' : 'bg-border',
                        )}
                        aria-hidden="true"
                      />
                    )}
                    <button
                      type="button"
                      onClick={() => setSelectedKey(step.key)}
                      aria-current={isSelected ? 'step' : undefined}
                      className={cn(
                        'relative flex w-full items-start gap-3 rounded-lg p-3 text-left transition-colors hover:bg-muted/60',
                        isSelected && 'bg-accent',
                      )}
                    >
                      <StepCircle step={step} />
                      <span className="min-w-0 flex-1">
                        <span className={cn('block text-sm font-medium', isSelected ? 'text-accent-foreground' : undefined)}>
                          {step.title}
                        </span>
                        <span className={cn('mt-0.5 block text-[11px]', visual.labelClass)}>{step.stateLabel}</span>
                        {step.completedAt && step.state === 'done' && (
                          <span className="mt-0.5 block text-[10px] text-muted-foreground">
                            完成于 {formatDateTime(step.completedAt)}
                          </span>
                        )}
                      </span>
                    </button>
                  </li>
                )
              })}
            </ol>
          </nav>
        </Card>

        <Card className="min-w-0 flex-1 gap-0 rounded-xl p-5 shadow-sm">
          <div className="mb-4 flex items-center gap-3 border-b pb-3">
            <StepCircle step={selectedStep} />
            <div>
              <h2 className="text-sm font-semibold">
                Step {selectedStep.index} · {selectedStep.title}
              </h2>
              <p className={cn('text-xs', stateVisual[selectedStep.state].labelClass)}>{selectedStep.stateLabel}</p>
            </div>
          </div>
          <StepPanel opportunity={opportunity} step={selectedStep} />
        </Card>
      </div>

      {/* 移动端：手风琴 */}
      <Card className="gap-0 rounded-xl px-4 py-1 shadow-sm lg:hidden">
        <Accordion value={[selectedKey]} onValueChange={(v) => setSelectedKey((v[0] as SopStep['key']) ?? selectedKey)}>
          {steps.map((step) => {
            const visual = stateVisual[step.state]
            return (
              <AccordionItem key={step.key} value={step.key}>
                <AccordionTrigger className="hover:no-underline">
                  <span className="flex items-center gap-3">
                    <StepCircle step={step} />
                    <span>
                      <span className="block text-sm font-medium">
                        Step {step.index} · {step.title}
                      </span>
                      <span className={cn('mt-0.5 block text-[11px] font-normal', visual.labelClass)}>
                        {step.stateLabel}
                        {step.completedAt && step.state === 'done' && ` · ${formatDateTime(step.completedAt)}`}
                      </span>
                    </span>
                  </span>
                </AccordionTrigger>
                <AccordionContent>
                  <StepPanel opportunity={opportunity} step={step} />
                </AccordionContent>
              </AccordionItem>
            )
          })}
        </Accordion>
      </Card>
    </>
  )
}
