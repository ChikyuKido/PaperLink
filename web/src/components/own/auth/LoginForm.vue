<script setup lang="ts">
import { computed, ref } from "vue"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { User, Lock, Eye, EyeOff } from "lucide-vue-next"
import {setAccessToken} from "@/auth/auth.ts";

const emit = defineEmits<{
  (e: "success"): void
  (e: "error", message: string): void
  (e: "status"): void
}>()

const LOGIN_ENDPOINT = "/api/v1/auth/login"

const username = ref("")
const password = ref("")
const showPassword = ref(false)

const isSubmitting = ref(false)
const submitted = ref(false)
const serverError = ref<string | null>(null)

const usernameError = computed(() => {
  if (!submitted.value) return ""
  if (!username.value) return "Username is required"
  return ""
})

const passwordError = computed(() => {
  if (!submitted.value) return ""
  if (!password.value) return "Password is required"
  return ""
})

const canSubmit = computed(
    () => !!username.value && !!password.value && !isSubmitting.value
)

async function onSubmit() {
  submitted.value = true
  serverError.value = null

  if (!canSubmit.value) return

  emit("status")
  isSubmitting.value = true

  try {
    const res = await fetch(LOGIN_ENDPOINT, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        username: username.value,
        password: password.value,
      }),
    })

    const body = await res.json().catch(() => null)

    if (!res.ok) {
      throw new Error(body?.error || body?.message || "Login failed")
    }
    setAccessToken(body.data.access)
    emit("success")

  } catch (e) {
    const message =
        e instanceof Error ? e.message : "Unexpected login error"
    serverError.value = message
    emit("error", message)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="onSubmit">
    <div class="relative space-y-1">
      <User class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-500" />
      <Input
          v-model="username"
          autocomplete="username"
          placeholder="Username"
          class="pl-9 bg-neutral-800 border-neutral-700 text-neutral-50 placeholder:text-neutral-500"
      />
      <p v-if="usernameError" class="text-sm text-red-500">
        {{ usernameError }}
      </p>
    </div>

    <div class="relative space-y-1">
      <Lock class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-500" />
      <Input
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          autocomplete="current-password"
          placeholder="Password"
          class="pl-9 pr-10 bg-neutral-800 border-neutral-700 text-neutral-50 placeholder:text-neutral-500"
      />
      <button
          type="button"
          class="absolute right-2 top-1/2 -translate-y-1/2 p-2 text-neutral-500 hover:text-neutral-200"
          @click="showPassword = !showPassword"
      >
        <Eye v-if="!showPassword" class="size-4" />
        <EyeOff v-else class="size-4" />
      </button>
      <p v-if="passwordError" class="text-sm text-red-500">
        {{ passwordError }}
      </p>
    </div>

    <Button
        type="submit"
        class="w-full bg-neutral-100 text-neutral-900 hover:bg-white"
        :disabled="!canSubmit"
    >
      {{ isSubmitting ? "Logging in…" : "Login" }}
    </Button>
  </form>
</template>
