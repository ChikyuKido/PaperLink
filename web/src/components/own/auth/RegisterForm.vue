<script setup lang="ts">
import { ref } from "vue"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { User, Lock, KeyRound, Eye, EyeOff } from "lucide-vue-next"

const emit = defineEmits<{
  (e: "success"): void
  (e: "status"): void
}>()

const REGISTER_ENDPOINT = "/api/v1/auth/register"

const username = ref("")
const password = ref("")
const confirm = ref("")
const inviteCode = ref("")

const showPassword = ref(false)
const isSubmitting = ref(false)
const submitted = ref(false)

const serverError = ref<string | null>(null)
const serverField = ref<"username" | null>(null)

function errorFor(field: "username" | "password" | "confirm" | "invite") {
  if (!submitted.value) return ""

  if (field === "username") {
    if (!username.value) return "Username is required"
    if (serverField.value === "username") return serverError.value
  }

  if (field === "password") {
    if (!password.value) return "Password is required"
    if (password.value.length < 8) return "Password must be at least 8 characters"
  }

  if (field === "confirm") {
    if (!confirm.value) return "Please confirm your password"
    if (confirm.value !== password.value) return "Passwords do not match"
  }

  if (field === "invite" && !inviteCode.value) {
    return "Invite code is required"
  }

  return ""
}

async function onSubmit() {
  submitted.value = true
  serverError.value = null
  serverField.value = null

  if (
      errorFor("username") ||
      errorFor("password") ||
      errorFor("confirm") ||
      errorFor("invite")
  ) return

  emit("status")
  isSubmitting.value = true

  try {
    const res = await fetch(REGISTER_ENDPOINT, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
        inviteCode: inviteCode.value,
      }),
    })

    const body = await res.json()

    if (!res.ok) throw body
    emit("success")
  } catch (e: any) {
    serverError.value = e?.error || "Registration failed"
    if (serverError.value.toLowerCase().includes("username")) {
      serverField.value = "username"
    }
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <form class="space-y-4" @submit.prevent="onSubmit">
    <div class="relative">
      <User class="absolute left-3 top-3 size-4 text-neutral-400" />
      <Input v-model="username" placeholder="Username" class="pl-9" />
      <p v-if="errorFor('username')" class="text-sm text-red-500">
        {{ errorFor("username") }}
      </p>
    </div>

    <div class="relative">
      <Lock class="absolute left-3 top-3 size-4 text-neutral-400" />
      <Input
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          placeholder="Password"
          class="pl-9 pr-10"
      />
      <button
          type="button"
          class="absolute right-2 top-2 p-2"
          @click="showPassword = !showPassword"
      >
        <Eye v-if="!showPassword" class="size-4" />
        <EyeOff v-else class="size-4" />
      </button>
      <p v-if="errorFor('password')" class="text-sm text-red-500">
        {{ errorFor("password") }}
      </p>
    </div>

    <div class="relative">
      <Lock class="absolute left-3 top-3 size-4 text-neutral-400" />
      <Input
          v-model="confirm"
          :type="showPassword ? 'text' : 'password'"
          placeholder="Confirm password"
          class="pl-9"
      />
      <p v-if="errorFor('confirm')" class="text-sm text-red-500">
        {{ errorFor("confirm") }}
      </p>
    </div>

    <div class="relative">
      <KeyRound class="absolute left-3 top-3 size-4 text-neutral-400" />
      <Input v-model="inviteCode" placeholder="Invite code" class="pl-9" />
      <p v-if="errorFor('invite')" class="text-sm text-red-500">
        {{ errorFor("invite") }}
      </p>
    </div>

    <p v-if="submitted && serverError && !serverField" class="text-sm text-red-600">
      {{ serverError }}
    </p>

    <Button type="submit" class="w-full" :disabled="isSubmitting">
      {{ isSubmitting ? "Creating account…" : "Create account" }}
    </Button>
  </form>
</template>
