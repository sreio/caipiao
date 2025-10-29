<template>
  <div class="trend-container">
    <!-- 控制栏 -->
    <el-card class="control-card">
      <div class="control-bar">
        <el-select v-model="limit" @change="loadTrend" style="width: 150px">
          <el-option label="最近30期" :value="30" />
          <el-option label="最近50期" :value="50" />
          <el-option label="最近100期" :value="100" />
        </el-select>
        
        <el-button type="primary" @click="loadTrend" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新数据
        </el-button>
      </div>
    </el-card>

    <!-- 走势图 -->
    <el-card class="trend-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">双色球号码走势图</span>
          <span class="period-info">共 {{ tableData.length }} 期</span>
        </div>
      </template>
      
      <div class="trend-table-wrapper">
        <table class="trend-table">
          <thead>
            <tr>
              <th class="period-col">期号</th>
              <th class="date-col">开奖日期</th>
              <th colspan="33" class="number-group-header red-header">红球区</th>
              <th colspan="16" class="number-group-header blue-header">蓝球区</th>
            </tr>
            <tr>
              <th class="period-col"></th>
              <th class="date-col"></th>
              <!-- 红球号码 1-33 -->
              <th v-for="num in 33" :key="`red-${num}`" class="number-header red-bg">
                {{ String(num).padStart(2, '0') }}
              </th>
              <!-- 蓝球号码 1-16 -->
              <th v-for="num in 16" :key="`blue-${num}`" class="number-header blue-bg">
                {{ String(num).padStart(2, '0') }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, rowIndex) in tableData" :key="row.issue" class="trend-row">
              <td class="period-cell">{{ row.issue }}</td>
              <td class="date-cell">{{ formatDate(row.draw_date) }}</td>
              
              <!-- 红球走势 -->
              <td v-for="num in 33" :key="`red-${row.issue}-${num}`" class="trend-cell">
                <div class="cell-content">
                  <span v-if="row.red_balls.includes(num)" class="ball-mark red-mark">
                    {{ String(num).padStart(2, '0') }}
                  </span>
                  <span v-else class="empty-mark">{{ num }}</span>
                  <!-- 连接线 -->
                  <svg v-if="rowIndex < tableData.length - 1 && hasConnection(rowIndex, num, 'red')" 
                       class="connect-line" viewBox="0 0 2 40">
                    <line x1="1" y1="0" x2="1" y2="40" stroke="#f56c6c" stroke-width="2"/>
                  </svg>
                </div>
              </td>
              
              <!-- 蓝球走势 -->
              <td v-for="num in 16" :key="`blue-${row.issue}-${num}`" class="trend-cell">
                <div class="cell-content">
                  <span v-if="row.blue_ball === num" class="ball-mark blue-mark">
                    {{ String(num).padStart(2, '0') }}
                  </span>
                  <span v-else class="empty-mark">{{ num }}</span>
                  <!-- 连接线 -->
                  <svg v-if="rowIndex < tableData.length - 1 && hasConnection(rowIndex, num, 'blue')" 
                       class="connect-line" viewBox="0 0 2 40">
                    <line x1="1" y1="0" x2="1" y2="40" stroke="#409eff" stroke-width="2"/>
                  </svg>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 红球频率统计 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">红球频率统计（Top 10）</span>
            </div>
          </template>
          <div ref="redFreqChartRef" style="width: 100%; height: 400px"></div>
          <el-table :data="topRedFreq" style="margin-top: 20px">
            <el-table-column prop="number" label="号码" min-width="100">
              <template #default="{ row }">
                <span class="ball red-ball">{{ row.number }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="出现次数" min-width="120" />
            <el-table-column label="出现率" min-width="120">
              <template #default="{ row }">
                {{ ((row.count / (limit * 6)) * 100).toFixed(1) }}%
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 蓝球频率统计 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">蓝球频率统计（Top 10）</span>
            </div>
          </template>
          <div ref="blueFreqChartRef" style="width: 100%; height: 400px"></div>
          <el-table :data="topBlueFreq" style="margin-top: 20px">
            <el-table-column prop="number" label="号码" min-width="100">
              <template #default="{ row }">
                <span class="ball blue-ball">{{ row.number }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="count" label="出现次数" min-width="120" />
            <el-table-column label="出现率" min-width="120">
              <template #default="{ row }">
                {{ ((row.count / limit) * 100).toFixed(1) }}%
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 遗漏值分析 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <!-- 红球遗漏 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">红球遗漏值（Top 10）</span>
              <el-tag type="warning" size="small">提醒：遗漏值较大的号码可能即将出现</el-tag>
            </div>
          </template>
          <el-alert
            v-if="topRedMissing.length > 0 && topRedMissing[0].missing > 20"
            title="长期遗漏提醒"
            :description="`号码 ${topRedMissing[0].number} 已连续 ${topRedMissing[0].missing} 期未出现！`"
            type="warning"
            show-icon
            style="margin-bottom: 20px"
          />
          <el-table :data="topRedMissing">
            <el-table-column prop="number" label="号码" min-width="100">
              <template #default="{ row }">
                <span class="ball red-ball">{{ row.number }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="missing" label="遗漏期数" min-width="140">
              <template #default="{ row }">
                <el-tag :type="getMissingTagType(row.missing)">
                  {{ row.missing }} 期
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <span v-if="row.missing === 0" style="color: #67c23a">●刚出现</span>
                <span v-else-if="row.missing < 10" style="color: #409eff">●正常</span>
                <span v-else-if="row.missing < 20" style="color: #e6a23c">●较长</span>
                <span v-else style="color: #f56c6c">●长期未出</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 蓝球遗漏 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span class="card-title">蓝球遗漏值（Top 10）</span>
              <el-tag type="warning" size="small">提醒：遗漏值较大的号码可能即将出现</el-tag>
            </div>
          </template>
          <el-alert
            v-if="topBlueMissing.length > 0 && topBlueMissing[0].missing > 10"
            title="长期遗漏提醒"
            :description="`号码 ${topBlueMissing[0].number} 已连续 ${topBlueMissing[0].missing} 期未出现！`"
            type="warning"
            show-icon
            style="margin-bottom: 20px"
          />
          <el-table :data="topBlueMissing">
            <el-table-column prop="number" label="号码" min-width="100">
              <template #default="{ row }">
                <span class="ball blue-ball">{{ row.number }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="missing" label="遗漏期数" min-width="140">
              <template #default="{ row }">
                <el-tag :type="getMissingTagType(row.missing)">
                  {{ row.missing }} 期
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <span v-if="row.missing === 0" style="color: #67c23a">●刚出现</span>
                <span v-else-if="row.missing < 5" style="color: #409eff">●正常</span>
                <span v-else-if="row.missing < 10" style="color: #e6a23c">●较长</span>
                <span v-else style="color: #f56c6c">●长期未出</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { shuangseqiuAPI } from '@/api/lottery'
import * as echarts from 'echarts'

const loading = ref(false)
const limit = ref(50)
const tableData = ref([])
const trendData = ref(null)

// 图表实例
const redFreqChartRef = ref()
const blueFreqChartRef = ref()
let redFreqChart = null
let blueFreqChart = null

// 计算属性
const topRedFreq = computed(() => {
  if (!trendData.value?.red_freq) return []
  return Object.entries(trendData.value.red_freq)
    .map(([number, count]) => ({ number: parseInt(number), count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 10)
})

const topBlueFreq = computed(() => {
  if (!trendData.value?.blue_freq) return []
  return Object.entries(trendData.value.blue_freq)
    .map(([number, count]) => ({ number: parseInt(number), count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 10)
})

const topRedMissing = computed(() => {
  if (!trendData.value?.red_missing) return []
  return Object.entries(trendData.value.red_missing)
    .map(([number, missing]) => ({ number: parseInt(number), missing }))
    .sort((a, b) => b.missing - a.missing)
    .slice(0, 10)
})

const topBlueMissing = computed(() => {
  if (!trendData.value?.blue_missing) return []
  return Object.entries(trendData.value.blue_missing)
    .map(([number, missing]) => ({ number: parseInt(number), missing }))
    .sort((a, b) => b.missing - a.missing)
    .slice(0, 10)
})

// 获取遗漏标签类型
const getMissingTagType = (missing) => {
  if (missing === 0) return 'success'
  if (missing < 10) return 'info'
  if (missing < 20) return 'warning'
  return 'danger'
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}-${date.getDate()}`
}

// 判断是否有连接线
const hasConnection = (rowIndex, num, type) => {
  if (rowIndex >= tableData.value.length - 1) return false
  
  const currentRow = tableData.value[rowIndex]
  const nextRow = tableData.value[rowIndex + 1]
  
  if (type === 'red') {
    return currentRow.red_balls.includes(num) && nextRow.red_balls.includes(num)
  } else {
    return currentRow.blue_ball === num && nextRow.blue_ball === num
  }
}

// 加载走势数据
const loadTrend = async () => {
  loading.value = true
  try {
    const data = await shuangseqiuAPI.getTrend(limit.value)
    trendData.value = data
    
    // 处理表格数据
    tableData.value = data.issues.map((issue, index) => ({
      issue: issue,
      draw_date: data.draw_dates ? data.draw_dates[index] : '',
      red_balls: data.red_balls[index] || [],
      blue_ball: data.blue_balls[index] || 0
    }))
    
    // 渲染图表
    renderFreqCharts()
  } catch (error) {
    console.error('加载走势数据失败:', error)
    ElMessage.error('加载走势数据失败')
  } finally {
    loading.value = false
  }
}


// 渲染频率图表
const renderFreqCharts = () => {
  // 红球频率图
  if (!redFreqChart) {
    redFreqChart = echarts.init(redFreqChartRef.value)
  }

  const redFreqData = topRedFreq.value
  const redFreqOption = {
    title: {
      text: '红球出现频率',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: redFreqData.map(item => item.number)
    },
    yAxis: {
      type: 'value',
      name: '出现次数'
    },
    series: [
      {
        data: redFreqData.map(item => ({
          value: item.count,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#f56c6c' },
              { offset: 1, color: '#e74c3c' }
            ])
          }
        })),
        type: 'bar',
        barWidth: '60%'
      }
    ]
  }

  redFreqChart.setOption(redFreqOption)

  // 蓝球频率图
  if (!blueFreqChart) {
    blueFreqChart = echarts.init(blueFreqChartRef.value)
  }

  const blueFreqData = topBlueFreq.value
  const blueFreqOption = {
    title: {
      text: '蓝球出现频率',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: blueFreqData.map(item => item.number)
    },
    yAxis: {
      type: 'value',
      name: '出现次数'
    },
    series: [
      {
        data: blueFreqData.map(item => ({
          value: item.count,
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#409eff' },
              { offset: 1, color: '#3498db' }
            ])
          }
        })),
        type: 'bar',
        barWidth: '60%'
      }
    ]
  }

  blueFreqChart.setOption(blueFreqOption)
}

// 窗口大小变化时重新渲染
const handleResize = () => {
  redFreqChart?.resize()
  blueFreqChart?.resize()
}

onMounted(() => {
  loadTrend()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  redFreqChart?.dispose()
  blueFreqChart?.dispose()
})
</script>

<style scoped>
.trend-container {
  padding: 20px;
}

.control-card {
  margin-bottom: 20px;
}

.control-bar {
  display: flex;
  gap: 12px;
  align-items: center;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.period-info {
  font-size: 14px;
  color: #909399;
}

/* 走势表格 */
.trend-card {
  margin-bottom: 20px;
}

.trend-table-wrapper {
  overflow-x: auto;
  max-height: 800px;
  overflow-y: auto;
}

.trend-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
  min-width: 1200px;
}

.trend-table thead {
  position: sticky;
  top: 0;
  z-index: 10;
  background: #fff;
}

.number-group-header {
  background: #f5f7fa;
  font-weight: 600;
  padding: 8px 4px;
  text-align: center;
  border: 1px solid #e4e7ed;
  font-size: 14px;
}

.red-header {
  background: linear-gradient(to bottom, #fef0f0, #fde2e2);
  color: #f56c6c;
}

.blue-header {
  background: linear-gradient(to bottom, #ecf5ff, #d9ecff);
  color: #409eff;
}

.period-col,
.date-col,
.number-header {
  padding: 6px 4px;
  text-align: center;
  border: 1px solid #e4e7ed;
  background: #fff;
  font-weight: 600;
}

.period-col {
  min-width: 80px;
  position: sticky;
  left: 0;
  z-index: 11;
  background: #fff;
}

.date-col {
  min-width: 60px;
  position: sticky;
  left: 80px;
  z-index: 11;
  background: #fff;
}

.number-header {
  width: 28px;
  min-width: 28px;
  font-size: 11px;
}

.red-bg {
  background: #fef0f0;
  color: #f56c6c;
}

.blue-bg {
  background: #ecf5ff;
  color: #409eff;
}

.trend-row {
  border-bottom: 1px solid #e4e7ed;
}

.trend-row:hover {
  background: #f5f7fa;
}

.period-cell {
  padding: 4px;
  text-align: center;
  font-weight: 600;
  font-size: 12px;
  border: 1px solid #e4e7ed;
  position: sticky;
  left: 0;
  background: #fff;
  z-index: 1;
}

.trend-row:hover .period-cell {
  background: #f5f7fa;
}

.date-cell {
  padding: 4px;
  text-align: center;
  font-size: 11px;
  color: #909399;
  border: 1px solid #e4e7ed;
  position: sticky;
  left: 80px;
  background: #fff;
  z-index: 1;
}

.trend-row:hover .date-cell {
  background: #f5f7fa;
}

.trend-cell {
  padding: 0;
  text-align: center;
  border: 1px solid #e4e7ed;
  height: 40px;
  position: relative;
  width: 28px;
  min-width: 28px;
}

.cell-content {
  position: relative;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ball-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  position: relative;
  z-index: 2;
}

.red-mark {
  background: linear-gradient(135deg, #f56c6c, #e74c3c);
  box-shadow: 0 1px 4px rgba(245, 108, 108, 0.4);
}

.blue-mark {
  background: linear-gradient(135deg, #409eff, #3498db);
  box-shadow: 0 1px 4px rgba(64, 158, 255, 0.4);
}

.empty-mark {
  font-size: 10px;
  color: #dcdfe6;
  font-weight: normal;
}

.connect-line {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translateX(-50%);
  width: 2px;
  height: 40px;
  pointer-events: none;
  z-index: 1;
}

/* 频率统计表格样式 */
.ball {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  font-size: 16px;
  font-weight: 600;
  color: #fff;
}

.red-ball {
  background: linear-gradient(135deg, #f56c6c, #e74c3c);
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.3);
}

.blue-ball {
  background: linear-gradient(135deg, #409eff, #3498db);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}
</style>

