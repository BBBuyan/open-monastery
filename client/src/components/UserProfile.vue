<script setup lang="ts">
defineProps<{ email: string }>()
</script>

<template>
  <div class="profile-page">
    <a class="return-link" href="#readings">← Back to books</a>
    <header class="profile-heading">
      <span class="profile-avatar" aria-hidden="true">
        <svg
          viewBox="0 0 32 32"
          width="40"
          height="40"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <circle cx="16" cy="11" r="5" />
          <path d="M6 28v-3a10 10 0 0 1 20 0v3" />
        </svg>
      </span>
      <h1>Profile</h1>
    </header>
    <section class="profile-contact" aria-labelledby="contact-title">
      <h2 id="contact-title">Contact information</h2>
      <dl class="profile-contact-values">
        <div>
          <dt>Email</dt>
          <dd>{{ email }}</dd>
        </div>
      </dl>
    </section>
  </div>
</template>
