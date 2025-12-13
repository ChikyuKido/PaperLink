<script setup lang="ts">
import { Button } from "@/components/ui/button"
import { Shield } from "lucide-vue-next"

const emit = defineEmits<{
  (e: "started"): void
  (e: "error", message: string): void
  (e: "status"): void
}>()

const OIDC_START_URL = "/api/auth/oidc/start"

function continueWithOidc() {
  emit("status")
  emit("started")
  try {
    window.location.assign(OIDC_START_URL)
  } catch {
    emit("error", "Could not start external login.")
  }
}
</script>

<template>
  <div class="space-y-4">
    <Button
        type="button"
        class="w-full bg-neutral-100 text-neutral-900 hover:bg-white"
        @click="continueWithOidc"
    >
      <Shield class="mr-2 size-4" />
      Continue with OIDC
    </Button>

    <p class="text-xs text-neutral-500">
      You will be redirected to your identity provider and returned to Paperlink after signing in.
    </p>
  </div>
</template>
