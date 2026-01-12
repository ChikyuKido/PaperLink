<script setup lang="ts">
import { computed, onMounted, ref } from "vue"
import { Plug, RefreshCcw, UserPlus } from "lucide-vue-next"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Separator } from "@/components/ui/separator"

import { createD4SAccount, deleteD4SAccount, syncD4SAccounts } from "@/lib/d4s_api"

type Notice = { type: "success" | "error"; message: string } | null

const notice = ref<Notice>(null)
let noticeTimer: number | undefined
function showNotice(n: Notice) {
  notice.value = n
  if (noticeTimer) window.clearTimeout(noticeTimer)
  if (n) noticeTimer = window.setTimeout(() => (notice.value = null), 3500)
}

// Digi4School integration state
const createOpen = ref(false)
const username = ref("")
const password = ref("")

// NOTE: backend doesn't expose an account list endpoint yet, so we keep a local list.
const accounts = ref<{ id: number }[]>([])
const selectedAccountIds = ref(new Set<number>())

const syncing = ref(false)
const selectedCount = computed(() => selectedAccountIds.value.size)

function addAccount(id: number) {
  accounts.value = [{ id }, ...accounts.value]
}

async function submitCreateAccount() {
  if (!username.value.trim() || !password.value) return
  try {
    const id = await createD4SAccount(username.value.trim(), password.value)
    addAccount(id)
    createOpen.value = false
    username.value = ""
    password.value = ""
    showNotice({ type: "success", message: `Account created (#${id}).` })
  } catch (e: any) {
    showNotice({ type: "error", message: e?.message ?? "Failed to create account" })
  }
}

async function onDeleteAccount(id: number) {
  try {
    await deleteD4SAccount(id)
    accounts.value = accounts.value.filter((a) => a.id !== id)
    selectedAccountIds.value.delete(id)
    showNotice({ type: "success", message: `Account deleted (#${id}).` })
  } catch (e: any) {
    showNotice({ type: "error", message: e?.message ?? "Failed to delete account" })
  }
}

async function onSync(ids: "all" | number[]) {
  syncing.value = true
  try {
    const taskId = await syncD4SAccounts(ids)
    showNotice({ type: "success", message: `Sync started (task ${taskId}).` })
  } catch (e: any) {
    showNotice({ type: "error", message: e?.message ?? "Failed to start sync" })
  } finally {
    syncing.value = false
  }
}

onMounted(() => {
  // future: load accounts from backend list endpoint
})
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
            <Plug class="h-5 w-5" />
          </div>
          <div>
            <h1 class="text-lg font-semibold tracking-tight">Integrations</h1>
            <p class="text-xs text-neutral-500 dark:text-neutral-400">Manage external services connected to your instance.</p>
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
        <div class="flex items-center justify-between gap-2">
          <div>
            <CardTitle class="text-sm">Digi4School</CardTitle>
            <CardDescription class="text-[11px]">
              Link accounts used for syncing the Digi4School book library.
            </CardDescription>
          </div>
          <Button
            class="rounded-full bg-emerald-700 text-white hover:bg-emerald-700/90"
            @click="createOpen = true"
          >
            <UserPlus class="h-4 w-4" />
            Add account
          </Button>
        </div>
      </CardHeader>

      <CardContent class="space-y-3">
        <div
          v-if="!accounts.length"
          class="rounded-xl border border-dashed border-neutral-300 bg-neutral-50 p-4 text-sm text-neutral-600 dark:border-neutral-700 dark:bg-neutral-900/40 dark:text-neutral-300"
        >
          No accounts yet. Add one, then run sync.
        </div>

        <div v-else class="space-y-2">
          <div
            v-for="acc in accounts"
            :key="acc.id"
            class="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-neutral-200 bg-white px-3 py-2 dark:border-neutral-800 dark:bg-neutral-900"
          >
            <div class="flex items-center gap-2">
              <button
                type="button"
                class="h-5 w-5 rounded border border-neutral-300 bg-white dark:border-neutral-700 dark:bg-neutral-950"
                :class="selectedAccountIds.has(acc.id) ? 'ring-2 ring-emerald-500/50 border-emerald-600' : ''"
                @click="
                  selectedAccountIds.has(acc.id)
                    ? selectedAccountIds.delete(acc.id)
                    : selectedAccountIds.add(acc.id)
                "
                :aria-label="`Select account ${acc.id}`"
              />

              <div>
                <p class="text-sm font-medium">Account #{{ acc.id }}</p>
                <p class="text-[11px] text-neutral-500 dark:text-neutral-400">Used for syncing</p>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <Button variant="outline" class="rounded-full" @click="onDeleteAccount(acc.id)">
                Delete
              </Button>
            </div>
          </div>
        </div>

        <Separator />

        <div class="flex flex-wrap items-center gap-2">
          <Button
            class="rounded-full bg-emerald-700 text-white hover:bg-emerald-700/90"
            :disabled="syncing"
            @click="onSync('all')"
          >
            <RefreshCcw class="h-4 w-4" />
            Sync all
          </Button>

          <Button
            variant="outline"
            class="rounded-full"
            :disabled="syncing || !selectedCount"
            @click="onSync(Array.from(selectedAccountIds))"
          >
            Sync selected ({{ selectedCount }})
          </Button>
        </div>
      </CardContent>
    </Card>

    <Dialog :open="createOpen" @update:open="(v) => (createOpen = v)">
      <DialogContent class="sm:max-w-[520px]">
        <DialogHeader>
          <DialogTitle>Add Digi4School account</DialogTitle>
          <DialogDescription>
            Credentials are used only for syncing.
          </DialogDescription>
        </DialogHeader>

        <div class="space-y-3">
          <div class="space-y-1.5">
            <p class="text-xs font-medium text-neutral-600 dark:text-neutral-300">Username</p>
            <Input v-model="username" placeholder="your.username" />
          </div>
          <div class="space-y-1.5">
            <p class="text-xs font-medium text-neutral-600 dark:text-neutral-300">Password</p>
            <Input v-model="password" type="password" placeholder="••••••••" />
          </div>
        </div>

        <DialogFooter class="mt-4">
          <Button variant="outline" class="rounded-full" @click="createOpen = false">Cancel</Button>
          <Button
            class="rounded-full bg-emerald-700 text-white hover:bg-emerald-700/90"
            @click="submitCreateAccount"
            :disabled="!username.trim() || !password"
          >
            Add
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
