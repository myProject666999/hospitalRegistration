<template>
  <div class="favorites-container">
    <el-tabs v-model="activeTab" class="favorites-tabs">
      <el-tab-pane label="收藏的医生" name="doctor">
        <div class="favorite-list" v-if="doctorFavorites.length > 0">
          <el-row :gutter="20">
            <el-col :span="8" v-for="item in doctorFavorites" :key="item.id">
              <el-card shadow="hover" class="favorite-card">
                <div class="card-header">
                  <div class="item-info" @click="goToDoctor(item.target_id)">
                    <el-avatar :size="60" class="item-avatar">
                      <el-icon :size="30"><User /></el-icon>
                    </el-avatar>
                    <div class="item-details">
                      <h4>{{ item.target_name || '医生' }}</h4>
                      <p class="info">{{ item.target_info || '医生' }}</p>
                    </div>
                  </div>
                  <el-button type="danger" link icon="Delete" @click="removeFavorite(item)" />
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>
        <el-empty v-else description="暂无收藏的医生" />
      </el-tab-pane>

      <el-tab-pane label="收藏的药品" name="medicine">
        <div class="favorite-list" v-if="medicineFavorites.length > 0">
          <el-row :gutter="20">
            <el-col :span="8" v-for="item in medicineFavorites" :key="item.id">
              <el-card shadow="hover" class="favorite-card">
                <div class="card-header">
                  <div class="item-info" @click="goToMedicine(item.target_id)">
                    <el-avatar :size="60" class="item-avatar medicine-avatar">
                      <el-icon :size="30"><MedicineBox /></el-icon>
                    </el-avatar>
                    <div class="item-details">
                      <h4>{{ item.target_name || '药品' }}</h4>
                      <p class="info">{{ item.target_info || '药品' }}</p>
                    </div>
                  </div>
                  <el-button type="danger" link icon="Delete" @click="removeFavorite(item)" />
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>
        <el-empty v-else description="暂无收藏的药品" />
      </el-tab-pane>

      <el-tab-pane label="收藏的项目" name="medical_project">
        <div class="favorite-list" v-if="projectFavorites.length > 0">
          <el-table :data="projectFavorites" stripe>
            <el-table-column prop="target_name" label="项目名称" />
            <el-table-column prop="target_info" label="项目描述" />
            <el-table-column label="操作" width="120">
              <template #default="scope">
                <el-button type="danger" link @click="removeFavorite(scope.row)">取消收藏</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <el-empty v-else description="暂无收藏的项目" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getFavorites, removeFavorite } from '@/api/user'
import { User, MedicineBox } from '@element-plus/icons-vue'

const router = useRouter()

const activeTab = ref('doctor')
const favorites = ref([])

const doctorFavorites = computed(() => {
  return favorites.value.filter(f => f.type === 'doctor')
})

const medicineFavorites = computed(() => {
  return favorites.value.filter(f => f.type === 'medicine')
})

const projectFavorites = computed(() => {
  return favorites.value.filter(f => f.type === 'medical_project')
})

const loadFavorites = async () => {
  try {
    const res = await getFavorites({ page_size: 100 })
    if (res.code === 200) {
      favorites.value = res.data?.list || []
    }
  } catch (error) {
    console.error('加载收藏列表失败', error)
  }
}

const goToDoctor = (id) => {
  router.push(`/doctors/${id}`)
}

const goToMedicine = (id) => {
  router.push(`/medicines/${id}`)
}

const removeFavorite = async (item) => {
  try {
    await ElMessageBox.confirm('确定要取消收藏吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  const res = await removeFavorite(item.id)
  if (res.code === 200) {
    ElMessage.success('已取消收藏')
    loadFavorites()
  } else {
    ElMessage.error(res.message || '操作失败')
  }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

onMounted(() => {
  loadFavorites()
})
</script>

<style scoped>
.favorites-container {
  padding: 20px;
}

.favorites-tabs {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
}

.favorite-list {
  margin-top: 20px;
}

.favorite-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.item-info {
  display: flex;
  gap: 15px;
  cursor: pointer;
  flex: 1;
}

.item-avatar {
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.medicine-avatar {
  color: #67c23a;
}

.item-details {
  flex: 1;
}

.item-details h4 {
  font-size: 16px;
  color: #333;
  margin-bottom: 5px;
}

.item-details .info {
  color: #909399;
  font-size: 14px;
  margin: 0;
}
</style>
