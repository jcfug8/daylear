<template>
  <div>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>
    
    <div v-if="!loading && circles.length === 0" class="text-center py-8">
      <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-account-group</v-icon>
      <h3 class="text-grey-lighten-1 mb-2">No circles found</h3>
      <p class="text-grey-lighten-1">Try switching to a different tab or check back later.</p>
    </div>

    <v-row v-if="!loading && circles.length > 0">
      <v-col lg="3" md="4" sm="6" cols="12" v-for="circle in circles" :key="circle.name">
        <v-card
          :to="{ name: 'circle', params: { circleId: circle.name } }"
          :title="circle.title"
          style="aspect-ratio: 8/6"
          hover
          class="circle-card"
        >
          <v-card-subtitle>
            {{ circle.name }}
          </v-card-subtitle>
          <v-img
            class="mt-4"
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
          v-if="circle.circleAccess?.state === 'ACCESS_STATE_PENDING'"
          color="success"
          class="accept-btn"
          @click.stop.prevent="$emit('accept', circle)"
          block
          :loading="acceptingCircleId === circle.name"
        >
          Accept
        </v-btn>
        <v-btn
          v-if="circle.circleAccess?.state === 'ACCESS_STATE_PENDING'"
          color="error"
          class="decline-btn"
          @click.stop.prevent="$emit('decline', circle)"
          block
        >
          Decline
        </v-btn>
      </v-col>
    </v-row>
  </div>
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
const emit = defineEmits(['accept', 'decline'])

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
  top: 8px;
  right: 8px;
  z-index: 1;
}

.accept-btn {
  margin-top: 12px;
}
</style> 