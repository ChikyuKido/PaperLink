<script setup lang="ts">
import { computed, ref } from "vue"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { User, Lock, Eye, EyeOff } from "lucide-vue-next"

const emit = defineEmits<{
  (e: "success"): void
  (e: "error", message: string): void
  (e: "status"): void
}>()

const LOGIN_ENDPOINT = "/api/auth/login"

const username = ref("")
const password = ref("")
const isSubmitting = ref(false)
const showPassword = ref(false)

const canSubmit = computed(() => !!username.value && !!password.value && !isSubmitting.value)

async function onSubmit() {
  if (!canSubmit.value) return
  emit("status")
  isSubmitting.value = true

  try {
    const res = await fetch(LOGIN_ENDPOINT, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ username: username.value, password: password.value }),
    })

    if (!res.ok) {
      const body = await res.json().catch(() => null)
      throw new Error(body?.message || "Login failed")
    }

    emit("success")
  } catch (e) {
    emit("error", e instanceof Error ? e.message : "Unexpected login error")
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="onSubmit">
    <div class="relative">
      <User class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-neutral-500" />
      <Input
          v-model="username"
          autocomplete="username"
          placeholder="Username"
          class="pl-9 bg-neutral-800 border-neutral-700 text-neutral-50 placeholder:text-neutral-500"
      />
    </div>

    <div class="relative">
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
          class="absolute right-2 top-1/2 -translate-y-1/2 rounded-md p-2 text-neutral-500 hover:text-neutral-200"
          @click="showPassword = !showPassword"
      >
        <Eye v-if="!showPassword" class="size-4" />
        <EyeOff v-else class="size-4" />
      </button>
    </div>

    <Button
        type="submit"
        class="w-full bg-neutral-100 text-neutral-900 hover:bg-white"
        :disabled="!canSubmit"
    >
      <span v-if="isSubmitting">Logging in…</span>
      <span v-else>Login</span>
    </Button>
  </form>
</template>
