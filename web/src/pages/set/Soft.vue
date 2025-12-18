<template>
  <div class="soft-config">
    <el-card shadow="hover">
      <div slot="header" class="card-header">
        <span class="header-title">
          <i class="el-icon-setting"></i>
          软件配置
        </span>
        <el-button 
          type="primary" 
          size="small"
          :disabled="!hasChanges"
          :loading="saving"
          @click="saveConfig">
          <i class="el-icon-check"></i>
          保存配置
        </el-button>
      </div>

      <el-alert
        v-if="hasChanges"
        title="提示：配置已修改，请点击保存按钮保存更改"
        type="warning"
        :closable="false"
        style="margin-bottom: 20px;">
      </el-alert>

      <el-table
        :data="displayData"
        border
        stripe
        :row-class-name="tableRowClassName"
        style="width: 100%">
        
        <el-table-column
          prop="info"
          label="信息"
          width="300"
          show-overflow-tooltip>
        </el-table-column>

        <el-table-column
          prop="name"
          label="配置"
          width="250"
          show-overflow-tooltip>
        </el-table-column>

        <el-table-column
          prop="env"
          label="环境变量"
          width="280"
          show-overflow-tooltip>
        </el-table-column>

        <el-table-column
          prop="val"
          label="数据"
          min-width="200">
          <template slot-scope="scope">
            <div class="config-value">
              <!-- 布尔值 -->
              <el-switch
                v-if="scope.row.type === 'bool'"
                v-model="scope.row.editValue"
                @change="handleValueChange(scope.row)"
                active-color="#13ce66">
              </el-switch>
              
              <!-- 数字值 -->
              <el-input-number
                v-else-if="scope.row.type === 'int'"
                v-model="scope.row.editValue"
                @change="handleValueChange(scope.row)"
                :min="0"
                :disabled="isReadOnly(scope.row.name)"
                size="small"
                controls-position="right"
                style="width: 100%;">
              </el-input-number>
              
              <!-- 字符串值 -->
              <el-input
                v-else
                v-model="scope.row.editValue"
                @input="handleValueChange(scope.row)"
                :disabled="isReadOnly(scope.row.name)"
                size="small"
                placeholder="请输入配置值">
              </el-input>

              <el-tag 
                v-if="scope.row.changed" 
                type="warning" 
                size="mini"
                style="margin-left: 10px;">
                已修改
              </el-tag>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="config-footer">
        <el-alert
          title="提示：修改配置后会自动热加载，无需重启服务。敏感配置（如密码、密钥）不可编辑。"
          type="info"
          :closable="false">
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Soft",
  created() {
    this.$emit('update:route_path', this.$route.path)
    this.$emit('update:route_name', ['基础信息', '软件配置'])
  },
  mounted() {
    this.getSoftInfo()
  },
  data() {
    return {
      soft_data: [],
      originalData: {},
      saving: false
    }
  },
  computed: {
    displayData() {
      return this.soft_data.map(item => {
        // 检查是否有修改
        const originalValue = this.originalData[item.name];
        const changed = originalValue !== undefined && originalValue !== item.editValue;
        
        return {
          ...item,
          changed,
          editValue: item.editValue !== undefined ? item.editValue : item.val
        }
      });
    },
    hasChanges() {
      return this.displayData.some(item => item.changed);
    }
  },
  methods: {
    getSoftInfo() {
      axios.get('/set/soft').then(resp => {
        var data = resp.data;
        console.log(data);
        this.soft_data = data.data.map(item => ({
          ...item,
          editValue: item.val,
          type: this.getValueType(item.val)
        }));
        
        // 保存原始数据用于对比
        this.originalData = {};
        this.soft_data.forEach(item => {
          this.originalData[item.name] = item.val;
        });
      }).catch(error => {
        if (error.response && error.response.status === 401) {
          return;
        }
        this.$message.error('获取配置失败');
        console.log(error);
      });
    },
    
    getValueType(value) {
      if (typeof value === 'boolean') {
        return 'bool';
      } else if (typeof value === 'number') {
        return 'int';
      } else {
        return 'string';
      }
    },
    
    isReadOnly(name) {
      // 根据后端返回的 sensitive 字段判断
      const item = this.soft_data.find(i => i.name === name);
      return item && item.sensitive === true;
    },
    
    handleValueChange() {
      this.$forceUpdate();
    },
    
    tableRowClassName({row}) {
      if (row.changed) {
        return 'warning-row';
      }
      return '';
    },
    
    saveConfig() {
      // 收集所有修改的配置
      const changedConfigs = {};
      this.displayData.forEach(item => {
        if (item.changed && !this.isReadOnly(item.name)) {
          changedConfigs[item.name] = item.editValue;
        }
      });
      
      if (Object.keys(changedConfigs).length === 0) {
        this.$message.info('没有需要保存的更改');
        return;
      }
      
      this.saving = true;
      
      axios.post('/set/soft/save', changedConfigs).then(resp => {
        const rdata = resp.data;
        if (rdata.code === 0) {
          this.$message.success('配置保存成功，已自动热加载');
          // 更新原始数据
          Object.keys(changedConfigs).forEach(key => {
            this.originalData[key] = changedConfigs[key];
          });
          // 刷新数据
          this.getSoftInfo();
        } else {
          this.$message.error(rdata.msg || '保存配置失败');
        }
      }).catch(error => {
        if (error.response && error.response.status === 401) {
          return;
        }
        this.$message.error('保存配置失败');
        console.log(error);
      }).finally(() => {
        this.saving = false;
      });
    }
  },
}
</script>

<style scoped>
.soft-config {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.header-title i {
  margin-right: 8px;
  color: #667eea;
}

.config-value {
  display: flex;
  align-items: center;
}

.config-footer {
  margin-top: 20px;
}

/* 修改行高亮 */
.soft-config >>> .warning-row {
  background-color: #fef0f0;
}

/* 表格样式优化 */
.soft-config >>> .el-table {
  font-size: 13px;
}

.soft-config >>> .el-table th {
  background-color: #f5f7fa;
  color: #606266;
  font-weight: 600;
}

.soft-config >>> .el-table td {
  padding: 12px 0;
}

/* 输入框样式 */
.soft-config >>> .el-input__inner {
  border-radius: 4px;
}

.soft-config >>> .el-input-number {
  width: 100%;
}

/* 禁用状态 */
.soft-config >>> .el-input.is-disabled .el-input__inner,
.soft-config >>> .el-input-number.is-disabled .el-input__inner {
  background-color: #f5f7fa;
  color: #c0c4cc;
  cursor: not-allowed;
}
</style>