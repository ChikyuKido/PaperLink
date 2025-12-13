<script setup lang="ts">
import { computed, ref } from "vue"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { User, Lock, KeyRound, Eye, EyeOff } from "lucide-vue-next"

const emit = defineEmits<{
  (e: "success"): void
  (e: "error", message: string): void
  (e: "status"): void
}>()

const REGISTER_ENDPOINT = "/api/auth/register"

const username = ref("")
const password = ref("")
const passwordConfirm = ref("")
const inviteCode = ref("")

const isSubmitting = ref(false)
const showPassword = ref(false)
const showConfirm = ref(false)

const passwordsMatch = computed(() => password.value === passwordConfirm.value)

const canSubmit = computed(() => {
  return (
      !!username.value &&
      password.value.length >= 8 &&
      passwordsMatch.value &&
      !!inviteCode.value &&
      !isSubmitting.value
  )
})

async function postJson<T>(url: string, payload: unknown): Promise<T> {
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify(payload),
  })

  if (!res.ok) {
    const body = await res.json().catch(() => null)
    throw new Error(body?.message || "Registration failed")
  }

  return (await res.json().catch(() => ({}))) as T
}

async function onSubmit() {
  if (!canSubmit.value) return
  emit("status")
  isSubmitting.value = true

  try {
    await postJson(REGISTER_ENDPOINT, {
      username: username.value,
      password: password.value,
      inviteCode: inviteCode.value,
    })
    emit("success")
  } catch (e) {
    emit("error", e instanceof Error ? e.message : "Unexpected registration error")
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="onSubmit">
    <div class="relative">
      <User class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-400" />
      <Input
          v-model="username"
          autocomplete="username"
          placeholder="Username"
          class="pl-9"
      />
    </div>

    <div class="relative">
      <Lock class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-400" />
      <Input
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          autocomplete="new-password"
          placeholder="Password (min 8 chars)"
          class="pl-9 pr-10"
      />
      <button
          type="button"
          class="absolute right-2 top-1/2 -translate-y-1/2 rounded-md p-2 text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-50"
          @click="showPassword = !showPassword"
      >
        <Eye v-if="!showPassword" class="size-4" />
        <EyeOff v-else class="size-4" />
      </button>
    </div>

    <div class="relative">
      <Lock class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-400" />
      <Input
          v-model="passwordConfirm"
          :type="showConfirm ? 'text' : 'password'"
          autocomplete="new-password"
          placeholder="Confirm password"
          class="pl-9 pr-10"
      />
      <button
          type="button"
          class="absolute right-2 top-1/2 -translate-y-1/2 rounded-md p-2 text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-50"
          @click="showConfirm = !showConfirm"
      >
        <Eye v-if="!showConfirm" class="size-4" />
        <EyeOff v-else class="size-4" />
      </button>
    </div>

    <div class="relative">
      <KeyRound class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-400" />
      <Input
          v-model="inviteCode"
          autocomplete="off"
          placeholder="Invite code"
          class="pl-9"
      />
    </div>

    <Button type="submit" class="w-full" :disabled="!canSubmit">
      <span v-if="isSubmitting">Creating account…</span>
      <span v-else>Create account</span>
    </Button>
  </form>
</template>
