<script setup>
import { ref } from 'vue';

// Data simulasi untuk tabel Job Status
const jobStatus = ref([
    { id: 1, name: 'Backup_Daily_DB', remote: 'Azure-West', status: 'Completed', progress: '100%', started: '2025-10-06 23:00', duration: '0h 12m' },
    { id: 2, name: 'Archive_Docs_Monthly', remote: 'S3-Archive', status: 'Running', progress: '75%', started: '2025-10-07 01:30', duration: '0h 45m' },
    { id: 3, name: 'Sync_Media_Weekly', remote: 'Local-NAS', status: 'Failed', progress: '15%', started: '2025-10-07 02:05', duration: '0h 03m' },
    { id: 4, name: 'VM_Snapshot_Hourly', remote: 'GCS-East', status: 'Pending', progress: '0%', started: '-', duration: '-' },
]);

// Logika untuk menentukan kelas status
const getStatusClass = (status) => {
    switch (status) {
        case 'Completed': return 'status-completed';
        case 'Running': return 'status-running';
        case 'Failed': return 'status-failed';
        case 'Pending': return 'status-pending';
        default: return '';
    }
}
</script>

<template>
    <div class="dashboard-content">
        <header class="main-header">
            <h1>Dashboard Overview</h1>
        </header>
        <div class="grid-container">
            <!-- Kotak 1: Total Job -->
            <div class="card card-small">
                <h2>Total Job</h2>
                <p class="data-value">120</p>
            </div>
            <!-- Kotak 2: Log Activity Summary -->
            <div class="card card-small">
                <h2>Log Activity</h2>
                <p class="data-value">2,567</p>
            </div>
            <!-- Kotak 3: Storage Used -->
            <div class="card card-small">
                <h2>Storage Used</h2>
                <p class="data-value">4.5 TB</p>
            </div>
            <!-- Kotak 4: Log Activity Detailed (Placeholder) -->
            <div class="card card-large log-activity-chart">
                <h2>Log Activity Chart</h2>
                <p>Grafik aktivitas log akan ditampilkan di sini.</p>
                <div class="chart-placeholder">~ Placeholder Grafik ~</div>
            </div>
            
            <!-- Kotak 5: Job Status Table (menggantikan placeholder) -->
            <div class="card card-large job-status-summary card-table">
                <h2>Recent Job Status</h2>
                <div class="table-responsive">
                    <table class="job-table">
                        <thead>
                            <tr>
                                <th>Job Name</th>
                                <th>Remote</th>
                                <th>Status</th>
                                <th>Progress</th>
                                <th>Started</th>
                                <th>Duration</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="job in jobStatus" :key="job.id">
                                <td>{{ job.name }}</td>
                                <td>{{ job.remote }}</td>
                                <td>
                                    <span :class="['status-badge', getStatusClass(job.status)]">{{ job.status }}</span>
                                </td>
                                <td>{{ job.progress }}</td>
                                <td>{{ job.started }}</td>
                                <td>{{ job.duration }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* --- HEADER KONTEN UTAMA --- */
.main-header {
  background-color: var(--card-bg);
  padding: 15px 20px;
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid var(--border-color);
}

.main-header h1 {
  font-size: 1.8rem;
  color: #333;
}

/* --- GRID CONTAINER (Dashboard) --- */
.grid-container {
  display: grid;
  gap: 20px;
  grid-template-columns: repeat(3, 1fr); /* 3 kolom default */
}

/* --- CARD UMUM --- */
.card {
  background-color: var(--card-bg);
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
  border: 1px solid var(--border-color);
}

.card h2 {
  font-size: 1.2rem;
  color: #555;
  margin-bottom: 15px;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 8px;
  font-weight: 500;
}

.card-small {
  text-align: center;
}

.data-value {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--primary-color);
  margin-top: 10px;
}

.card-large {
  grid-column: span 2; /* Mencakup 2 kolom */
}

/* Placeholder Styling */
.log-activity-chart .chart-placeholder {
    height: 200px;
    border: 1px dashed #ccc;
    background-color: #f9f9f9;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #888;
    margin-top: 15px;
    border-radius: 4px;
}

/* --- TABLE STYLING (Job Status) --- */
.card-table {
    padding: 0; /* Hapus padding agar tabel memenuhi card */
}

.card-table h2 {
    padding: 20px;
    margin: 0;
    border-bottom: 1px solid var(--border-color);
}

.table-responsive {
    overflow-x: auto; /* Untuk responsif pada layar kecil */
    padding: 0 20px 20px 20px;
}

.job-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 10px;
    font-size: 0.9rem;
}

.job-table th, .job-table td {
    padding: 12px 15px;
    text-align: left;
    border-bottom: 1px solid #f0f0f0;
}

.job-table th {
    background-color: #f8f9fa;
    color: #495057;
    font-weight: 600;
    text-transform: uppercase;
}

.job-table tr:hover {
    background-color: #f5f5f5;
}

/* Status Badges */
.status-badge {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 4px;
    font-weight: 600;
    font-size: 0.8rem;
    color: white;
}

.status-completed {
    background-color: #28a745; /* Green */
}

.status-running {
    background-color: #ffc107; /* Yellow */
    color: #333;
}

.status-failed {
    background-color: #dc3545; /* Red */
}

.status-pending {
    background-color: #6c757d; /* Gray */
}

/* Penyesuaian Responsif */
@media (max-width: 1024px) {
    .grid-container {
        grid-template-columns: repeat(2, 1fr);
    }
    .card-large {
        grid-column: span 2;
    }
}
@media (max-width: 600px) {
    .grid-container {
        grid-template-columns: 1fr;
    }
    .card-small, .card-large {
        grid-column: span 1;
    }
}
</style>
