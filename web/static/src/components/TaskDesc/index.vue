<template>
  <el-dialog title="任务详情" :visible.sync="show" width="70%" center>
    <span>需要注意的是内容是默认不居中的</span>
      <span slot="footer" class="dialog-footer">
      <el-button @click="show = false">关 闭</el-button>
    </span>
  </el-dialog>
</template>

<script>
    export default {
      data() {
        return {
          show: false
        }
      },
      methods: {
        showTaskDesc() {
          this.show = true
        }
      }
    }
</script>

<style scoped>

</style>
