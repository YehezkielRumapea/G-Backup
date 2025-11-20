<template>
    <tr>
        <td>
            <strong>{{ remote.remote_name }}</strong>
            <div v-if="remote.system_message" class="system-warning">
                ⚠️ {{ remote.system_message }}
            </div>
        </td>
        
        <td>
            <span class="status" :class="remote.status_connect.toLowerCase()">
                {{ remote.status_connect }}
            </span>
        </td>

        <td>
            {{ storageDisplay }}
            
            <div class="storage-bar-container">
                <div 
                    class="storage-progress" 
                    :style="{ width: usagePercentage + '%' }" 
                >
                </div>
            </div>
        </td>
        
        <td>
            <span :class="{'job-count-high': remote.active_job_count > 0}">
                {{ remote.active_job_count }}
            </span>
        </td>
        
        <td>{{ formatLastChecked(remote.last_checked_at) }}</td>
    </tr>
</template>

<script setup>
import { computed } from 'vue';

// WAJIB: Menerima objek 'remote' dari parent (Remotes.vue)
const props = defineProps({
    remote: {
        type: Object,
        required: true
    }
});

// --- LOGIKA PERHITUNGAN & FORMATTING ---

// 1. Hitungan Persentase Penggunaan
const usagePercentage = computed(() => {
    // Menggunakan nama field snake_case dari JSON
    const used = props.remote.used_storage_gb; 
    const total = props.remote.total_storage_gb;

    if (!total || total === 0 || isNaN(total)) {
        return 0;
    }
    // Menghitung persentase
    return (used / total) * 100;
});

// 2. Format String Tampilan
const storageDisplay = computed(() => {
    const used = props.remote.used_storage_gb;
    const total = props.remote.total_storage_gb;
    const percentage = usagePercentage.value.toFixed(0);

    const usedFormatted = used.toFixed(2);
    const totalFormatted = total.toFixed(2);

    // Solusi untuk display <0.01 GB (agar tampilan 0.00 tidak menyesatkan)
    const displayUsed = (used > 0 && usedFormatted === "0.00") ? " <0.01" : usedFormatted;

    return `${displayUsed} (${percentage}%) Used ${totalFormatted} GB `;
});

// 3. Format Tanggal
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
/* --- STYLING BARIS REMOTE --- */
/* Pastikan styling ini sudah di-copy dari Remotes.vue jika Anda belum memilikinya di sini */

.status {
    padding: 4px 8px;
    border-radius: 12px;
    font-weight: bold;
    font-size: 0.8rem;
    color: white;
}
.status.connected { background-color: #2ecc71; }
.status.disconnected { background-color: #e74c3c; }

.job-count-high {
    font-weight: bold;
    color: #f39c12;
}

.storage-bar-container {
    max-width: 200px; 
    height: 8px;
    background-color: #ffffff;
    border-radius: 4px;
    overflow: hidden;
    margin-top: 4px;
}
.storage-progress {
    height: 100%;
    background-color: #260aa3; 
    transition: width 0.3s ease;
}
</style>