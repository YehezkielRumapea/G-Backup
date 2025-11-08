<template>
  <tr>
    <td>
      <strong>{{ remote.remote_name }}</strong>
    </td>
    
    <td>
      <span class="status" :class="is_connected ? 'connected' : 'disconnected'">
        {{ remote.status_connect }}
      </span>
    </td>
    
    <td>
      <div v-if="remote.total_storage_gb > 0">
        <span>{{ storageUsage }}</span>
        <progress 
          :value="remote.used_storage_gb" 
          :max="remote.total_storage_gb">
        </progress>
      </div>
      <span v-else>N/A</span>
    </td>
    
    <td>{{ remote.job_runs }}</td>
    
    <td>{{ remote.last_checked_at }}</td>
  </tr>
</template>

<script setup>
import { computed } from 'vue'

// Terima data 'remote' dari parent (Remotes.vue)
const props = defineProps({
  remote: {
    type: Object,
    required: true
  }
})

// Computed property untuk format status
const is_connected = computed(() => {
  return props.remote.status_connect === 'CONNECTED'
})

// Computed property untuk format teks storage
const storageUsage = computed(() => {
  const used = props.remote.used_storage_gb.toFixed(2)
  const total = props.remote.total_storage_gb.toFixed(2)
  const percentage = ((props.remote.used_storage_gb / props.remote.total_storage_gb) * 100).toFixed(0)
  return `${used} dari ${total} GB (${percentage}%) Digunakan`
})
</script>

<style scoped>
.status {
  padding: 4px 8px;
  border-radius: 12px;
  font-weight: bold;
  font-size: 0.8rem;
}
.status.connected {
  background-color: #2ecc71; /* Hijau */
  color: white;
}
.status.disconnected {
  background-color: #e74c3c; /* Merah */
  color: white;
}
progress {
  width: 100%;
}
</style>