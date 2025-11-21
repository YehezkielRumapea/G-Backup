<template>
    <tr>
        <td>
            <div class="remote-icon"><strong>{{ remote.remote_name }}</strong></div>
            <div v-if="remote.system_message" class="warning">
                {{ remote.system_message }}
            </div>
        </td>
        
        <td>
            <span class="status" :class="remote.status_connect.toLowerCase()">
                {{ remote.status_connect }}
            </span>
        </td>

        <td>
            <div class="storage-info">
                <span>{{ usedFormatted }} / {{ totalFormatted }} GB</span>
                <div class="storage-row">
                    <div class="storage-bar">
                        <div 
                            class="bar-fill" 
                            :style="{ width: usagePercentage + '%' }"
                        ></div>
                    </div>
                    <span class="percentage">{{ usagePercentage.toFixed(0) }}%</span>
                </div>
            </div>
        </td>
        
        <td class="center">{{ remote.active_job_count }}</td>
        
        <td class="text-muted">{{ formatLastChecked(remote.last_checked_at) }}</td>
    </tr>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    remote: {
        type: Object,
        required: true
    }
});

const used = computed(() => props.remote.used_storage_gb || 0);
const total = computed(() => props.remote.total_storage_gb || 0);

const usagePercentage = computed(() => {
    if (!total.value || total.value === 0) return 0;
    return (used.value / total.value) * 100;
});

const usedFormatted = computed(() => {
    if (used.value > 0 && used.value < 0.01) return '<0.01';
    return used.value.toFixed(2);
});

const totalFormatted = computed(() => total.value.toFixed(2));

function formatLastChecked(isoDate) {
    if (!isoDate) return 'N/A';
    try {
        const date = new Date(isoDate);
        return date.toLocaleDateString('en-GB', { 
            year: 'numeric', 
            month: '2-digit', 
            day: '2-digit', 
            hour: '2-digit', 
            minute: '2-digit' 
        });
    } catch (e) {
        return isoDate;
    }
}
</script>

<style scoped>
.status {
    display: inline-block;
    padding: 3px 8px;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    color: white;
}

.status.connected {
    background-color: #27ae60;
}

.status.disconnected {
    background-color: #e74c3c;
}

.warning {
    font-size: 0.75rem;
    color: #e67e22;
    margin-top: 0.25rem;
}

.storage-info {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
}

.storage-info > span:first-child {
    font-size: 0.73rem;
    font-weight: 500;
}

.storage-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.storage-bar {
    width: 150px;
    height: 4px;
    background-color: #ecf0f1;
    border-radius: 2px;
    overflow: hidden;
}

.bar-fill {
    height: 100%;
    background-color: #3498db;
    transition: width 0.3s ease;
}

.percentage {
    font-size: 0.75rem;
    color: #7f8c8d;
    font-weight: 500;
    min-width: 2rem;
}

.center {
    text-align: center;
}

.text-muted {
    color: #95a5a6;
    font-size: 0.85rem;
}

td {
    padding: 0.75rem;
    vertical-align: middle;
}

strong {
    color: #1a1a1a;
}

.remote-icon {
    font-size: 0.90rem;
    margin-bottom: 0.25rem;
}
</style>