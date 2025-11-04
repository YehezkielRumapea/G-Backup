<script setup>
import { ref, onMounted, computed } from 'vue';
import apiClient from '@/api'; // Gunakan "jembatan" API yang sudah kita buat

// --- STATE UNTUK KONTROL UI (TAMPILAN) ---
const showBackupModal = ref(false);
const showRestoreModal = ref(false);
const isSubmitting = ref(false); // Untuk menonaktifkan tombol saat proses

// --- STATE UNTUK DATA DARI API ---
const remotes = ref([]); // Akan diisi dengan daftar remote dari API
const logs = ref([]);    // Akan diisi dengan log terbaru dari API
const jobs = ref([]);    // Akan diisi dengan daftar job dari API
const isLoading = ref(true); // Status loading data awal
const errorMsg = ref('');  // Pesan error jika gagal mengambil data

// --- STATE UNTUK FORM BACKUP ---
// Strukturnya disesuaikan agar cocok dengan struct `BackupRequest` di backend Go (json tag)
const backupForm = ref({
    job_name: '',
    job_type: 'FILE', // Tipe default: 'FILE' atau 'DB'
    source_path: '', // Untuk FILE: path di server, Untuk DB: akan diisi nama DB

    // Detail koneksi SSH jika menggunakan 2 VM (untuk 1 VM, ini bisa diabaikan oleh backend)
    ssh_host: 'localhost', // <-- Karena pakai 1 VM, kita set ke localhost
    ssh_user: 'idven',     // <-- Sesuaikan dengan username Ubuntu Anda

    // Detail Remote Rclone & Tujuan
    remote_name: '', // Akan diisi dari API
    destination_path: '/gbackup_files/', // Default path tujuan di Cloud

    // Penjadwalan (kosongkan untuk manual)
    schedule_cron: '',

    // Kredensial DB (jika job_type = 'DB')
    db_name: '', // Hanya untuk input di UI, nilainya akan disalin ke source_path
    db_user: '',
    db_pass: '',
});

// --- STATE UNTUK FORM RESTORE ---
// Strukturnya disesuaikan dengan struct `RestoreRequest` di backend Go
const restoreForm = ref({
    job_name: 'Manual Restore',
    remote_name: '',
    source_path: '', // Path file backup di Google Drive
    restore_path: '', // Path tujuan restore di server atau nama DB tujuan
    job_type: 'FILE',
    destDbName: '', // Hanya untuk input di UI
    db_user: '',    // Kredensial DB server tujuan
    db_pass: ''
});

// --- FUNGSI UNTUK MENGAMBIL DATA AWAL DARI BACKEND ---
const fetchData = async () => {
    isLoading.value = true;
    errorMsg.value = '';
    try {
        // Panggil beberapa endpoint API sekaligus untuk efisiensi
        const [remotesRes, logsRes, jobsRes] = await Promise.all([
            apiClient.get('/monitoring/remotes'), // Panggil endpoint remotes
            apiClient.get('/monitoring/logs'),    // Panggil endpoint logs
            apiClient.get('/jobs/scheduled')      // Panggil endpoint jobs
        ]);

        // Isi state dengan data dari API (pastikan data tidak null)
// Isi state dengan data dari API (pastikan data tidak null)
remotes.value = Array.isArray(remotesRes.data) ? remotesRes.data : [];

// --- BLOK YANG DIPERBAIKI ---
let receivedLogs = logsRes.data;
// 1. Periksa dulu apakah data yang diterima adalah sebuah array
if (Array.isArray(receivedLogs)) {
    // 2. Jika ya, baru lakukan pengurutan dan pemotongan
    logs.value = receivedLogs.sort((a, b) => new Date(b.Timestamp) - new Date(a.Timestamp)).slice(0, 5);
} else {
    // 3. Jika bukan, jadikan saja array kosong untuk mencegah error
    console.warn("Data log yang diterima dari backend bukanlah sebuah array:", receivedLogs);
    logs.value = [];
}
// --- AKHIR BLOK YANG DIPERBAIKI ---

jobs.value = Array.isArray(jobsRes.data) ? jobsRes.data : [];

        // Set remote pertama sebagai pilihan default di form jika belum dipilih
        if (remotes.value.length > 0) {
            if (!backupForm.value.remote_name) backupForm.value.remote_name = remotes.value[0].RemoteName;
            if (!restoreForm.value.remote_name) restoreForm.value.remote_name = remotes.value[0].RemoteName;
        }

    } catch (error) {
        console.error("Gagal memuat data dashboard:", error);
        if (error.response?.status !== 401) { // Jangan tampilkan error jika hanya karena token expired
            errorMsg.value = "Gagal memuat data dari server. Pastikan backend berjalan dan terhubung.";
        }
    } finally {
        isLoading.value = false;
    }
};

// Panggil fetchData() saat komponen pertama kali dimuat
onMounted(fetchData);

// --- FUNGSI UTAMA: MENGIRIM JOB BARU KE BACKEND ---
const saveBackupJob = async () => {
    isSubmitting.value = true;

    // Siapkan payload (data yang akan dikirim), sesuai struktur `BackupRequest` di Go
    let payload = { ...backupForm.value };

    // **Logika "Satu Nalar"**:
    // Jika tipe backup adalah 'DB', field 'source_path' di backend Go
    // seharusnya berisi nama database, bukan path file. Kita salin dari 'db_name'.
    if (payload.job_type === 'DB') {
        payload.source_path = payload.db_name;
    }

    // Hapus field yang hanya untuk UI (db_name) agar tidak dikirim ke backend
    delete payload.db_name;

    console.log("Mengirim payload backup ke /jobs/new:", payload);

    try {
        // Kirim data ke endpoint backend Go
        const response = await apiClient.post('/jobs/new', payload);

        alert(response.data.message || "Job backup manual berhasil dimulai!");
        closeModal();
        fetchData(); // Muat ulang data (terutama log) di dashboard

    } catch (error) {
        console.error("Error saat menyimpan job backup:", error.response?.data || error.message);
        alert("Gagal membuat job backup: " + (error.response?.data?.error || "Error tidak diketahui"));
    } finally {
        isSubmitting.value = false;
    }
};

// --- FUNGSI UNTUK MENGIRIM PERINTAH RESTORE KE BACKEND ---
const startRestore = async () => {
    isSubmitting.value = true;
    
    let payload = { ...restoreForm.value };
    
    // Logika "Satu Nalar": Jika tipe DB, `restore_path` adalah nama DB tujuan
    if (payload.job_type === 'DB') {
        payload.restore_path = payload.destDbName;
    }
    
    delete payload.destDbName;
    
    console.log("Mengirim payload restore ke /restore:", payload);

    try {
        const response = await apiClient.post('/jobs/restore', payload); // Sesuaikan endpoint jika berbeda
        alert(response.data.message || "Job restore berhasil dimulai!");
        closeModal();
        fetchData();
    } catch (error) {
        console.error("Error saat memulai restore:", error.response?.data || error.message);
        alert("Gagal memulai restore: " + (error.response?.data?.error || "Error tidak diketahui"));
    } finally {
        isSubmitting.value = false;
    }
};

// --- FUNGSI UTILITAS UI ---
const closeModal = () => {
  showBackupModal.value = false;
  showRestoreModal.value = false;
};
const formatTimestamp = (timestamp) => {
    if (!timestamp) return '-';
    try { return new Date(timestamp).toLocaleString('id-ID', { hour12: false }); } catch(e) { return timestamp; }
};

// --- COMPUTED PROPERTIES UNTUK STATISTIK DI DASHBOARD ---
const totalRemotes = computed(() => remotes.value.length);
const totalJobs = computed(() => jobs.value.length);
const nextJob = computed(() => {
    if (!Array.isArray(jobs.value)) return null; // Penjaga jika jobs bukan array
    const futureJobs = jobs.value
        .map(j => ({ ...j, nextRunDate: new Date(j.nextRun) }))
        .filter(j => j.nextRunDate && j.nextRunDate > new Date())
        .sort((a, b) => a.nextRunDate - b.nextRunDate);
    return futureJobs.length > 0 ? futureJobs[0] : null;
});
</script>

<template>
    <div class="dashboard-container">
        <header class="main-header">
            <h1>Dashboard</h1>
        </header>

        <div v-if="isLoading" class="text-center p-8">Memuat data dashboard...</div>
        <div v-else-if="errorMsg" class="error-message card">
            <p><strong>Error:</strong> {{ errorMsg }}</p>
            <button @click="fetchData" class="btn btn-primary" style="margin-top: 10px;">Coba Lagi</button>
        </div>

        <div v-else class="dashboard-grid">
            <div class="left-column">
                <div class="top-stats">
                    <div class="stat-card">
                        <h3>Total Remote</h3>
                        <p>{{ totalRemotes }}</p>
                    </div>
                    <div class="stat-card">
                        <h3>Total Job</h3>
                        <p>{{ totalJobs }}</p>
                    </div>
                </div>

                <div class="card quick-actions">
                    <h3>Quick Action</h3>
                    <div class="action-buttons">
                        <button @click="showBackupModal = true" class="action-btn backup-btn">Create Manual Backup</button>
                        <button @click="showRestoreModal = true" class="action-btn restore-btn">Start Manual Restore</button>
                    </div>
                </div>

                <div class="card next-job">
                    <h3>Next Scheduled Job</h3>
                    <div v-if="nextJob">
                        <p><strong>{{ nextJob.name }}</strong> ({{ nextJob.jobType }})</p>
                        <p>Scheduled for: <strong>{{ formatTimestamp(nextJob.nextRun) }}</strong></p>
                    </div>
                     <div v-else>
                        <p>Tidak ada job terjadwal berikutnya.</p>
                    </div>
                </div>
            </div>

            <div class="right-column">
                <div class="card log-activity">
                    <h3>Recent Log Activity</h3>
                    <ul v-if="logs.length > 0">
                        <li v-for="log in logs" :key="log.ID">
                            <span class="log-time">{{ formatTimestamp(log.Timestamp) }}</span>
                            <span class="log-message">
                                <strong :class="{ 'text-success': log.Status === 'SUCCESS', 'text-danger': log.Status !== 'SUCCESS' }">
                                    [{{ log.OperationType }} - {{ log.Status }}]
                                </strong>
                                {{ log.Message }}
                            </span>
                        </li>
                    </ul>
                     <p v-else>Belum ada aktivitas log.</p>
                </div>
            </div>
        </div>

        <div v-if="showBackupModal" class="modal-overlay" @click.self="closeModal">
            <div class="modal-content">
                <div class="modal-header">
                    <h2>Manual Backup Configuration</h2>
                    <button @click="closeModal" class="close-button">&times;</button>
                </div>
                <form @submit.prevent="saveBackupJob">
                    <div class="form-group">
                        <label for="backupJobName">Job Name*</label>
                        <input type="text" id="backupJobName" v-model="backupForm.job_name" required>
                    </div>

                    <div class="form-group">
                        <label for="backupJobType">Target Type</label>
                        <select id="backupJobType" v-model="backupForm.job_type">
                            <option value="FILE">File / Folder</option>
                            <option value="DB">Database (MariaDB/MySQL)</option>
                        </select>
                    </div>

                    <div v-if="backupForm.job_type === 'FILE'">
                        <div class="form-group">
                            <label for="backupSourcePathFile">Source Path (Lokal di Server)*</label>
                            <input type="text" id="backupSourcePathFile" v-model="backupForm.source_path" required placeholder="/home/idven/dokumen">
                        </div>
                    </div>

                    <div v-if="backupForm.job_type === 'DB'" class="db-group">
                        <div class="form-group">
                            <label for="backupDbName">DB Name*</label>
                            <input type="text" id="backupDbName" v-model="backupForm.db_name" required>
                        </div>
                        <div class="form-group">
                            <label for="backupDbUser">DB Username*</label>
                            <input type="text" id="backupDbUser" v-model="backupForm.db_user" required>
                        </div>
                         <div class="form-group">
                            <label for="backupDbPass">DB Password</label>
                            <input type="password" id="backupDbPass" v-model="backupForm.db_pass">
                        </div>
                    </div>

                    <div class="form-group">
                        <label for="backupRemoteName">Cloud Remote Destination*</label>
                        <select id="backupRemoteName" v-model="backupForm.remote_name" required>
                            <option disabled value="">-- Pilih Remote --</option>
                            <option v-for="remote in remotes" :key="remote.RemoteName" :value="remote.RemoteName">
                                {{ remote.RemoteName }}
                            </option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="backupDestPath">Cloud Destination Path*</label>
                        <input type="text" id="backupDestPath" v-model="backupForm.destination_path" required placeholder="/folder-backup/">
                    </div>

                    <div class="modal-actions">
                         <button type="button" class="btn btn-secondary" @click="closeModal" :disabled="isSubmitting">Cancel</button>
                        <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
                            {{ isSubmitting ? 'Starting Job...' : 'Start Manual Backup' }}
                        </button>
                    </div>
                </form>
            </div>
        </div>

        <div v-if="showRestoreModal" class="modal-overlay" @click.self="closeModal">
            <div class="modal-content">
                 <div class="modal-header">
                    <h2>Manual Restore Configuration</h2>
                    <button @click="closeModal" class="close-button">&times;</button>
                </div>
                <form @submit.prevent="startRestore">
                     <div class="form-group">
                        <label for="restoreRemoteName">Cloud Remote Source*</label>
                         <select id="restoreRemoteName" v-model="restoreForm.remote_name" required>
                             <option disabled value="">-- Pilih Remote --</option>
                             <option v-for="remote in remotes" :key="remote.RemoteName" :value="remote.RemoteName">
                                {{ remote.RemoteName }}
                            </option>
                        </select>
                    </div>
                     <div class="form-group">
                        <label for="restoreSourcePath">Source Path (di Cloud)*</label>
                        <input type="text" id="restoreSourcePath" v-model="restoreForm.source_path" required placeholder="/folder-backup/namafile.zip">
                    </div>

                     <div class="form-group">
                        <label for="restoreJobType">Restore Target Type</label>
                        <select id="restoreJobType" v-model="restoreForm.job_type">
                            <option value="FILE">File / Folder</option>
                            <option value="DB">Database</option>
                        </select>
                    </div>

                    <div v-if="restoreForm.job_type === 'FILE'">
                        <div class="form-group">
                            <label for="restoreDestPathFile">Destination Path (Lokal di Server)*</label>
                            <input type="text" id="restoreDestPathFile" v-model="restoreForm.restore_path" required placeholder="/home/idven/restored_files/">
                        </div>
                    </div>

                    <div v-if="restoreForm.job_type === 'DB'" class="db-group">
                         <div class="form-group">
                            <label for="restoreDestDbName">Destination DB Name*</label>
                            <input type="text" id="restoreDestDbName" v-model="restoreForm.destDbName" required>
                        </div>
                        <div class="form-group">
                            <label for="restoreDbUser">DB Username (Lokal)*</label>
                            <input type="text" id="restoreDbUser" v-model="restoreForm.db_user" required>
                        </div>
                         <div class="form-group">
                            <label for="restoreDbPass">DB Password (Lokal)</label>
                            <input type="password" id="restoreDbPass" v-model="restoreForm.db_pass">
                        </div>
                    </div>

                    <div class="modal-actions">
                        <button type="button" @click="closeModal" class="btn btn-secondary" :disabled="isSubmitting">Cancel</button>
                        <button type="submit" class="btn btn-success" :disabled="isSubmitting">
                             {{ isSubmitting ? 'Starting Restore...' : 'Start Manual Restore' }}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* --- STYLE DASHBOARD --- */
.main-header { margin-bottom: 20px; }
.main-header h1 { font-size: 2rem; font-weight: 700; color: #333; }
.dashboard-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
.left-column, .right-column { display: flex; flex-direction: column; gap: 20px; }
.card { background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); border: 1px solid var(--border-color); }
.card h3 { margin-top: 0; margin-bottom: 15px; font-size: 1.1rem; font-weight: 600; border-bottom: 1px solid #eee; padding-bottom: 10px; }
.top-stats { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.stat-card { background-color: #fff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); border: 1px solid var(--border-color); text-align: center;}
.stat-card h3 { border-bottom: none; font-size: 1rem; color: #6c757d; margin-bottom: 5px; }
.stat-card p { font-size: 2rem; font-weight: 700; color: var(--primary-color); margin: 0; }
.quick-actions .action-buttons { display: grid; grid-template-columns: 1fr 1fr; gap: 15px; }
.action-btn { padding: 20px 10px; border: none; border-radius: 8px; cursor: pointer; font-size: 1rem; font-weight: 600; color: #fff; text-align: center; }
.backup-btn { background-color: var(--primary-color); }
.backup-btn:hover { background-color: var(--primary-hover); }
.restore-btn { background-color: var(--success-color); }
.restore-btn:hover { background-color: #1a6e3a; }
.next-job p { font-size: 0.95rem; color: #333; line-height: 1.5; margin: 5px 0;}
.log-activity { flex-grow: 1; }
.log-activity ul { list-style: none; padding: 0; margin: 0; max-height: 350px; overflow-y: auto;}
.log-activity li { padding: 8px 0; border-bottom: 1px solid #f0f0f0; display: flex; flex-direction: column; gap: 2px; font-size: 0.9rem;}
.log-activity li:last-child { border-bottom: none; }
.log-time { color: #6c757d; font-size: 0.8rem; }
.log-message { color: #343a40; word-break: break-word; }
.log-message .text-success { color: var(--success-color); font-weight: bold; }
.log-message .text-danger { color: var(--danger-color); font-weight: bold; }
@media (max-width: 992px) { .dashboard-grid { grid-template-columns: 1fr; } }
.error-message { color: #dc3545; }
.text-center { text-align: center; }

/* --- STYLE MODAL --- */
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.6); display: flex; justify-content: center; align-items: center; z-index: 1000; padding: 20px; }
.modal-content { background-color: #fff; border-radius: 8px; box-shadow: 0 5px 15px rgba(0,0,0,0.3); width: 90%; max-width: 550px; max-height: 90vh; display: flex; flex-direction: column; }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 15px 20px; border-bottom: 1px solid #eee; }
.modal-header h2 { margin: 0; font-size: 1.25rem; }
.close-button { background: none; border: none; font-size: 1.7rem; cursor: pointer; color: #888; line-height: 1; }
.modal-content form { padding: 20px; overflow-y: auto; display: flex; flex-direction: column; gap: 15px; }
.form-group label { display: block; margin-bottom: 5px; font-weight: 500; font-size: 0.9rem; }
.form-group input, .form-group select { width: 100%; padding: 10px; border: 1px solid #ccc; border-radius: 4px; font-size: 1rem; box-sizing: border-box; }
.db-group { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 15px; }
.modal-actions { margin-top: 20px; display: flex; justify-content: flex-end; gap: 10px; padding: 15px 20px; border-top: 1px solid #eee; }
.btn-primary, .btn-secondary, .btn-success { padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; font-weight: 500; font-size: 0.9rem; transition: background-color 0.2s, opacity 0.2s;}
.btn-primary:hover:not(:disabled) { background-color: var(--primary-hover); }
.btn-success:hover:not(:disabled) { background-color: #1a6e3a; }
.btn-secondary { background-color: #6c757d; color: white; }
.btn-secondary:hover:not(:disabled) { background-color: #5a6268; }
button:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
