<template>
  <div class="profile-container">
    <el-card>
      <template #header>
        <span>个人信息</span>
      </template>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 500px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" disabled />
        </el-form-item>
        <el-form-item label="真实姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="form.gender">
            <el-radio :value="1">男</el-radio>
            <el-radio :value="2">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="身份证号" prop="id_card">
          <el-input v-model="form.id_card" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="form.address" type="textarea" placeholder="请输入地址" :rows="3" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">保存修改</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getUserInfo, updateUserInfo } from '@/api/user'

const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  username: '',
  name: '',
  phone: '',
  gender: 1,
  id_card: '',
  address: ''
})

const rules = {
  name: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const loadUserInfo = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 200) {
      Object.assign(form, res.data)
    }
  } catch (error) {
    console.error('获取用户信息失败', error)
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await updateUserInfo({
      name: form.name,
      phone: form.phone,
      gender: form.gender,
      id_card: form.id_card,
      address: form.address
    })
    if (res.code === 200) {
      ElMessage.success('修改成功')
    } else {
      ElMessage.error(res.message || '修改失败')
    }
  } catch (error) {
    ElMessage.error('修改失败')
  }
  loading.value = false
}

onMounted(() => {
  if (userStore.userInfo) {
    Object.assign(form, userStore.userInfo)
  }
  loadUserInfo()
})
</script>

<style scoped>
.profile-container {
  padding: 20px;
}
</style>
