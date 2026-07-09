'use client'

import { Pencil, Plus } from 'lucide-react'
import { useMemo, useState } from 'react'
import { TemplateDialog } from '@/components/template-dialog'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { useAppStore } from '@/lib/app-store'
import { templateCategories } from '@/lib/mock-data'
import type { ReplyTemplate } from '@/lib/types'
import { cn } from '@/lib/utils'

export default function TemplatesPage() {
  const { templates } = useAppStore()
  const [activeCategory, setActiveCategory] = useState('全部')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingTemplate, setEditingTemplate] = useState<ReplyTemplate | null>(null)

  const filtered = useMemo(
    () => (activeCategory === '全部' ? templates : templates.filter((t) => t.category === activeCategory)),
    [templates, activeCategory],
  )

  const openCreate = () => {
    setEditingTemplate(null)
    setDialogOpen(true)
  }

  const openEdit = (template: ReplyTemplate) => {
    setEditingTemplate(template)
    setDialogOpen(true)
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-6 md:px-8">
      <header className="mb-6 flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-xl font-semibold tracking-tight md:text-2xl">回复模板库</h1>
          <p className="mt-1 text-sm text-muted-foreground">维护常用回复话术，支持变量占位符</p>
        </div>
        <Button onClick={openCreate} className="gap-1.5">
          <Plus className="size-4" />
          新建模板
        </Button>
      </header>

      <div className="mb-5 flex gap-2 overflow-x-auto pb-1" role="tablist" aria-label="模板分类">
        {templateCategories.map((category) => (
          <button
            key={category}
            type="button"
            role="tab"
            aria-selected={activeCategory === category}
            onClick={() => setActiveCategory(category)}
            className={cn(
              'shrink-0 rounded-full border px-3.5 py-1.5 text-xs font-medium transition-colors',
              activeCategory === category
                ? 'border-primary bg-primary text-primary-foreground'
                : 'bg-card text-muted-foreground hover:border-primary/40 hover:text-foreground',
            )}
          >
            {category}
          </button>
        ))}
      </div>

      <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
        {filtered.map((template) => (
          <Card key={template.id} className="group gap-2.5 rounded-xl p-4 shadow-sm transition-shadow hover:shadow-md">
            <div className="flex items-start justify-between gap-2">
              <div className="min-w-0">
                <p className="truncate text-sm font-semibold">{template.title}</p>
                <Badge variant="secondary" className="mt-1.5 h-5 rounded-md px-1.5 text-[10px] font-normal">
                  {template.category}
                </Badge>
              </div>
              <Button
                variant="ghost"
                size="icon"
                className="size-7 opacity-0 transition-opacity group-hover:opacity-100 focus-visible:opacity-100"
                onClick={() => openEdit(template)}
                aria-label={`编辑模板 ${template.title}`}
              >
                <Pencil className="size-3.5" />
              </Button>
            </div>
            <p className="line-clamp-3 text-xs leading-relaxed text-muted-foreground">{template.content}</p>
          </Card>
        ))}
      </div>

      <TemplateDialog open={dialogOpen} onOpenChange={setDialogOpen} template={editingTemplate} />
    </div>
  )
}
