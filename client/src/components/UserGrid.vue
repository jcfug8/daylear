<template>
  <div>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4" />
    <div v-if="!loading && users.length === 0" class="text-center py-8">
      <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-account</v-icon>
      <h3 class="text-grey-lighten-1 mb-2">{{ emptyText }}</h3>
    </div>
    <v-row v-if="!loading && users.length > 0">
      <v-col lg="3" md="4" sm="6" cols="12" v-for="user in users" :key="user.name">
        <v-card
          class="user-card"
          hover
          style="aspect-ratio: 8/6"
          :to="canNavigateToUser(user) ? { name: 'user', params: { userId: user.name } } : undefined"
        >
        <v-card-title class="pt-4 pb-1">
            <div>
              <span v-if="user.givenName || user.familyName">
                {{ user.givenName }} {{ user.familyName }}
                <span class="text-grey-darken-1">({{ user.username }})</span>
              </span>
              <span v-else>
                {{ user.username }}
              </span>
            </div>
          </v-card-title>
          <v-card-subtitle>
            <div v-if="user.bio" class="text-body-2 mb-1" style="max-height: 2.5em; overflow: hidden; text-overflow: ellipsis; white-space: pre-line;">
              {{ user.bio.length > 80 ? user.bio.slice(0, 80) + '…' : user.bio }}
            </div>
          </v-card-subtitle>
          <v-img
            class="mt-4"
            style="background-color: lightgray"
            height="100%"
            :src="user.imageUri"
            cover
          >
            <template v-slot:placeholder>
              <v-row class="fill-height ma-0" align="center" justify="center">
                <v-icon size="48" color="grey-lighten-1">mdi-account</v-icon>
              </v-row>
            </template>
          </v-img>
        </v-card>
        <template v-if="showActions && user.access?.state === 'ACCESS_STATE_PENDING'">
          <template v-if="!isRequester(user)">
            <v-btn
              color="success"
              class="accept-btn"
              @click.stop.prevent="$emit('accept', user)"
              block
            >
              Accept
            </v-btn>
            <v-btn
              color="error"
              class="decline-btn"
              @click.stop.prevent="$emit('decline', user)"
              block
            >
              Decline
            </v-btn>
          </template>
          <template v-else>
            <v-btn
              color="warning"
              class="pending-btn"
              block
              disabled
            >
              Pending
            </v-btn>
            <v-btn
              color="error"
              class="decline-btn"
              @click.stop.prevent="$emit('decline', user)"
              block
            >
              Cancel
            </v-btn>
          </template>
        </template>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
const props = defineProps({
  users: { type: Array as () => any[], required: true },
  loading: { type: Boolean, default: false },
  emptyText: { type: String, default: '' },
  showActions: { type: Boolean, default: false },
  acceptingUserId: { type: String, default: null },
})

const authStore = useAuthStore()
const currentUserName = computed(() => authStore.activeAccount?.name)

defineEmits(['accept', 'decline'])

function canNavigateToUser(user: any) {
  return (
    user.visibility === 'VISIBILITY_LEVEL_PUBLIC' ||
    user.visibility === 'VISIBILITY_LEVEL_RESTRICTED' ||
    user.access?.state === 'ACCESS_STATE_ACCEPTED'
  )
}

function isRequester(user: any) {
  return user.access?.requester === currentUserName.value
}
</script>

<style scoped>
.user-card {
  transition: transform 0.2s ease-in-out;
  position: relative;
}
.user-card:hover {
  transform: translateY(-4px);
}
.accept-btn {
  margin-top: 12px;
}
.pending-btn {
  margin-top: 12px;
  cursor: not-allowed;
}
</style> 