<script setup lang="ts">
import { ref, computed } from "vue"
import { Ticket, Shield, Copy, Loader2 } from "lucide-vue-next"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { apiFetch } from "@/auth/api"
import { Input } from "@/components/ui/input"

type Notice = { type: "success" | "error"; message: string } | null

const notice = ref<Notice>(null)
const isCreating = ref(false)

const inviteCode = ref<string | null>(null)
const inviteExpiresAt = ref<number | null>(null)
const inviteUses = ref<number | null>(null)

const validDays = ref<number>(3)
const uses = ref<number>(1)

const validDaysText = computed(() => String(validDays.value ?? ""))
const usesText = computed(() => String(uses.value ?? ""))

function formatExpiresAt(ts: number | null) {
  if (!ts) return "—"
  const d = new Date(ts * 1000)
  if (Number.isNaN(d.getTime())) return "—"
  return d.toLocaleString()
}

function showNotice(n: Notice) {
  notice.value = n
  if (n) window.setTimeout(() => (notice.value = null), 3500)
}

async function createInvite() {
  isCreating.value = true
  inviteCode.value = null
  inviteExpiresAt.value = null
  inviteUses.value = null
  try {
    const res = await apiFetch("/api/v1/invite/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        validDays: Number(validDays.value) || 0,
        uses: Number(uses.value) || 0,
      }),
    })
    const body = await res.json().catch(() => null)

    if (!res.ok) {
      showNotice({ type: "error", message: body?.error ?? `Failed to create invite (${res.status})` })
      return
    }

    const code = body?.data?.code as string | undefined
    const expiresAt = body?.data?.expiresAt as number | undefined
    const createdUses = body?.data?.uses as number | undefined

    if (!code) {
      showNotice({ type: "error", message: "Invite created, but no code returned" })
      return
    }

    inviteCode.value = code
    inviteExpiresAt.value = typeof expiresAt === "number" ? expiresAt : null
    inviteUses.value = typeof createdUses === "number" ? createdUses : null
    showNotice({ type: "success", message: "Invite code created." })
  } catch (e: any) {
    showNotice({ type: "error", message: e?.message ?? "Failed to create invite" })
  } finally {
    isCreating.value = false
  }
}

async function copyCode() {
  if (!inviteCode.value) return
  try {
    await navigator.clipboard.writeText(inviteCode.value)
    showNotice({ type: "success", message: "Copied to clipboard." })
  } catch {
    showNotice({ type: "error", message: "Could not copy to clipboard." })
  }
}
</script>

<template>
  <div class="mx-auto max-w-6xl px-4 lg:px-6 py-5 lg:py-7 space-y-4">
    <section
      class="rounded-2xl border border-neutral-200 bg-white shadow-sm shadow-neutral-200/70 overflow-hidden dark:border-neutral-800 dark:bg-neutral-900 dark:shadow-none"
    >
      <div
        class="px-4 sm:px-6 py-4 bg-gradient-to-r from-neutral-50 via-white to-emerald-50/70 dark:from-neutral-900 dark:via-neutral-900 dark:to-emerald-900/30"
      >
        <div class="flex items-center gap-3">
          <div
            class="inline-flex h-10 w-10 items-center justify-center rounded-2xl bg-emerald-600/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-200"
          >
            <Shield class="h-5 w-5" />
          </div>
          <div>
            <h1 class="text-lg font-semibold tracking-tight">Admin · Invites</h1>
            <p class="text-xs text-neutral-500 dark:text-neutral-400">Create one-time invite codes for registration.</p>
          </div>
        </div>

        <div
          v-if="notice"
          class="mt-4 rounded-xl border px-4 py-3 text-sm"
          :class="
            notice.type === 'success'
              ? 'border-emerald-600/30 bg-emerald-600/10 text-emerald-900 dark:text-emerald-200 dark:bg-emerald-500/10'
              : 'border-red-600/30 bg-red-600/10 text-red-900 dark:text-red-200 dark:bg-red-500/10'
          "
        >
          {{ notice.message }}
        </div>
      </div>
    </section>

    <Card class="border border-neutral-200 dark:border-neutral-800">
      <CardHeader>
        <div class="flex items-center gap-2">
          <span class="inline-flex h-8 w-8 items-center justify-center rounded-full bg-emerald-700/10 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-200">
            <Ticket class="h-4 w-4" />
          </span>
          <div>
            <CardTitle class="text-sm">Create invite code</CardTitle>
            <CardDescription class="text-[11px]">Admins only · Uses <span class="font-mono">POST /api/v1/invite/create</span></CardDescription>
          </div>
        </div>
      </CardHeader>

      <CardContent class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div class="space-y-1.5">
            <p class="text-[11px] font-medium text-neutral-600 dark:text-neutral-300">Expires in (days)</p>
            <Input
              :model-value="validDaysText"
              inputmode="numeric"
              class="h-9"
              placeholder="3"
              @update:model-value="(v: string) => (validDays.value = Number(v))"
            />
            <p class="text-[11px] text-neutral-500 dark:text-neutral-400">Leave at 3 for default.</p>
          </div>

          <div class="space-y-1.5">
            <p class="text-[11px] font-medium text-neutral-600 dark:text-neutral-300">Max uses</p>
            <Input
              :model-value="usesText"
              inputmode="numeric"
              class="h-9"
              placeholder="1"
              @update:model-value="(v: string) => (uses.value = Number(v))"
            />
            <p class="text-[11px] text-neutral-500 dark:text-neutral-400">Leave at 1 for single-use.</p>
          </div>
        </div>

        <div class="flex flex-col sm:flex-row gap-2 sm:items-center sm:justify-between">
          <Button
            class="rounded-xl bg-emerald-700 text-white hover:bg-emerald-700/90"
            :disabled="isCreating"
            @click="createInvite"
          >
            <Loader2 v-if="isCreating" class="h-4 w-4 animate-spin" />
            <span v-if="!isCreating">Create invite</span>
            <span v-else>Creating…</span>
          </Button>

          <p class="text-[11px] text-neutral-500 dark:text-neutral-400">
            Share the code with the person you want to invite.
          </p>
        </div>

        <Separator />

        <div class="space-y-2">
          <p class="text-xs font-medium text-neutral-700 dark:text-neutral-200">Latest code</p>

          <div
            class="flex items-center justify-between gap-2 rounded-xl border border-neutral-200 bg-neutral-50 px-3 py-2 dark:border-neutral-800 dark:bg-neutral-900/40"
          >
            <p class="min-w-0 truncate font-mono text-sm text-neutral-900 dark:text-neutral-50">
              {{ inviteCode ?? '—' }}
            </p>

            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8 rounded-lg text-neutral-600 hover:bg-neutral-100 dark:text-neutral-300 dark:hover:bg-neutral-800"
              :disabled="!inviteCode"
              @click="copyCode"
            >
              <Copy class="h-4 w-4" />
            </Button>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-[11px] text-neutral-500 dark:text-neutral-400">
            <p>
              Expires at: <span class="text-neutral-900 dark:text-neutral-50">{{ formatExpiresAt(inviteExpiresAt) }}</span>
            </p>
            <p>
              Uses: <span class="text-neutral-900 dark:text-neutral-50">{{ inviteUses ?? '—' }}</span>
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
