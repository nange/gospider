<template>
  <div class="panel">
    <panel-title :title="$route.meta.title"></panel-title>
    <div class="panel-body"
         v-loading="loadData"
         element-loading-text="拼命加载中">
      <el-row>
        <el-col :span="12">
          <el-form ref="form" :model="form" :rules="rules" label-width="120px">
            <el-form-item label="显示名称:" prop="show_name">
              <el-input v-model="form.show_name" placeholder="请输入内容"></el-input>
            </el-form-item>
            <el-form-item label="主机地址:" prop="host">
              <el-input v-model="form.host" placeholder="请输入内容, 默认127.0.0.1"></el-input>
            </el-form-item>
            <el-form-item label="端口:" prop="port">
              <el-input-number v-model="form.port" :controls="false"></el-input-number>
            </el-form-item>
            <el-form-item label="用户名:" prop="user">
              <el-input v-model="form.user" placeholder="请输入内容, 默认root"></el-input>
            </el-form-item>
            <el-form-item label="密码:" prop="password">
              <el-input type="password" v-model="form.password" placeholder="请输入内容"></el-input>
            </el-form-item>
            <el-form-item label="数据库名:" prop="db_name">
              <el-input v-model="form.db_name" placeholder="请输入内容"></el-input>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="submitForm" :loading="submitLoading" :disabled="submitDisable">立即提交</el-button>
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script type="text/javascript">
  import {panelTitle} from 'components'

  export default{
    data(){
      return {
        form: {
          port: 3306
        },
        loadData: false,
        submitLoading: false,
        submitDisable: false,
        rules: {
          show_name: [{required: true, message: '显示名不能为空', trigger: 'blur'}],
          db_name: [{required: true, message: '数据库名不能为空', trigger: 'blur'}]
        }
      }
    },
    methods: {
      submitForm() {
        this.$refs.form.validate((valid) => {
          if (!valid) return false
          this.submitLoading = true
          this.$fetch.api_sysdb.create(this.form)
            .then((ret) => {
              this.$message.success('导出数据库记录创建成功!  ID:' + ret.id + '  3秒钟后跳转到数据库管理页面!')
              this.submitLoading = false
              this.submitDisable = true
              setTimeout(() => this.$router.push({name: 'sysDBMgt'}), 3000)
            })
            .catch(() => {
              this.submitLoading = false
            })
        })
      }
    },
    components: {
      panelTitle
    }
  }
</script>
