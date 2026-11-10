<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'

const account = ref('history')
const editing = ref(false)
const email = ref('sample@example.com')
const phone = ref('+976 9911 2233')
const draftEmail = ref('')
const draftPhone = ref('')
const notice = ref('')
const emailInput = ref<HTMLInputElement>()
const editButton = ref<HTMLButtonElement>()
const saved = new Map<string, { email: string; phone: string }>()

watch(account, (value) => {
  const contact = saved.get(value)
  email.value =
    contact?.email ?? (value === 'empty' ? '' : 'sample@example.com')
  phone.value = contact?.phone ?? (value === 'empty' ? '' : '+976 9911 2233')
  editing.value = false
  notice.value = ''
  if (location.hash !== '#profile') location.hash = '#profile'
})
async function edit() {
  draftEmail.value = email.value
  draftPhone.value = phone.value
  notice.value = ''
  editing.value = true
  await nextTick()
  emailInput.value?.focus()
}
async function close() {
  editing.value = false
  await nextTick()
  editButton.value?.focus()
}
function save() {
  email.value = draftEmail.value.trim()
  phone.value = draftPhone.value.trim()
  saved.set(account.value, { email: email.value, phone: phone.value })
  notice.value = 'Demo contact details updated for this preview only.'
  void close()
}
</script>

<template>
  <div class="profile-page">
    <a class="return-link" href="#readings">← Back to books</a>
    <header class="profile-heading">
      <span class="profile-avatar" aria-hidden="true"
        ><svg
          viewBox="0 0 32 32"
          width="40"
          height="40"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <circle cx="16" cy="11" r="5" />
          <path d="M6 28v-3a10 10 0 0 1 20 0v3" /></svg
      ></span>
      <div>
        <h1>Profile Settings</h1>
        <p>Manage your email and phone number.</p>
      </div>
    </header>
    <div class="request-demo">
      <p>
        <strong>Demo profile</strong> · Sample information. No account is signed
        in. Changes last only while this profile is open.
      </p>
      <label for="profile-demo-account">Demo account</label>
      <select id="profile-demo-account" v-model="account">
        <option value="history">Sample account · Contact details</option>
        <option value="empty">Sample account · No contact details</option>
      </select>
    </div>
    <section class="profile-contact" aria-labelledby="contact-title">
      <div class="profile-section-heading">
        <h2 id="contact-title">Contact information</h2>
        <button
          v-if="!editing"
          ref="editButton"
          class="button secondary"
          type="button"
          @click="edit"
        >
          Edit contact details
        </button>
      </div>
      <p class="profile-name">Sample user</p>
      <form v-if="editing" class="profile-form" @submit.prevent="save">
        <p id="contact-demo-notice">
          Preview editing your email and phone. These changes won’t update an
          account.
        </p>
        <label for="profile-email">Email <span>(optional)</span></label>
        <input
          id="profile-email"
          ref="emailInput"
          v-model="draftEmail"
          type="email"
          autocomplete="email"
          maxlength="254"
          aria-describedby="contact-demo-notice"
        />
        <label for="profile-phone">Phone number <span>(optional)</span></label>
        <input
          id="profile-phone"
          v-model="draftPhone"
          type="tel"
          autocomplete="tel"
          maxlength="40"
          aria-describedby="contact-demo-notice"
        />
        <div class="profile-form-actions">
          <button class="button primary" type="submit">Save demo changes</button
          ><button class="button secondary" type="button" @click="close">
            Cancel
          </button>
        </div>
      </form>
      <dl v-else class="profile-contact-values">
        <div>
          <dt>Email</dt>
          <dd>{{ email || 'Not added' }}</dd>
        </div>
        <div>
          <dt>Phone number</dt>
          <dd>{{ phone || 'Not added' }}</dd>
        </div>
      </dl>
      <p role="status" class="profile-notice">{{ notice }}</p>
    </section>
  </div>
</template>
