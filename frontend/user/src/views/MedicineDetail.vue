<template>
  <div class="medicine-detail-container">
    <el-breadcrumb separator="/" class="breadcrumb">
      <el-breadcrumb-item :to="{ path: '/home' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item :to="{ path: '/medicines' }">药品信息</el-breadcrumb-item>
      <el-breadcrumb-item>药品详情</el-breadcrumb-item>
    </el-breadcrumb>

    <el-card v-loading="loading" class="detail-card">
      <div class="medicine-header">
        <div class="medicine-avatar">
          <el-icon :size="80"><MedicineBox /></el-icon>
        </div>
        <div class="medicine-info">
          <div class="title-row">
            <h2>{{ medicine.name }}</h2>
            <el-tag v-if="medicine.is_prescription" type="danger">处方药</el-tag>
            <el-tag v-else type="success">非处方药</el-tag>
          </div>
          <p class="price">价格：<span>¥{{ medicine.price || 0 }}</span></p>
          <p class="specification">规格：{{ medicine.specification || '暂无' }}</p>
          <p class="manufacturer">生产厂家：{{ medicine.manufacturer || '暂无' }}</p>
        </div>
      </div>
    </el-card>

    <el-card class="section-card">
      <template #header>
        <span>适应症</span>
      </template>
      <p>{{ medicine.indication || '暂无信息' }}</p>
    </el-card>

    <el-card class="section-card">
      <template #header>
        <span>用法用量</span>
      </template>
      <p>{{ medicine.dosage || '暂无信息' }}</p>
    </el-card>

    <el-card class="section-card">
      <template #header>
        <span>不良反应</span>
      </template>
      <p>{{ medicine.adverse_reaction || '暂无信息' }}</p>
    </el-card>

    <el-card class="section-card">
      <template #header>
        <span>禁忌</span>
      </template>
      <p>{{ medicine.contraindication || '暂无信息' }}</p>
    </el-card>

    <el-card class="section-card">
      <template #header>
        <span>注意事项</span>
      </template>
      <p>{{ medicine.precautions || '暂无信息' }}</p>
    </el-card>

    <div class="action-bar">
      <el-button type="primary" :loading="favoriteLoading" @click="handleAddFavorite">
        收藏药品
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getMedicineDetail } from '@/api/public'
import { addFavorite } from '@/api/user'
import { useUserStore } from '@/stores/user'
import { MedicineBox } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const medicineId = route.params.id
const loading = ref(false)
const favoriteLoading = ref(false)
const medicine = ref({})

const loadMedicineDetail = async () => {
  loading.value = true
  try {
    const res = await getMedicineDetail(medicineId)
    if (res.code === 200) {
      medicine.value = res.data || {}
    }
  } catch (error) {
    console.error('加载药品详情失败', error)
  }
  loading.value = false
}

const handleAddFavorite = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  favoriteLoading.value = true
  try {
    const res = await addFavorite({
      type: 'medicine',
      target_id: medicineId
    })
    if (res.code === 200) {
      ElMessage.success('收藏成功')
    } else {
      ElMessage.error(res.message || '收藏失败')
    }
  } catch (error) {
    ElMessage.error('收藏失败')
  }
  favoriteLoading.value = false
}

onMounted(() => {
  loadMedicineDetail()
})
</script>

<style scoped>
.medicine-detail-container {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.breadcrumb {
  margin-bottom: 20px;
}

.detail-card {
  margin-bottom: 20px;
}

.medicine-header {
  display: flex;
  gap: 30px;
}

.medicine-avatar {
  width: 150px;
  height: 150px;
  background: #f0f2f5;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409eff;
}

.medicine-info {
  flex: 1;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.title-row h2 {
  font-size: 24px;
  color: #333;
  margin: 0;
}

.price {
  font-size: 18px;
  margin-bottom: 10px;
}

.price span {
  color: #f56c6c;
  font-weight: bold;
}

.specification,
.manufacturer {
  color: #606266;
  margin-bottom: 5px;
}

.section-card {
  margin-bottom: 20px;
}

.action-bar {
  text-align: center;
  margin-top: 30px;
}
</style>
