<template>
  <div>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4" />
    <div v-if="!loading && users.length === 0" class="text-center py-8">
      <v-icon size="64" color="grey-lighten-1" class="mb-4">mdi-account</v-icon>
      <h3 class="text-grey-lighten-1 mb-2">{{ emptyText }}</h3>
    </div>
    <v-row v-if="!loading && users.length > 0">
      <v-col lg="3" md="4" sm="6" cols="12" v-for="user in users" :key="user.name">
        <v-card>
          <v-card-title>
            <v-avatar size="32" class="me-2">
              <v-img v-if="user.imageUri" :src="user.imageUri" />
              <v-icon v-else>mdi-account</v-icon>
            </v-avatar>
            {{ user.username || user.name }}
          </v-card-title>
          <v-card-subtitle>
            {{ user.givenName }} {{ user.familyName }}
          </v-card-subtitle>
          <v-card-text>
            <div v-if="user.bio" class="text-body-2 mb-1" style="max-height: 2.5em; overflow: hidden; text-overflow: ellipsis; white-space: pre-line;">
              {{ user.bio.length > 80 ? user.bio.slice(0, 80) + '…' : user.bio }}
            </div>
            <span v-if="user.visibility">Visibility: {{ user.visibility.replace('VISIBILITY_LEVEL_', '').toLowerCase() }}</span>
          </v-card-text>
          <v-card-actions v-if="showActions">
            <v-btn color="success" block :loading="acceptingUserId === user.name" @click="$emit('accept', user)">Accept</v-btn>
            <v-btn color="error" block @click="$emit('decline', user)">Decline</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  users: { type: Array as () => any[], required: true },
  loading: { type: Boolean, default: false },
  emptyText: { type: String, default: '' },
  showActions: { type: Boolean, default: false },
  acceptingUserId: { type: String, default: null },
})

defineEmits(['accept', 'decline'])
</script> 