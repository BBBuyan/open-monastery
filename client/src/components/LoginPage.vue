<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import type { LoginMethod } from '../demo-auth'

const props = defineProps<{ signup: boolean }>()
const emit = defineEmits<{ enterDemo: [method: LoginMethod] }>()
const selected = ref<LoginMethod | null>(null)
const preview = ref<HTMLElement>()
const choices = ref<HTMLElement>()
const labels = {
  email: 'Email',
  phone: 'Phone number',
  google: 'Google',
  facebook: 'Facebook',
}
watch(
  () => props.signup,
  () => {
    selected.value = null
  },
)
async function choose(method: LoginMethod) {
  selected.value = method
  await nextTick()
  preview.value?.focus()
}
async function back() {
  const method = selected.value
  selected.value = null
  await nextTick()
  choices.value
    ?.querySelector<HTMLButtonElement>(`[data-method="${method}"]`)
    ?.focus()
}
</script>

<template>
  <section class="login-page" aria-labelledby="login-title">
    <a class="return-link" href="#readings">← Back to books</a>
    <h1 id="login-title">
      {{ signup ? 'Create an account' : 'Welcome back' }}
    </h1>
    <p class="login-intro">
      {{
        signup
          ? 'Choose how you’d like to sign up.'
          : 'Choose how you’d like to log in.'
      }}
    </p>
    <p class="login-demo" id="login-demo-notice">
      <strong>Demo only.</strong> No account will be created, no verification
      message will be sent, and no provider connection will be made.
    </p>
    <div
      v-if="!selected"
      ref="choices"
      class="login-options"
      aria-describedby="login-demo-notice"
    >
      <button
        class="button secondary"
        type="button"
        data-method="email"
        @click="choose('email')"
      >
        <span aria-hidden="true">✉</span> Continue with email
      </button>
      <button
        class="button secondary"
        type="button"
        data-method="phone"
        @click="choose('phone')"
      >
        <svg
          viewBox="0 0 24 24"
          width="20"
          height="20"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          aria-hidden="true"
        >
          <rect x="6" y="2" width="12" height="20" rx="2" />
          <path d="M10 18h4" />
        </svg>
        Continue with phone number
      </button>
      <div class="login-divider"><span>or use a provider</span></div>
      <button
        class="provider-google"
        type="button"
        data-method="google"
        aria-label="Sign in with Google"
        @click="choose('google')"
      >
        <img src="/google_g_logo.svg" alt="" width="28" height="28" />
        Continue with Google
      </button>
      <button
        class="provider-facebook"
        type="button"
        data-method="facebook"
        @click="choose('facebook')"
      >
        <img
          src="/facebook_icon.svg"
          alt=""
          width="24"
          height="24"
        />
        Continue with Facebook
      </button>
    </div>
    <section
      v-else
      ref="preview"
      class="login-preview"
      tabindex="-1"
      aria-labelledby="demo-method-title"
    >
      <h2 id="demo-method-title">
        {{ labels[selected] }} {{ signup ? 'sign-up' : 'login' }} preview
      </h2>
      <p>
        Authentication will be available later. You can explore Profile Settings
        and sample prayer requests using a temporary demo session.
      </p>
      <button
        class="button primary"
        type="button"
        @click="emit('enterDemo', selected)"
      >
        Enter demo session
      </button>
      <button class="button secondary" type="button" @click="back">
        Choose another option
      </button>
    </section>
    <p class="login-mode">
      {{ signup ? 'Already have an account?' : 'New to Open Monastery?' }}
      <a :href="signup ? '#login' : '#signup'">{{
        signup ? 'Log in' : 'Sign up'
      }}</a>
    </p>
  </section>
</template>
