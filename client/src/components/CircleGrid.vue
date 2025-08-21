<template>
  <v-container>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>
    
    <div v-if="!loading && circles.length === 0" class="text-center py-8">
      <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-account-group</v-icon>
      <h3 class="text-grey-lighten-1 mb-2">No circles found</h3>
      <p class="text-grey-lighten-1">Try switching to a different tab or check back later.</p>
    </div>

    <v-row v-if="!loading && circles.length > 0">
      <v-col class="pa-1" md="3" sm="4" cols="6" v-for="circle in circles" :key="circle.name">
        <v-card
          :to="'/'+circle.name"
          style="aspect-ratio: 8/6;border-color: lightgrey;border-width: 1.5px;border-style: solid;"
          hover
          class="circle-card"
        >
          <v-card-title style="font-size: 1rem;">
            {{ circle.title }}
            <v-icon 
                v-if="circle.favorited"
                size="24" 
                class="favorite-heart"
              >
              mdi-heart
            </v-icon>
          </v-card-title>
          <v-card-subtitle style="font-size: 0.8rem;">
            <div v-if="circle.description" class="text-body-2 mb-1" style="max-height: 2.5em; overflow: hidden; text-overflow: ellipsis; white-space: pre-line;">
              {{ circle.description.length > 80 ? circle.description.slice(0, 80) + '…' : circle.description }}
            </div>
          </v-card-subtitle>
          <v-img
            class="mt-2"
            style="background-color: lightgray"
            height="100%"
            :src="circle.imageUri"
            cover
          >
            <template v-slot:placeholder>
              <v-row class="fill-height ma-0" align="center" justify="center">
                <v-icon size="48" color="grey-lighten-1">mdi-account-group</v-icon>
              </v-row>
            </template>
          </v-img>
          
          <!-- Permission level indicator -->
          <v-chip
            v-if="circle.circleAccess?.permissionLevel"
            size="small"
            :color="getPermissionColor(circle.circleAccess?.permissionLevel)"
            class="permission-chip"
          >
            {{ getPermissionText(circle.circleAccess?.permissionLevel) }}
          </v-chip>
        </v-card>
        <v-btn
          v-if="circle.circleAccess?.acceptTarget === 'ACCEPT_TARGET_RECIPIENT' && circle.circleAccess?.state === 'ACCESS_STATE_PENDING'"
          color="success"
          class="accept-btn"
          @click.stop.prevent="$emit('accept', circle)"
          block
          :loading="acceptingCircleId === circle.name"
        >
          Accept
        </v-btn>
        <v-btn
          v-if="circle.circleAccess?.acceptTarget === 'ACCEPT_TARGET_RECIPIENT' && circle.circleAccess?.state === 'ACCESS_STATE_PENDING'"
          color="error"
          class="decline-btn"
          @click.stop.prevent="$emit('decline', circle)"
          block
        >
          Decline
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import type { Circle } from '@/genapi/api/circles/circle/v1alpha1'
import { defineEmits } from 'vue'

interface Props {
  circles: Circle[]
  loading?: boolean
  acceptingCircleId?: string | null
}

defineProps<Props>()
defineEmits(['accept', 'decline'])

function getPermissionColor(permission: string) {
  switch (permission) {
    case 'PERMISSION_LEVEL_ADMIN':
      return 'primary'
    case 'PERMISSION_LEVEL_WRITE':
      return 'orange'
    case 'PERMISSION_LEVEL_READ':
      return 'blue-grey'
    case 'PERMISSION_LEVEL_PUBLIC':
      return 'green'
    default:
      return 'grey'
  }
}

function getPermissionText(permission: string) {
  switch (permission) {
    case 'PERMISSION_LEVEL_ADMIN':
      return 'Admin'
    case 'PERMISSION_LEVEL_WRITE':
      return 'Write'
    case 'PERMISSION_LEVEL_READ':
      return 'Read'
    case 'PERMISSION_LEVEL_PUBLIC':
      return 'Public'
    default:
      return permission.replace('PERMISSION_LEVEL_', '').toLowerCase()
  }
}
</script>

<style scoped>
.circle-card {
  transition: transform 0.2s ease-in-out;
  position: relative;
}

.circle-card:hover {
  transform: translateY(-4px);
}

.permission-chip {
  position: absolute;
  bottom: 8px;
  right: 8px;
  z-index: 1;
}

.accept-btn {
  margin-top: 12px;
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