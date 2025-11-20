<template>
  <div class="chart-wrapper">
    <div class="donut-container">
      <apexchart
        type="donut"
        height="300"
        :options="donutOptions"
        :series="donutSeries"
      ></apexchart>
      <div class="center-stats">
        <div class="stat-item">
          <span class="stat-label">Used</span>
          <span class="stat-value">{{ totalUsedFormatted }}</span>
        </div>
        <div class="divider"></div>
        <div class="stat-item">
          <span class="stat-label">Free</span>
          <span class="stat-value">{{ totalFreeFormatted }}</span>
        </div>
      </div>
    </div>

    <div class="legend-container">
      <div v-for="(item, index) in driveLegends" :key="index" class="legend-item">
        <span class="legend-dot" :style="{ backgroundColor: item.color }"></span>
        <div class="legend-text">
          <span class="legend-name">{{ item.name }}</span>
          <span class="legend-size">{{ item.used }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import VueApexCharts from "vue3-apexcharts";

const apexchart = VueApexCharts;

const props = defineProps({
  series: { type: Array, required: true },
  labels: { type: Array, required: true },
  colors: { type: Array, required: true }
});

// Hitung total used dan free dari semua drive
const totalUsed = computed(() => {
  let used = 0;
  props.labels.forEach((label, index) => {
    if (!label.includes("(Free)")) {
      used += props.series[index];
    }
  });
  return used;
});

const totalCapacity = computed(() => {
  return props.series.reduce((a, b) => a + b, 0);
});

const totalFree = computed(() => {
  return totalCapacity.value - totalUsed.value;
});

const totalUsedFormatted = computed(() => {
  return formatSize(totalUsed.value);
});

const totalFreeFormatted = computed(() => {
  return formatSize(totalFree.value);
});

// Donut series: [used segments] + [total free]
const donutSeries = computed(() => {
  const series = [];
  props.labels.forEach((label, index) => {
    if (!label.includes("(Free)")) {
      series.push(props.series[index]);
    }
  });
  series.push(totalFree.value);
  return series;
});

// Donut labels
const donutLabels = computed(() => {
  const labels = [];
  props.labels.forEach((label, index) => {
    if (!label.includes("(Free)")) {
      labels.push(label);
    }
  });
  labels.push("Free Space");
  return labels;
});

// Donut colors
const donutColors = computed(() => {
  const colors = [];
  props.labels.forEach((label, index) => {
    if (!label.includes("(Free)")) {
      colors.push(props.colors[index]);
    }
  });
  colors.push("#e8e8e8");
  return colors;
});

// Legend items untuk drive
const driveLegends = computed(() => {
  const legends = [];
  props.labels.forEach((label, index) => {
    if (!label.includes("(Free)")) {
      legends.push({
        name: label,
        used: formatSize(props.series[index]),
        color: props.colors[index]
      });
    }
  });
  return legends;
});

function formatSize(bytes) {
  if (bytes >= 1024) {
    return (bytes / 1024).toFixed(2) + " TB";
  }
  return bytes.toFixed(2) + " GB";
}

const donutOptions = computed(() => ({
  chart: {
    type: "donut",
    foreColor: "#1a1a1a",
    toolbar: { show: false },
    animations: {
      enabled: true,
      speed: 900,
      animateGradually: {
        enabled: true,
        delay: 150
      }
    }
  },

  labels: donutLabels.value,
  colors: donutColors.value,

  legend: {
    show: false
  },

  stroke: {
    width: 2.5,
    colors: ["#ffffff"]
  },

  plotOptions: {
    pie: {
      donut: {
        size: "68%",
        labels: {
          show: false
        }
      },
      expandOnClick: false,
      dataLabels: {
        offset: 40,
        minAngleToShowLabel: 5
      }
    }
  },

  states: {
    hover: {
      filter: { type: "none" }
    },
    active: {
      filter: { type: "none" }
    }
  },

  dataLabels: {
    enabled: true,
    formatter: (val, opts) => {
      const label = donutLabels.value[opts.seriesIndex];
      const size = donutSeries.value[opts.seriesIndex];
      
      const total = donutSeries.value.reduce((a, b) => a + b, 0);
      const percentage = ((size / total) * 100).toFixed(0);
      
      let formattedSize;
      if (size >= 1024) {
        formattedSize = (size / 1024).toFixed(2) + " TB";
      } else {
        formattedSize = size.toFixed(2) + " GB";
      }
      
      return label + "\n" + formattedSize + " (" + percentage + "%)";
    },
    style: {
      fontSize: "12px",
      fontWeight: "600",
      colors: ["#1a1a1a"]
    },
    dropShadow: {
      enabled: false
    },
    connectorColor: "#999",
    connectorPadding: 6
  },

  tooltip: {
    enabled: true,
    theme: "light",
    y: {
      formatter: (v) => {
        return formatSize(v);
      }
    }
  }
}));
</script>

<style scoped>
.chart-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  gap: 10rem;
  padding: 0;
  font-size: 100px;
}

.donut-container {
  position: relative;
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.center-stats {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  align-items: center;
}

.stat-label {
  font-size: 0.75rem;
  color: #999;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a1a;
  line-height: 1;
}

.divider {
  width: 30px;
  height: 1px;
  background: #e5e5e5;
}

/* Legend Container */
.legend-container {
  flex: 0 0 auto;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-height: 300px;
  overflow-y: auto;
  padding-right: 0.5rem;
}

.legend-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.legend-dot {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 0.15rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.legend-text {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  min-width: 0;
}

.legend-name {
  font-size: 1.050rem;
  font-weight: 500;
  color: #1a1a1a;
  word-break: break-word;
}

.legend-size {
  font-size: 0.75rem;
  color: #999;
  font-weight: 400;
  font-family: monospace;
}

/* Scrollbar */
.legend-container::-webkit-scrollbar {
  width: 4px;
}

.legend-container::-webkit-scrollbar-track {
  background: transparent;
}

.legend-container::-webkit-scrollbar-thumb {
  background: #d0d0d0;
  border-radius: 2px;
}

.legend-container::-webkit-scrollbar-thumb:hover {
  background: #999;
}

/* Responsive */
@media (max-width: 1024px) {
  .chart-wrapper {
    flex-direction: column;
    gap: 1.5rem;
  }

  .donut-container {
    height: 280px;
  }

  .legend-container {
    max-height: 200px;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .legend-item {
    flex: 0 0 calc(50% - 0.375rem);
  }
}

@media (max-width: 768px) {
  .chart-wrapper {
    gap: 1rem;
  }

  .donut-container {
    height: 250px;
  }

  .stat-value {
    font-size: 1.25rem;
  }

  .legend-container {
    max-height: 150px;
    gap: 0.5rem;
  }

  .legend-item {
    flex: 0 0 auto;
  }
}
</style>