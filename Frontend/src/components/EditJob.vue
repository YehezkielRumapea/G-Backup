<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
        @click.self="handleClose"
      >
        <div class="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-indigo-600 to-purple-600 p-6 flex justify-between items-center">
            <h2 class="text-2xl font-bold text-white">
              Edit Job #{{ jobId }}
            </h2>
            <button
              @click="handleClose"
              class="text-white hover:bg-white hover:bg-opacity-20 rounded-full p-2 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Content -->
          <div class="p-6 overflow-y-auto max-h-[calc(90vh-140px)]">
            <!-- Loading State -->
            <div v-if="loading" class="flex items-center justify-center py-12">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
              <span class="ml-3 text-gray-600">Loading job data...</span>
            </div>

            <!-- Error State -->
            <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
              <p class="text-red-800">Error: {{ error }}</p>
              <button
                @click="loadJobData"
                class="mt-2 text-red-600 hover:text-red-800 font-medium"
              >
                Try Again
              </button>
            </div>

            <!-- Form -->
            <form v-else @submit.prevent="handleSubmit" class="space-y-6">
              <!-- Basic Info -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- Job Name -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Job Name <span class="text-red-500">*</span>
                  </label>
                  <input
                    v-model="formData.job_name"
                    type="text"
                    required
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  />
                  <p v-if="isFieldModified('job_name')" class="text-xs text-orange-600 mt-1">
                    ✓ Modified
                  </p>
                </div>

                <!-- Remote Name -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Remote Name <span class="text-red-500">*</span>
                  </label>
                  <input
                    v-model="formData.remote_name"
                    type="text"
                    required
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  />
                  <p v-if="isFieldModified('remote_name')" class="text-xs text-orange-600 mt-1">
                    ✓ Modified
                  </p>
                </div>

                <!-- Operation Mode -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Operation Mode <span class="text-red-500">*</span>
                  </label>
                  <select
                    v-model="formData.operation_mode"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  >
                    <option value="BACKUP">Backup</option>
                    <option value="RESTORE">Restore</option>
                  </select>
                  <p v-if="isFieldModified('operation_mode')" class="text-xs text-orange-600 mt-1">
                    ✓ Modified
                  </p>
                </div>

                <!-- Rclone Mode -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Rclone Mode <span class="text-red-500">*</span>
                  </label>
                  <select
                    v-model="formData.rclone_mode"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  >
                    <option value="copy">Copy</option>
                    <option value="sync">Sync</option>
                  </select>
                  <p v-if="isFieldModified('rclone_mode')" class="text-xs text-orange-600 mt-1">
                    ✓ Modified
                  </p>
                </div>
              </div>

              <!-- Paths -->
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Source Path <span class="text-red-500">*</span>
                  </label>
                  <input
                    v-model="formData.source_path"
                    type="text"
                    placeholder="/path/to/source"
                    required
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  />
                  <p v-if="isFieldModified('source_path')" class="text-xs text-orange-600 mt-1">
                    ✓ Modified
                  </p>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Destination Path <span class="text-red-500">*</span>
                  </label>
                  <input
                    v-model="formData.destination_path"
                    type="text"
                    placeholder="backup/folder"
                    required
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                  />
                  <p v-if="isFieldModified('destination_path')" class="text-xs text-orange-600 mt-1">
                    ✓ Modified
                  </p>
                </div>
              </div>

              <!-- Schedule -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Schedule (Cron Expression)
                  <span class="text-gray-500 font-normal ml-2 text-xs">
                    Leave empty for manual job
                  </span>
                </label>
                <input
                  v-model="formData.schedule_cron"
                  type="text"
                  placeholder="0 2 * * * (daily at 2 AM)"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                />
                <p v-if="isFieldModified('schedule_cron')" class="text-xs text-orange-600 mt-1">
                  ✓ Modified
                </p>
              </div>

              <!-- Scripts -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Pre-Script (Optional)
                  </label>
                  <textarea
                    v-model="formData.pre_script"
                    rows="4"
                    placeholder="echo 'Before backup'"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent font-mono text-sm"
                  />
                  <p v-if="isFieldModified('pre_script')" class="text-xs text-orange-600 mt-1">
                    ✓ Modified
                  </p>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Post-Script (Optional)
                  </label>
                  <textarea
                    v-model="formData.post_script"
                    rows="4"
                    placeholder="echo 'After backup'"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent font-mono text-sm"
                  />
                  <p v-if="isFieldModified('post_script')" class="text-xs text-orange-600 mt-1">
                    ✓ Modified
                  </p>
                </div>
              </div>

              <!-- Is Active -->
              <div class="flex items-center">
                <input
                  v-model="formData.is_active"
                  type="checkbox"
                  id="is_active"
                  class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
                />
                <label for="is_active" class="ml-2 text-sm text-gray-700">
                  Job Active
                </label>
                <span v-if="isFieldModified('is_active')" class="ml-2 text-xs text-orange-600">
                  ✓ Modified
                </span>
              </div>

              <!-- Changed Fields Summary -->
              <div v-if="changedFieldsCount > 0" class="bg-orange-50 border border-orange-200 rounded-lg p-4">
                <p class="text-sm font-medium text-orange-800 mb-2">
                  📝 Changes detected ({{ changedFieldsCount }} fields):
                </p>
                <ul class="text-sm text-orange-700 space-y-1">
                  <li v-for="field in changedFieldsList" :key="field">
                    • {{ formatFieldName(field) }}
                  </li>
                </ul>
              </div>

              <!-- Error Message -->
              <div v-if="submitError" class="bg-red-50 border border-red-200 rounded-lg p-4">
                <p class="text-red-800">{{ submitError }}</p>
              </div>
            </form>
          </div>

          <!-- Footer -->
          <div class="bg-gray-50 px-6 py-4 flex justify-end gap-3">
            <button
              type="button"
              @click="handleClose"
              :disabled="saving"
              class="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-100 transition-colors disabled:opacity-50"
            >
              Cancel
            </button>
            <button
              @click="handleSubmit"
              :disabled="saving || loading || changedFieldsCount === 0"
              class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <div v-if="saving" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              {{ saving ? 'Saving...' : `Save Changes (${changedFieldsCount})` }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { jobAPI } from '@/services/api';

const props = defineProps({
  jobId: {
    type: Number,
    required: true,
  },
  isOpen: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['close', 'success']);

// State
const loading = ref(false);
const saving = ref(false);
const error = ref(null);
const submitError = ref(null);
const originalData = ref(null);

const formData = ref({
  job_name: '',
  operation_mode: 'BACKUP',
  rclone_mode: 'copy',
  source_path: '',
  destination_path: '',
  remote_name: '',
  schedule_cron: '',
  pre_script: '',
  post_script: '',
  is_active: true,
});

// Watch modal open
watch(() => props.isOpen, (newValue) => {
  if (newValue && props.jobId) {
    loadJobData();
  }
});

// Load job data
const loadJobData = async () => {
  try {
    loading.value = true;
    error.value = null;

    const response = await jobAPI.getById(props.jobId);
    const data = response.data;

    // Store original
    originalData.value = { ...data };

    // Pre-fill form
    formData.value = {
      job_name: data.job_name || '',
      operation_mode: data.operation_mode || 'BACKUP',
      rclone_mode: data.rclone_mode || 'copy',
      source_path: data.source_path || '',
      destination_path: data.destination_path || '',
      remote_name: data.remote_name || '',
      schedule_cron: data.schedule_cron || '',
      pre_script: data.pre_script || '',
      post_script: data.post_script || '',
      is_active: data.is_active ?? true,
    };

    loading.value = false;
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Failed to load job data';
    loading.value = false;
  }
};

// Check if field modified
const isFieldModified = (fieldName) => {
  if (!originalData.value) return false;
  return formData.value[fieldName] !== originalData.value[fieldName];
};

// Get changed fields
const getChangedFields = () => {
  if (!originalData.value) return {};

  const changes = {};
  
  Object.keys(formData.value).forEach(key => {
    if (formData.value[key] !== originalData.value[key]) {
      changes[key] = formData.value[key];
    }
  });

  return changes;
};

// Computed
const changedFieldsCount = computed(() => {
  return Object.keys(getChangedFields()).length;
});

const changedFieldsList = computed(() => {
  return Object.keys(getChangedFields());
});

// Format field name
const formatFieldName = (fieldName) => {
  return fieldName.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
};

// Handle submit
const handleSubmit = async () => {
  const changes = getChangedFields();

  if (Object.keys(changes).length === 0) {
    submitError.value = 'No changes detected';
    return;
  }

  try {
    saving.value = true;
    submitError.value = null;

    await jobAPI.update(props.jobId, changes);

    emit('success');
    emit('close');
    
  } catch (err) {
    submitError.value = err.response?.data?.error || err.message || 'Failed to update job';
  } finally {
    saving.value = false;
  }
};

// Handle close
const handleClose = () => {
  if (!saving.value) {
    emit('close');
  }
};
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .bg-white,
.modal-leave-active .bg-white {
  transition: transform 0.3s ease;
}

.modal-enter-from .bg-white,
.modal-leave-to .bg-white {
  transform: scale(0.95);
}
</style>