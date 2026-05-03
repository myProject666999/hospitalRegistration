<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="6" v-for="item in statsCards" :key="item.title">
        <el-card shadow="hover" class="stat-card" :style="{ background: item.color }">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon :size="48"><component :is="item.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-number">{{ item.value }}</p>
              <p class="stat-title">{{ item.title }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>预约趋势</span>
          </template>
          <div ref="appointmentChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>科室预约分布</span>
          </template>
          <div ref="departmentChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>最新预约</span>
          </template>
          <el-table :data="recentAppointments" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="user_name" label="用户" />
            <el-table-column prop="doctor_name" label="医生" />
            <el-table-column prop="department_name" label="科室" />
            <el-table-column prop="date" label="日期" />
            <el-table-column prop="time_slot" label="时间段">
              <template #default="scope">
                {{ scope.row.time_slot === 'morning' ? '上午' : scope.row.time_slot === 'afternoon' ? '下午' : '晚间' }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { getStatistics } from '@/api/auth'
import { 
  User, 
  UserFilled, 
  Calendar, 
  Star,
  Document
} from '@element-plus/icons-vue'

const appointmentChartRef = ref(null)
const departmentChartRef = ref(null)

let appointmentChart = null
let departmentChart = null

const statsCards = ref([
  { title: '总用户数', value: '0', icon: 'UserFilled', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { title: '总医生数', value: '0', icon: 'User', color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { title: '总预约数', value: '0', icon: 'Calendar', color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
  { title: '平均好评率', value: '0%', icon: 'Star', color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' }
])

const recentAppointments = ref([])

const getStatusType = (status) => {
  const types = ['info', 'warning', 'success', 'danger']
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = ['待确认', '已确认', '已完成', '已取消']
  return texts[status] || '未知'
}

const loadStatistics = async () => {
  try {
    const res = await getStatistics()
    if (res.code === 200) {
      const data = res.data
      statsCards.value[0].value = data.total_users || 0
      statsCards.value[1].value = data.total_doctors || 0
      statsCards.value[2].value = data.total_appointments || 0
      statsCards.value[3].value = data.avg_rating ? (data.avg_rating * 20).toFixed(1) + '%' : '0%'
    }
  } catch (error) {
    console.error('加载统计数据失败', error)
  }
}

const initCharts = () => {
  appointmentChart = echarts.init(appointmentChartRef.value)
  departmentChart = echarts.init(departmentChartRef.value)

  const appointmentOption = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      data: [12, 19, 15, 22, 18, 25, 8],
      type: 'line',
      smooth: true,
      areaStyle: {
        opacity: 0.3
      }
    }]
  }

  const departmentOption = {
    tooltip: {
      trigger: 'item'
    },
    legend: {
      bottom: '5%',
      left: 'center'
    },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false
      },
      emphasis: {
        label: {
          show: true,
          fontSize: '14',
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: [
        { value: 35, name: '内科' },
        { value: 25, name: '外科' },
        { value: 20, name: '儿科' },
        { value: 15, name: '妇科' },
        { value: 5, name: '其他' }
      ]
    }]
  }

  appointmentChart.setOption(appointmentOption)
  departmentChart.setOption(departmentOption)
}

const handleResize = () => {
  appointmentChart?.resize()
  departmentChart?.resize()
}

onMounted(() => {
  loadStatistics()
  setTimeout(initCharts, 100)
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  appointmentChart?.dispose()
  departmentChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.dashboard-container {
  padding: 0;
}

.stat-card {
  color: #fff;
  border: none;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  opacity: 0.8;
}

.stat-number {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 5px;
}

.stat-title {
  font-size: 14px;
  opacity: 0.9;
}

.chart-container {
  height: 300px;
}
</style>
