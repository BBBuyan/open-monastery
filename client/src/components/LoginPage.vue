<script setup lang="ts">
import { onUnmounted, ref } from 'vue'
import { ApiError, logIn, registerUser } from '../api'
import { sessionFromToken, type Session } from '../session'

const props = defineProps<{ signup: boolean; registeredEmail: string }>()
const emit = defineEmits<{
  registered: [email: string]
  signedIn: [session: Session]
}>()
const firstName = ref('')
const lastName = ref('')
const email = ref(props.registeredEmail)
const password = ref('')
const pending = ref(false)
const error = ref('')
let controller: AbortController | undefined

async function submit() {
  if (pending.value) return
  pending.value = true
  error.value = ''
  controller = new AbortController()
  const { signal } = controller
  try {
    if (props.signup) {
      await registerUser(
        {
          first_name: firstName.value.trim(),
          last_name: lastName.value.trim(),
          email: email.value.trim(),
          password: password.value,
        },
        signal,
      )
      if (!signal.aborted) emit('registered', email.value.trim())
    } else {
      const token = await logIn(email.value.trim(), password.value, signal)
      if (!signal.aborted)
        emit('signedIn', sessionFromToken(token, email.value.trim()))
    }
  } catch (failure) {
    if (signal.aborted) return
    const status = failure instanceof ApiError ? failure.status : undefined
    if (status === 401 && !props.signup) {
      error.value = 'The email or password is incorrect.'
    } else if (status === 409 && props.signup) {
      error.value = 'An account already uses this email address.'
    } else if ((status === 400 || status === 422) && props.signup) {
      error.value =
        'Check your name, email address, and password, then try again.'
    } else {
      error.value = props.signup
        ? 'Your account couldn’t be created. Please try again.'
        : 'You couldn’t be signed in. Please try again.'
    }
  } finally {
    pending.value = false
  }
}

onUnmounted(() => controller?.abort())
</script>

<template>
  <section class="login-page" aria-labelledby="login-title">
    <a class="return-link" href="#readings">← Browse monasteries</a>
    <h1 id="login-title">{{ signup ? 'Create an account' : 'Log in' }}</h1>
    <p v-if="!signup && registeredEmail" class="login-notice" role="status">
      Your account has been created. Log in with your email and password.
    </p>
    <form class="login-form" :aria-busy="pending" @submit.prevent="submit">
      <fieldset :disabled="pending">
        <div v-if="signup" class="login-name-fields">
          <label>
            First name
            <input
              v-model.trim="firstName"
              name="first_name"
              autocomplete="given-name"
              maxlength="255"
              required
            />
          </label>
          <label>
            Last name
            <input
              v-model.trim="lastName"
              name="last_name"
              autocomplete="family-name"
              maxlength="255"
              required
            />
          </label>
        </div>
        <label>
          Email
          <input
            v-model.trim="email"
            name="email"
            type="email"
            autocomplete="username"
            maxlength="254"
            required
          />
        </label>
        <label>
          Password
          <input
            v-model="password"
            name="password"
            type="password"
            :autocomplete="signup ? 'new-password' : 'current-password'"
            required
          />
        </label>
        <p v-if="error" class="login-notice" role="alert">{{ error }}</p>
        <button class="button primary" type="submit">
          {{
            signup
              ? pending
                ? 'Creating account…'
                : 'Create account'
              : pending
                ? 'Logging in…'
                : 'Log in'
          }}
        </button>
      </fieldset>
    </form>
    <p class="login-mode">
      {{ signup ? 'Already have an account?' : 'New to Open Monastery?' }}
      <a :href="signup ? '#login' : '#signup'">{{
        signup ? 'Log in' : 'Sign up'
      }}</a>
    </p>
  </section>
</template>
