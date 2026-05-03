<template>
  <div class="medicines-container">
    <div class="search-bar">
      <el-input v-model="keyword" placeholder="搜索药品名称" clearable style="width: 300px" @keyup.enter="handleSearch">
        <template #append>
          <el-button @click="handleSearch"><el-icon><Search /></el-icon></el-button>
        </template>
      </el-input>
    </div>

    <div class="medicine-list">
      <el-row :gutter="20">
        <el-col :span="8" v-for="medicine in medicines" :key="medicine.id">
          <el-card shadow="hover" class="medicine-card" @click="goToDetail(medicine.id)">
            <div class="medicine-header">
              <h4>{{ medicine.name }}</h4>
              <el-tag v-if="medicine.is_prescription" type="danger">处方药</el-tag>
              <el-tag v-else type="success">非处方药</el-tag>
            </div>
            <div class="medicine-info">
              <p><span>规格：</span>{{ medicine.specification || '暂无' }}</p>
              <p><span>生产厂家：</span>{{ medicine.manufacturer || '暂无' }}</p>
              <p><span>适应症：</span>{{ medicine.indication || '暂无' }}</p>
            </div>
            <div class="medicine-footer">
              <span class="price">¥{{ medicine.price || 0 }}</span>
              <el-button type="primary" size="small" @click.stop="addFavorite(medicine)">
                收藏
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[9, 18, 27]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getMedicines } from '@/api/public'
import { addFavorite } from '@/api/user'
import { useUserStore } from '@/stores/user'
import { Search } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const medicines = ref([])
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(9)
const total = ref(0)

const goToDetail = (id) => {
  router.push(`/medicines/${id}`)
}

const handleSearch = () => {
  currentPage.value = 1
  loadMedicines()
}

const handleSizeChange = () => {
  loadMedicines()
}

const handleCurrentChange = () => {
  loadMedicines()
}

const loadMedicines = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: keyword.value
    }
    const res = await getMedicines(params)
    if (res.code === 200) {
      medicines.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载药品列表失败', error)
  }
}

const handleAddFavorite = async (medicine) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  try {
    const res = await addFavorite({
      type: 'medicine',
      target_id: medicine.id
    })
    if (res.code === 200) {
      ElMessage.success('收藏成功')
    } else {
      ElMessage.error(res.message || '收藏失败')
    }
  } catch (error) {
    ElMessage.error('收藏失败')
  }
}

onMounted(() => {
  loadMedicines()
})
</script>

<style scoped>
.medicines-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.search-bar {
  margin-bottom: 20px;
}

.medicine-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.medicine-card:hover {
  transform: translateY(-5px);
}

.medicine-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.medicine-header h4 {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin: 0;
}

.medicine-info {
  font-size: 13px;
  color: #606266;
}

.medicine-info p {
  margin-bottom: 5px;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.medicine-info span {
  color: #909399;
}

.medicine-footer {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  color: #f56c6c;
  font-size: 18px;
  font-weight: bold;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
