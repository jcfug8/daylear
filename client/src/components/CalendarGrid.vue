<template>
  <v-container>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>
    
    <div v-if="!loading && calendars.length === 0" class="text-center py-8">
      <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-calendar</v-icon>
      <h3 class="text-grey-lighten-1 mb-2">No calendars found</h3>
      <p class="text-grey-lighten-1">Try switching to a different tab or check back later.</p>
    </div>

    <v-row v-if="!loading && calendars.length > 0">
      <v-col class="pa-1" md="3" sm="4" cols="6" v-for="calendar in calendars" :key="calendar.name">
        <v-card
          :to="'/'+calendar.name"
          style="border-color: lightgrey;border-width: 1.5px;border-style: solid;"
          hover
          class="calendar-card pb-10"
        >
          <v-card-title style="font-size: 1rem;">
            {{ calendar.title }}
          </v-card-title>
          <v-card-subtitle style="font-size: 0.8rem;">
            <div v-if="calendar.description" class="text-body-2 mb-1" style="max-height: 2.5em; overflow: hidden; text-overflow: ellipsis; white-space: pre-line;">
              {{ calendar.description.length > 80 ? calendar.description.slice(0, 80) + '…' : calendar.description }}
            </div>
          </v-card-subtitle>
          
          <!-- Permission level indicator -->
          <v-chip
            v-if="calendar.calendarAccess?.permissionLevel"
            size="small"
            :color="getPermissionColor(calendar.calendarAccess?.permissionLevel)"
            class="permission-chip"
          >
            {{ getPermissionText(calendar.calendarAccess?.permissionLevel) }}
          </v-chip>
        </v-card>
        <v-btn
          v-if="calendar.calendarAccess?.acceptTarget === 'ACCEPT_TARGET_RECIPIENT' && calendar.calendarAccess?.state === 'ACCESS_STATE_PENDING'"
          color="success"
          class="accept-btn"
          @click.stop.prevent="$emit('accept', calendar)"
          block
          :loading="acceptingCalendarId === calendar.name"
        >
          Accept
        </v-btn>
        <v-btn
          v-if="calendar.calendarAccess?.acceptTarget === 'ACCEPT_TARGET_RECIPIENT' && calendar.calendarAccess?.state === 'ACCESS_STATE_PENDING'"
          color="error"
          class="decline-btn"
          @click.stop.prevent="$emit('decline', calendar)"
          block
        >
          Decline
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import type { Calendar } from '@/genapi/api/calendars/calendar/v1alpha1'
import { defineEmits } from 'vue'

interface Props {
  calendars: Calendar[]
  loading?: boolean
  acceptingCalendarId?: string | null
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
.calendar-card {
  transition: transform 0.2s ease-in-out;
  position: relative;
}

.calendar-card:hover {
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
</style>
