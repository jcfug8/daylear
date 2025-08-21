<template>
  <v-container>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4" />
    <div v-if="!loading && users.length === 0" class="text-center py-8">
      <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-account</v-icon>
      <h3 class="text-grey-lighten-1 mb-2">{{ emptyText }}</h3>
    </div>
    <v-row v-if="!loading && users.length > 0">
      <v-col class="pa-1" md="3" sm="4" cols="6" v-for="user in users" :key="user.name">
        <v-card
          class="user-card"
          hover
          style="aspect-ratio: 8/6;border-color: lightgrey;border-width: 1.5px;border-style: solid;"
          :to="'/'+user.name"
        >
        <v-card-title class="pt-4 pb-1" style="font-size: 1rem;">
            <div>
              <span v-if="user.givenName || user.familyName">
                {{ user.givenName }} {{ user.familyName }}
                <span class="text-grey-darken-1">({{ user.username }})</span>
              </span>
              <span v-else>
                {{ user.username }}
              </span>
            </div>
            <v-icon 
                v-if="user.favorited"
                size="24" 
                class="favorite-heart"
              >
              mdi-heart
              </v-icon>
          </v-card-title>
          <v-card-subtitle style="font-size: 0.8rem;">
            <div v-if="user.bio" class="text-body-2 mb-1" style="max-height: 2.5em; overflow: hidden; text-overflow: ellipsis; white-space: pre-line;">
              {{ user.bio.length > 80 ? user.bio.slice(0, 80) + '…' : user.bio }}
            </div>
          </v-card-subtitle>
          <v-img
            class="mt-2"
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
  </v-container>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import type { User } from '@/genapi/api/users/user/v1alpha1'

defineProps({
  users: { type: Array as () => User[], required: true },
  loading: { type: Boolean, default: false },
  emptyText: { type: String, default: '' },
  showActions: { type: Boolean, default: false },
  acceptingUserId: { type: String, default: null },
})

const authStore = useAuthStore()
const { user: currentUser } = storeToRefs(authStore)

defineEmits(['accept', 'decline'])

function isRequester(user: User) {
  return user.access?.requester === currentUser.value?.name
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

.favorite-heart {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.6));
  border-radius: 50%;
  padding: 4px;
  transition: all 0.2s ease-in-out;
  color: red;
}
</style> 