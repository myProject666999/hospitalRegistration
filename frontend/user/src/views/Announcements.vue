<template>
  <div class="announcements-container">
    <el-card class="announcement-card" v-for="announcement in announcements" :key="announcement.id" shadow="hover" @click="goToDetail(announcement.id)">
      <div class="announcement-header">
        <h3>{{ announcement.title }}</h3>
        <span class="date">{{ formatDate(announcement.created_at) }}</span>
      </div>
      <div class="announcement-content">
        <p>{{ announcement.content || '暂无内容' }}</p>
      </div>
      <div class="announcement-footer">
        <el-button type="primary" link>查看详情 <el-icon><ArrowRight /></el-icon></el-button>
      </div>
    </el-card>

    <div class="pagination-container" v-if="total > 0">
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

    <el-empty v-if="announcements.length === 0 && !loading" description="暂无公告信息" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getAnnouncements } from '@/api/public'
import { ArrowRight } from '@element-plus/icons-vue'

const router = useRouter()

const announcements = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString('zh-CN')
}

const goToDetail = (id) => {
  router.push(`/announcements/${id}`)
}

const handleSizeChange = () => {
  loadAnnouncements()
}

const handleCurrentChange = () => {
  loadAnnouncements()
}

const loadAnnouncements = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    const res = await getAnnouncements(params)
    if (res.code === 200) {
      announcements.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载公告列表失败', error)
  }
  loading.value = false
}

onMounted(() => {
  loadAnnouncements()
})
</script>

<style scoped>
.announcements-container {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.announcement-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.announcement-card:hover {
  transform: translateY(-3px);
}

.announcement-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.announcement-header h3 {
  font-size: 18px;
  color: #333;
  margin: 0;
}

.announcement-header .date {
  color: #909399;
  font-size: 14px;
}

.announcement-content {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 15px;
}

.announcement-content p {
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.announcement-footer {
  text-align: right;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
