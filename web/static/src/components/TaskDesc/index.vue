<template>
  <el-dialog title="任务详情" :visible.sync="show" width="70%" center>
    <el-row border fit highlight-current-row>
      <el-form label-position="right" label-width="150px" label-suffix=":" class="task-desc-form">
        <el-col :span="12">
          <el-form-item :label="$t('task.id')">
            <span>{{row.id}}</span>
          </el-form-item>
          <el-form-item :label="$t('task.rule')">
            <span>{{row.task_rule_name}}</span>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('task.name')">
            <span>{{row.task_name}}</span>
          </el-form-item>
          <el-form-item :label="$t('task.desc')">
            <span>{{row.task_desc}}</span>
          </el-form-item>
        </el-col>
      </el-form>
    </el-row>

    <span slot="footer" class="dialog-footer">
      <el-button @click="show = false">关 闭</el-button>
    </span>
  </el-dialog>
</template>

<script>
    export default {
      data() {
        return {
          show: false,
          row: {}
        }
      },
      methods: {
        showTaskDesc(row) {
          this.show = true
          this.row = row
        }
      }
    }
</script>

<style scoped>
  .task-desc-form .el-form-item {
    margin-bottom: 0;
  }
</style>
