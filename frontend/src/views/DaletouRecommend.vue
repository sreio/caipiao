<template>
  <div class="recommend-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-title">大乐透号码推荐</span>
        </div>
      </template>

      <!-- 控制栏 -->
      <div class="control-bar">
        <div class="control-item">
          <span class="label">生成注数：</span>
          <el-input-number
            v-model="ticketCount"
            :min="1"
            :max="10"
            :step="1"
            size="default"
            style="width: 120px"
          />
        </div>

        <div class="control-item">
          <el-radio-group v-model="mode">
            <el-radio label="random">随机推荐</el-radio>
            <el-radio label="manual">自选号码</el-radio>
          </el-radio-group>
        </div>

        <el-button type="primary" @click="generateNumbers" :loading="generating">
          <el-icon><Refresh /></el-icon>
          生成号码
        </el-button>

        <el-button v-if="tickets.length > 0" type="success" @click="handleExport">
          <el-icon><Download /></el-icon>
          生成预览图
        </el-button>

        <el-button v-if="tickets.length > 0" @click="clearAll">
          <el-icon><Delete /></el-icon>
          清空
        </el-button>
      </div>

      <!-- 自选号码区域 -->
      <div v-if="mode === 'manual'" class="manual-select-area">
        <el-card shadow="never" class="select-card">
          <template #header>
            <span>自选号码（第{{ currentTicketIndex + 1 }}注）</span>
          </template>

          <div class="ball-selection">
            <div class="ball-group">
              <div class="group-title">前区（选择5个）</div>
              <div class="ball-list">
                <span
                  v-for="num in 35"
                  :key="`front-${num}`"
                  class="ball-item front-ball-item"
                  :class="{ selected: selectedFront.includes(num) }"
                  @click="toggleFrontBall(num)"
                >
                  {{ String(num).padStart(2, '0') }}
                </span>
              </div>
            </div>

            <div class="ball-group">
              <div class="group-title">后区（选择2个）</div>
              <div class="ball-list">
                <span
                  v-for="num in 12"
                  :key="`back-${num}`"
                  class="ball-item back-ball-item"
                  :class="{ selected: selectedBack.includes(num) }"
                  @click="toggleBackBall(num)"
                >
                  {{ String(num).padStart(2, '0') }}
                </span>
              </div>
            </div>

            <div class="manual-actions">
              <el-button type="primary" @click="addManualTicket" :disabled="!canAddManual">
                添加这一注
              </el-button>
              <el-button @click="clearManual">清空当前</el-button>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 号码展示区域 -->
      <div v-if="tickets.length > 0" class="tickets-area">
        <div class="tickets-list">
          <el-card
            v-for="(ticket, index) in tickets"
            :key="index"
            class="ticket-card"
            shadow="hover"
          >
            <template #header>
              <div class="ticket-header">
                <span>第{{ index + 1 }}注</span>
                <el-button
                  type="danger"
                  size="small"
                  text
                  @click="removeTicket(index)"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
              </div>
            </template>

            <div class="ticket-content" :ref="el => setTicketRef(el, index)">
              <div class="ticket-balls">
                <div class="front-balls-group">
                  <span
                    v-for="num in ticket.front"
                    :key="`front-${num}`"
                    class="ball front-ball"
                  >
                    {{ String(num).padStart(2, '0') }}
                  </span>
                </div>
                <span class="separator">+</span>
                <div class="back-balls-group">
                  <span
                    v-for="num in ticket.back"
                    :key="`back-${num}`"
                    class="ball back-ball"
                  >
                    {{ String(num).padStart(2, '0') }}
                  </span>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>

      <!-- 空状态 -->
      <el-empty v-if="tickets.length === 0" description="暂无推荐号码，请点击生成号码" />
    </el-card>

    <!-- 预览图对话框 -->
    <el-dialog
      v-model="showPreview"
      title="预览图"
      width="600px"
      @close="closePreview"
    >
      <div class="preview-container" ref="previewRef">
        <div class="preview-header">
          <h2>大乐透号码推荐</h2>
          <div class="preview-date">{{ previewDate }}</div>
        </div>
        <div class="preview-tickets">
          <div
            v-for="(ticket, index) in tickets"
            :key="index"
            class="preview-ticket"
          >
            <div class="preview-ticket-label">第{{ index + 1 }}注</div>
            <div class="preview-ticket-balls">
              <span
                v-for="num in ticket.front"
                :key="`front-${num}`"
                class="preview-ball preview-front-ball"
              >
                {{ String(num).padStart(2, '0') }}
              </span>
              <span class="preview-separator">+</span>
              <span
                v-for="num in ticket.back"
                :key="`back-${num}`"
                class="preview-ball preview-back-ball"
              >
                {{ String(num).padStart(2, '0') }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="closePreview">取消</el-button>
        <el-button type="primary" @click="downloadPreview">保存图片</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import html2canvas from 'html2canvas'

const mode = ref('random') // random | manual
const ticketCount = ref(5)
const tickets = ref([])
const generating = ref(false)

// 自选号码
const currentTicketIndex = ref(0)
const selectedFront = ref([])
const selectedBack = ref([])

// 预览图
const showPreview = ref(false)
const previewRef = ref(null)
const ticketRefs = ref([])

const previewDate = computed(() => {
  const now = new Date()
  return `${now.getFullYear()}年${now.getMonth() + 1}月${now.getDate()}日`
})

const canAddManual = computed(() => {
  return selectedFront.value.length === 5 && selectedBack.value.length === 2
})

// 设置ticket引用
const setTicketRef = (el, index) => {
  if (el) {
    ticketRefs.value[index] = el
  }
}

// 随机生成号码
const generateRandomNumbers = () => {
  const newTickets = []
  for (let i = 0; i < ticketCount.value; i++) {
    // 生成5个不重复的前区号码（1-35）
    const front = []
    while (front.length < 5) {
      const num = Math.floor(Math.random() * 35) + 1
      if (!front.includes(num)) {
        front.push(num)
      }
    }
    front.sort((a, b) => a - b)

    // 生成2个不重复的后区号码（1-12）
    const back = []
    while (back.length < 2) {
      const num = Math.floor(Math.random() * 12) + 1
      if (!back.includes(num)) {
        back.push(num)
      }
    }
    back.sort((a, b) => a - b)

    newTickets.push({ front, back })
  }
  return newTickets
}

// 生成号码
const generateNumbers = () => {
  if (mode.value === 'random') {
    generating.value = true
    setTimeout(() => {
      tickets.value = generateRandomNumbers()
      generating.value = false
      ElMessage.success(`已生成${ticketCount.value}注号码`)
    }, 300)
  } else {
    ElMessage.warning('请使用自选号码功能添加号码')
  }
}

// 切换前区号码
const toggleFrontBall = (num) => {
  const index = selectedFront.value.indexOf(num)
  if (index > -1) {
    selectedFront.value.splice(index, 1)
  } else {
    if (selectedFront.value.length < 5) {
      selectedFront.value.push(num)
      selectedFront.value.sort((a, b) => a - b)
    } else {
      ElMessage.warning('最多只能选择5个前区号码')
    }
  }
}

// 切换后区号码
const toggleBackBall = (num) => {
  const index = selectedBack.value.indexOf(num)
  if (index > -1) {
    selectedBack.value.splice(index, 1)
  } else {
    if (selectedBack.value.length < 2) {
      selectedBack.value.push(num)
      selectedBack.value.sort((a, b) => a - b)
    } else {
      ElMessage.warning('最多只能选择2个后区号码')
    }
  }
}

// 添加自选号码
const addManualTicket = () => {
  if (!canAddManual.value) {
    ElMessage.warning('请选择5个前区号码和2个后区号码')
    return
  }

  tickets.value.push({
    front: [...selectedFront.value],
    back: [...selectedBack.value]
  })

  clearManual()
  ElMessage.success('已添加')
}

// 清空自选
const clearManual = () => {
  selectedFront.value = []
  selectedBack.value = []
}

// 删除一注
const removeTicket = (index) => {
  tickets.value.splice(index, 1)
  ElMessage.success('已删除')
}

// 清空所有
const clearAll = () => {
  tickets.value = []
  clearManual()
  ElMessage.success('已清空')
}

// 生成预览图
const handleExport = async () => {
  if (tickets.value.length === 0) {
    ElMessage.warning('请先生成号码')
    return
  }
  showPreview.value = true
}

// 关闭预览
const closePreview = () => {
  showPreview.value = false
}

// 下载预览图
const downloadPreview = async () => {
  if (!previewRef.value) return

  try {
    const canvas = await html2canvas(previewRef.value, {
      backgroundColor: '#fff',
      scale: 2,
      useCORS: true,
      logging: false
    })

    const url = canvas.toDataURL('image/png')
    const link = document.createElement('a')
    link.download = `大乐透推荐_${new Date().getTime()}.png`
    link.href = url
    link.click()

    ElMessage.success('图片已保存')
  } catch (error) {
    console.error('生成图片失败:', error)
    ElMessage.error('生成图片失败，请稍后重试')
  }
}
</script>

<style scoped>
.recommend-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.control-bar {
  display: flex;
  gap: 20px;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.control-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.label {
  font-size: 14px;
  color: #606266;
}

/* 自选号码区域 */
.manual-select-area {
  margin: 20px 0;
}

.select-card {
  border: 1px solid #e4e7ed;
}

.ball-selection {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.ball-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.group-title {
  font-size: 14px;
  font-weight: 600;
  color: #606266;
}

.ball-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.ball-item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid #dcdfe6;
  background: #fff;
}

.front-ball-item {
  color: #f56c6c;
}

.front-ball-item:hover {
  border-color: #f56c6c;
  background: #fef0f0;
}

.front-ball-item.selected {
  background: linear-gradient(135deg, #f56c6c, #e74c3c);
  color: #fff;
  border-color: #f56c6c;
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.3);
}

.back-ball-item {
  color: #409eff;
}

.back-ball-item:hover {
  border-color: #409eff;
  background: #ecf5ff;
}

.back-ball-item.selected {
  background: linear-gradient(135deg, #409eff, #3498db);
  color: #fff;
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.manual-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

/* 号码展示区域 */
.tickets-area {
  margin-top: 20px;
}

.tickets-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.ticket-card {
  border: 1px solid #e4e7ed;
}

.ticket-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: #303133;
}

.ticket-content {
  padding: 10px 0;
}

.ticket-balls {
  display: flex;
  align-items: center;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.front-balls-group {
  display: flex;
  gap: 8px;
}

.back-balls-group {
  display: flex;
  gap: 8px;
}

.ball {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  font-size: 16px;
  font-weight: 600;
  color: #fff;
}

.front-ball {
  background: linear-gradient(135deg, #f56c6c, #e74c3c);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.3);
}

.back-ball {
  background: linear-gradient(135deg, #409eff, #3498db);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.separator {
  font-size: 20px;
  font-weight: 600;
  color: #909399;
  margin: 0 4px;
}

/* 预览图 */
.preview-container {
  background: #fff;
  padding: 30px;
  border-radius: 8px;
}

.preview-header {
  text-align: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid #e4e7ed;
}

.preview-header h2 {
  margin: 0 0 10px 0;
  font-size: 24px;
  color: #303133;
}

.preview-date {
  font-size: 14px;
  color: #909399;
}

.preview-tickets {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.preview-ticket {
  padding: 20px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  background: #fafafa;
}

.preview-ticket-label {
  font-size: 14px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 12px;
}

.preview-ticket-balls {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
}

.preview-ball {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.preview-front-ball {
  background: linear-gradient(135deg, #f56c6c, #e74c3c);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.3);
}

.preview-back-ball {
  background: linear-gradient(135deg, #409eff, #3498db);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.preview-separator {
  font-size: 24px;
  font-weight: 600;
  color: #909399;
  margin: 0 8px;
}
</style>

