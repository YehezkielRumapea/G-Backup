<template>
  <div class="donut-chart-container">
    <apexchart
      type="donut"
      height="200"
      :options="chartOptions"
      :series="series"
    ></apexchart>
  </div>
</template>

<script setup>
import { ref, watchEffect, computed } from 'vue';
import VueApexCharts from 'vue3-apexcharts';

// Daftarkan komponen ApexCharts
const apexchart = VueApexCharts;

// Terima data dari Dashboard.vue
const props = defineProps({
  series: {
    type: Array,
    required: true, // Data pemakaian, misal [1.61, 28.39] (Terpakai, Kosong)
  },
  labels: {
    type: Array,
    required: true, // Nama, misal ['Terpakai', 'Kosong']
  },
  // Total Kapasitas (misal 30)
  totalAvailable: {
    type: Number,
    required: true,
  }
});

// State untuk tema (light/dark)
const currentTheme = ref(localStorage.getItem('theme') || 'light');

// Opsi untuk Donut Chart
const chartOptions = computed(() => {
  const totalUsed = props.series[0] || 0;
  const totalAvail = props.totalAvailable || 0;
  const percentage = (totalAvail > 0) ? (totalUsed / totalAvail) * 100 : 0;

  return {
    chart: {
      type: 'donut',
      foreColor: currentTheme.value === 'dark' ? '#e0e0e0' : '#2c3e50'
    },
    theme: {
      mode: currentTheme.value
    },
    stroke: {
      width: 0,
    },
    plotOptions: {
      pie: {
        donut: {
          // Lingkaran lebih tipis (lubang 80%)
          size: '80%', 
          labels: {
            show: true,
            // Teks GB Terpakai (Besar) di tengah
            value: {
                show: true,
                color: currentTheme.value === 'dark' ? '#e0e0e0' : '#2c3e50',
                fontSize: '1.5rem',
                fontWeight: 'bold',
                formatter: () => { 
                    // Tampilkan total GB Terpakai
                    return totalUsed.toFixed(2) + ' GB';
                }
            },
            // Teks Total Kapasitas (Kecil) di tengah
            total: {
              show: true,
              // Tampilkan label "Total Kapasitas"
              label: 'Total Capacity', 
              color: currentTheme.value === 'dark' ? '#a0a0a0' : '#6c757d',
              formatter: () => {
                // Tampilkan total kapasitas yang tersedia
                return totalAvail.toFixed(0) + ' GB'; 
              }
            },
          }
        },
        expandOnClick: false,
      },
    },
    // Pembagian warna (Terpakai vs Kosong)
    labels: props.labels, 
    // Warna: [Ungu (Terpakai), Abu-abu (Kosong)]
    colors: ['#667eea', currentTheme.value === 'dark' ? '#2c2c2c' : '#f0f0f0'], 
    dataLabels: {
      enabled: false, // Matikan label di slice
    },
    legend: {
      show: true, // Tampilkan legenda (Terpakai / Kosong)
      position: 'bottom',
      horizontalAlign: 'center',
    },
    tooltip: {
      enabled: true,
      y: {
        // Tooltip menampilkan GB dan Persentase
        formatter: (val, { seriesIndex, w }) => {
            const total = w.globals.seriesTotals.reduce((a, b) => a + b, 0);
            const percentage = (val / total) * 100;
            return val.toFixed(2) + " GB (" + percentage.toFixed(1) + "%)";
        },
      }
    }
  };
});

// Watcher untuk ganti tema chart (light/dark)
watchEffect(() => {
  const theme = document.documentElement.getAttribute('data-theme') || 'light';
  currentTheme.value = theme;
});
</script>

<style scoped>
.donut-chart-container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
min-height : 200px;
}
</style>