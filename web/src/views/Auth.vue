<script setup lang="ts">
import {ref, computed, watch} from "vue"
import { useRouter } from "vue-router"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { Alert, AlertDescription } from "@/components/ui/alert"

import LoginForm from "@/components/own/auth/LoginForm.vue"
import RegisterForm from "@/components/own/auth/RegisterForm.vue"
import ExternalAuthForm from "@/components/own/auth/ExternalAuthForm.vue"

type Tab = "login" | "register" | "external"

const router = useRouter()
const activeTab = ref<Tab>("login")


const errorMessage = ref("")
const successMessage = ref("")

const title = computed(() => {
  if (activeTab.value === "login") return "Sign in to Paperlink"
  if (activeTab.value === "register") return "Create a Paperlink account"
  return "Continue with an external provider"
})

const description = computed(() => {
  if (activeTab.value === "login") return "Access and collaborate on your PDFs."
  if (activeTab.value === "register") return "Create an account using your invite code."
  return "Use your organization identity provider to sign in."
})

function clear() {
  errorMessage.value = ""
  successMessage.value = ""
}

function onLoginSuccess() {
  router.push("/")
}

function onRegisterSuccess() {
  successMessage.value = "Account created. You can now sign in."
  activeTab.value = "login"
}
watch(activeTab, () => {
  clear()
})
</script>

<template>
  <div class="min-h-screen grid grid-cols-1 lg:grid-cols-2 bg-neutral-950 text-neutral-50">
    <div
        class="hidden lg:flex flex-col justify-center px-20
             bg-neutral-950 border-r border-neutral-800"
    >
      <h1 class="text-5xl font-semibold tracking-wide">PAPERLINK</h1>
      <p class="mt-3 text-sm text-neutral-400 tracking-wide">
        SELF-HOSTED PDF MANAGEMENT
      </p>
    </div>

    <div class="flex items-center justify-center px-6">
      <div
          class="w-full max-w-2xl rounded-2xl border border-neutral-800
               bg-neutral-900/70 shadow-xl
               ring-1 ring-emerald-500/10"
      >
        <Tabs v-model="activeTab" class="w-full">
          <div class="p-5 pb-4">
            <TabsList
                class="grid w-full grid-cols-3 rounded-xl bg-neutral-800/60 p-1"
            >
              <TabsTrigger
                  value="login"
                  class="
                  rounded-lg text-neutral-300
                  data-[state=active]:bg-neutral-700
                  data-[state=active]:text-white
                  data-[state=active]:shadow
                  data-[state=active]:ring-1
                  data-[state=active]:ring-emerald-500/30
                "
              >
                Login
              </TabsTrigger>

              <TabsTrigger
                  value="register"
                  class="
                  rounded-lg text-neutral-300
                  data-[state=active]:bg-neutral-700
                  data-[state=active]:text-white
                  data-[state=active]:shadow
                  data-[state=active]:ring-1
                  data-[state=active]:ring-emerald-500/30
                "
              >
                Register
              </TabsTrigger>

              <TabsTrigger
                  value="external"
                  class="
                  rounded-lg text-neutral-300
                  data-[state=active]:bg-neutral-700
                  data-[state=active]:text-white
                  data-[state=active]:shadow
                  data-[state=active]:ring-1
                  data-[state=active]:ring-emerald-500/30
                "
              >
                External
              </TabsTrigger>
            </TabsList>
          </div>

          <div class="px-10 pb-10 pt-6">
            <h2 class="text-2xl font-semibold tracking-tight">
              {{ title }}
            </h2>

            <div class="mt-3 mb-4 h-px w-12 bg-emerald-500/60 rounded-full" />

            <p class="text-sm text-neutral-400">
              {{ description }}
            </p>

            <div class="mt-6 space-y-3">
              <Alert v-if="errorMessage" variant="destructive">
                <AlertDescription>{{ errorMessage }}</AlertDescription>
              </Alert>

              <Alert
                  v-if="successMessage"
                  class="border-emerald-500/40 bg-emerald-900/20 text-emerald-200"
              >
                <AlertDescription>{{ successMessage }}</AlertDescription>
              </Alert>
            </div>

            <div class="mt-8">
              <TabsContent value="login">
                <LoginForm
                    @success="onLoginSuccess"
                    @error="(e) => (errorMessage = e)"
                    @status="clear"
                />
              </TabsContent>

              <TabsContent value="register">
                <RegisterForm
                    @success="onRegisterSuccess"
                    @error="(e) => (errorMessage = e)"
                    @status="clear"
                />
              </TabsContent>

              <TabsContent value="external">
                <ExternalAuthForm
                    @started="clear"
                    @error="(e) => (errorMessage = e)"
                    @status="clear"
                />
              </TabsContent>
            </div>
          </div>
        </Tabs>
      </div>
    </div>
  </div>
</template>
