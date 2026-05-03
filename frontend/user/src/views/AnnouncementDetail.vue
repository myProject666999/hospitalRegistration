<template>
  <div class="announcement-detail-container">
    <el-breadcrumb separator="/" class="breadcrumb">
      <el-breadcrumb-item :to="{ path: '/home' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item :to="{ path: '/announcements' }">公告信息</el-breadcrumb-item>
      <el-breadcrumb-item>公告详情</el-breadcrumb-item>
    </el-breadcrumb>

    <el-card v-loading="loading">
      <div class="announcement-header">
        <h2>{{ announcement.title }}</h2>
        <div class="meta">
          <span>发布时间：{{ formatDate(announcement.created_at) }}</span>
        </div>
      </div>
      <div class="announcement-content">
        <p>{{ announcement.content || '暂无内容' }}</p>
      </div>
    </el-card>

    <div class="back-btn">
      <el-button type="primary" @click="router.back()">返回列表</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getAnnouncementDetail } from '@/api/public'

const route = useRoute()
const router = useRouter()

const announcementId = route.params.id
const loading = ref(false)
const announcement = ref({})

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleString('zh-CN')
}

const loadAnnouncementDetail = async () => {
  loading.value = true
  try {
    const res = await getAnnouncementDetail(announcementId)
    if (res.code === 200) {
      announcement.value = res.data || {}
    }
  } catch (error) {
    console.error('加载公告详情失败', error)
  }
  loading.value = false
}

onMounted(() => {
  loadAnnouncementDetail()
})
</script>

<style scoped>
.announcement-detail-container {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.breadcrumb {
  margin-bottom: 20px;
}

.announcement-header {
  text-align: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.announcement-header h2 {
  font-size: 24px;
  color: #333;
  margin-bottom: 15px;
}

.announcement-header .meta {
  color: #909399;
  font-size: 14px;
}

.announcement-content {
  color: #333;
  line-height: 1.8;
  font-size: 16px;
}

.announcement-content p {
  margin-bottom: 15px;
  text-indent: 2em;
}

.back-btn {
  text-align: center;
  margin-top: 30px;
}
</style>
