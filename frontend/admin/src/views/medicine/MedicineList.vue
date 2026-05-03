<template>
  <div class="page-container">
    <div class="page-header">
      <h3>药品管理</h3>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新增药品
      </el-button>
    </div>

    <el-card>
      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="药品名称" />
        <el-table-column prop="code" label="药品编码" width="120" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="spec" label="规格" width="120" />
        <el-table-column prop="price" label="价格">
          <template #default="scope">
            <span class="price">¥{{ scope.row.price || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
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
          @size-change="loadList"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑药品' : '新增药品'" width="700px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="药品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入药品名称" />
        </el-form-item>
        <el-form-item label="药品编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入药品编码" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分类" prop="category">
              <el-input v-model="form.category" placeholder="请输入分类" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="规格" prop="spec">
              <el-input v-model="form.spec" placeholder="请输入规格" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="单位" prop="unit">
              <el-input v-model="form.unit" placeholder="请输入单位" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="价格" prop="price">
              <el-input-number v-model="form.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="库存" prop="stock">
              <el-input-number v-model="form.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="正常" :value="1" />
                <el-option label="禁用" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="生产厂家" prop="manufacturer">
          <el-input v-model="form.manufacturer" placeholder="请输入生产厂家" />
        </el-form-item>
        <el-form-item label="用法用量" prop="usage">
          <el-input v-model="form.usage" type="textarea" placeholder="请输入用法用量" :rows="2" />
        </el-form-item>
        <el-form-item label="适应症" prop="indication">
          <el-input v-model="form.indication" type="textarea" placeholder="请输入适应症" :rows="2" />
        </el-form-item>
        <el-form-item label="药品描述" prop="description">
          <el-input v-model="form.description" type="textarea" placeholder="请输入药品描述" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import { Plus } from '@element-plus/icons-vue'

const list = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const dialogVisible = ref(false)
const submitLoading = ref(false)
const formRef = ref(null)
const isEdit = ref(false)
const currentId = ref(null)

const form = reactive({
  name: '',
  code: '',
  category: '',
  spec: '',
  unit: '',
  price: 0,
  manufacturer: '',
  usage: '',
  indication: '',
  description: '',
  stock: 0,
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入药品名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入药品编码', trigger: 'blur' }]
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/admin/medicines',
      method: 'get',
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    if (res.code === 200) {
      list.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } catch (error) {
    console.error('加载列表失败', error)
  }
  loading.value = false
}

const resetForm = () => {
  Object.assign(form, {
    name: '',
    code: '',
    category: '',
    spec: '',
    unit: '',
    price: 0,
    manufacturer: '',
    usage: '',
    indication: '',
    description: '',
    stock: 0,
    status: 1
  })
}

const handleAdd = () => {
  isEdit.value = false
  currentId.value = null
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  currentId.value = row.id
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该药品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await request({
      url: `/admin/medicines/${row.id}`,
      method: 'delete'
    })
    if (res.code === 200) {
      ElMessage.success('删除成功')
      loadList()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const res = await request({
      url: isEdit.value ? `/admin/medicines/${currentId.value}` : '/admin/medicines',
      method: isEdit.value ? 'put' : 'post',
      data: form
    })
    if (res.code === 200) {
      ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
      dialogVisible.value = false
      loadList()
    } else {
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
  submitLoading.value = false
}

onMounted(() => {
  loadList()
})
</script>

<style scoped>
.page-container {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h3 {
  margin: 0;
  font-size: 18px;
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
