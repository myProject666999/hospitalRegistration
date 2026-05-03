<template>
  <div class="projects-container">
    <div class="search-bar">
      <el-input v-model="keyword" placeholder="搜索诊疗项目" clearable style="width: 300px" @keyup.enter="handleSearch">
        <template #append>
          <el-button @click="handleSearch"><el-icon><Search /></el-icon></el-button>
        </template>
      </el-input>
    </div>

    <el-table :data="projects" stripe v-loading="loading">
      <el-table-column prop="name" label="项目名称" />
      <el-table-column prop="description" label="项目描述">
        <template #default="scope">
          <span>{{ scope.row.description || '暂无描述' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="price" label="价格">
        <template #default="scope">
          <span class="price">¥{{ scope.row.price || 0 }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="scope">
          <el-button type="primary" link @click="handleFavorite(scope.row)">收藏</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
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
import { getMedicalProjects } from '@/api/public'
import { addFavorite } from '@/api/user'
import { useUserStore } from '@/stores/user'
import { Search } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const projects = ref([])
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

const handleSearch = () => {
  currentPage.value = 1
  loadProjects()
}

const handleSizeChange = () => {
  loadProjects()
}

const handleCurrentChange = () => {
  loadProjects()
}

const loadProjects = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: keyword.value
    }
    const res = await getMedicalProjects(params)
    if (res.code === 200) {
      projects.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载诊疗项目失败', error)
  }
  loading.value = false
}

const handleFavorite = async (project) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  try {
    const res = await addFavorite({
      type: 'medical_project',
      target_id: project.id
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
  loadProjects()
})
</script>

<style scoped>
.projects-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.search-bar {
  margin-bottom: 20px;
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
