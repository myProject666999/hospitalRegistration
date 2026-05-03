<template>
  <div class="password-container">
    <el-card>
      <template #header>
        <span>修改密码</span>
      </template>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 400px">
        <el-form-item label="原密码" prop="old_password">
          <el-input v-model="form.old_password" type="password" placeholder="请输入原密码" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="form.new_password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm_password">
          <el-input v-model="form.confirm_password" type="password" placeholder="请确认新密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">确认修改</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { updatePassword } from '@/api/user'

const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  old_password: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为6-20个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await updatePassword({
      old_password: form.old_password,
      new_password: form.new_password
    })
    if (res.code === 200) {
      ElMessage.success('密码修改成功')
      form.old_password = ''
      form.new_password = ''
      form.confirm_password = ''
    } else {
      ElMessage.error(res.message || '修改失败')
    }
  } catch (error) {
    ElMessage.error('修改失败')
  }
  loading.value = false
}
</script>

<style scoped>
.password-container {
  padding: 20px;
}
</style>
