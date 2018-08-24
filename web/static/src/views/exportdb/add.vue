<template>
  <div class="app-container">
    <el-row v-loading="loadData" border fit highlight-current-row>
      <el-col :span="12">
        <el-form ref="form" :model="form" :rules="rules" label-width="120px">
          <el-form-item :label="$t('exportdb.showname')" prop="show_name">
            <el-input v-model="form.show_name" placeholder="请输入内容"></el-input>
          </el-form-item>
          <el-form-item :label="$t('exportdb.host')" prop="host">
            <el-input v-model="form.host" placeholder="请输入内容, 默认127.0.0.1"></el-input>
          </el-form-item>
          <el-form-item :label="$t('exportdb.port')" prop="port">
            <el-input-number v-model="form.port" :controls="false"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('exportdb.user')" prop="user">
            <el-input v-model="form.user" placeholder="请输入内容, 默认root"></el-input>
          </el-form-item>
          <el-form-item :label="$t('exportdb.password')" prop="password">
            <el-input type="password" v-model="form.password" placeholder="请输入内容"></el-input>
          </el-form-item>
          <el-form-item :label="$t('exportdb.dbname')" prop="db_name">
            <el-input v-model="form.db_name" placeholder="请输入内容"></el-input>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="submitForm" :loading="submitLoading" :disabled="submitDisable">{{$t('exportdb.add')}}</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>
<script type="text/javascript">
  import { createExportDb } from '@/api/exportdb'
  import waves from '@/directive/waves' // 水波纹指令
  export default{
    directives: {
      waves
    },
    data() {
      return {
        form: {
          port: 3306
        },
        loadData: false,
        submitLoading: false,
        submitDisable: false,
        rules: {
          show_name: [{ required: true, message: '显示名不能为空', trigger: 'blur' }],
          db_name: [{ required: true, message: '数据库名不能为空', trigger: 'blur' }]
        }
      }
    },
    methods: {
      submitForm() {
        this.$refs.form.validate((valid) => {
          if (!valid) return false
          this.submitLoading = true
          createExportDb(this.form).then((response) => {
            const data = response.data
            this.$message.success('导出数据库记录创建成功!  ID:' + data.id + '  3秒钟后跳转到数据库管理页面!')
            this.submitLoading = false
            this.submitDisable = true
            setTimeout(() => this.$router.push({ name: 'expDbList' }), 3000)
          }).catch(() => {
            this.submitLoading = false
          })
        })
      }
    }
  }
</script>
