<template>
  <div class="panel">
    <panel-title :title="$route.meta.title"></panel-title>
    <div class="panel-body"
         v-loading="load_data"
         element-loading-text="拼命加载中">
      <el-row>
        <el-col :span="12">
          <el-form ref="form" :model="form" :rules="rules" label-width="120px">
            <el-form-item label="任务名称:" prop="task_name">
              <el-input v-model="form.task_name" placeholder="请输入内容"></el-input>
            </el-form-item>
            <el-form-item label="任务规则名:" prop="task_rule_name">
              <el-select v-model="form.task_rule_name" placeholder="请选择">
                <el-option
                  v-for="item in ruleOpts"
                  :key="item"
                  :label="item"
                  :value="item">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="任务描述:">
              <el-input type="textarea" placeholder="请输入内容" :rows="2" v-model="form.task_desc">
              </el-input>
            </el-form-item>
            <el-form-item label="定时执行:">
              <el-input v-model="form.cron_spec" placeholder="兼容crontab语法"></el-input>
            </el-form-item>
            <el-form-item label="User Agent:">
              <el-input v-model="form.opt_user_agent" placeholder="User Agent"></el-input>
            </el-form-item>
            <el-form-item label="爬虫最大深度:">
              <el-input-number v-model="form.opt_max_depth" :controls="false"></el-input-number>
            </el-form-item>
            <el-form-item label="允许访问的域名:">
              <el-input placeholder="默认空,不限制" v-model="form.opt_allowed_domains"></el-input>
            </el-form-item>
            <el-form-item label="URL过滤:">
              <el-input placeholder="默认空,不限制,可指定正则表达式" v-model="form.opt_url_filters"></el-input>
            </el-form-item>
            <el-form-item label="最大body值:">
              <el-input-number v-model="form.opt_max_body_size" :controls="false"></el-input-number>
            </el-form-item>
            <el-form-item label="导出类型:" prop="output_type">
              <el-select v-model="form.output_type" placeholder="请选择">
                <el-option key="mysql" label="MYSQL" value="mysql"></el-option>
                <el-option key="csv" label="CSV" value="csv"></el-option>
              </el-select>
              <el-row>
                <el-col :span="8"> <el-input v-model="form.output_mysql_host" placeholder="host"></el-input> </el-col>
                <el-col :span="8"> <el-input-number v-model="form.output_mysql_port" :controls="false"></el-input-number> </el-col>
                <el-col :span="8"> <el-input v-model="form.output_mysql_user" placeholder="user"></el-input> </el-col>
                <el-col :span="8"> <el-input v-model="form.output_mysql_password" placeholder="password"></el-input> </el-col>
                <el-col :span="8"> <el-input v-model="form.output_mysql_dbname" placeholder="dbname"></el-input> </el-col>
              </el-row>
            </el-form-item>

            <el-form-item label="频率限制:">
              <el-checkbox v-model="form.limit_enable">频率限制</el-checkbox>
            </el-form-item>
            <el-form-item label="域名glob匹配:">
              <el-input v-model="form.limit_domain_glob" placeholder="默认*,匹配所有域名"></el-input>
            </el-form-item>
            <el-form-item label="延迟:">
              <el-input-number v-model="form.limit_delay" :controls="false"></el-input-number>
            </el-form-item>
            <el-form-item label="随机延迟:">
              <el-input-number v-model="form.limit_random_delay" :controls="false"></el-input-number>
            </el-form-item>
            <el-form-item label="请求并发度:">
              <el-input-number v-model="form.limit_parallelism" :controls="false"></el-input-number>
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="on_submit_form" :loading="on_submit_loading" :disabled="submit_disable">立即提交</el-button>
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
          opt_user_agent: navigator.userAgent
        },
        ruleOpts: [],
        route_id: this.$route.params.id,
        load_data: false,
        on_submit_loading: false,
        submit_disable: false,
        rules: {
          task_name: [{required: true, message: '任务名不能为空', trigger: 'blur'}],
          task_rule_name: [{required: true, message: '请选择规则名称', trigger: 'change'}],
          output_type: [{required: true, message: '请选择规导出类型', trigger: 'change'}]
        }
      }
    },
    created(){
      this.getRules()
      this.route_id && this.get_form_data()
    },
    methods: {
      //获取数据
      get_form_data() {
        this.load_data = true
        this.$fetch.api_table.get({
          id: this.route_id
        })
          .then(({data}) => {
            this.form = data
            this.load_data = false
          })
          .catch(() => {
            this.load_data = false
          })
      },
      // 获取所有rules
      getRules() {
        this.$fetch.api_table.getRules()
          .then((data) => {
            this.ruleOpts = data.data
          })
      },
      //提交
      on_submit_form() {
        this.$refs.form.validate((valid) => {
          if (!valid) return false
          this.on_submit_loading = true
          this.$fetch.api_table.save(this.form)
            .then((ret) => {
              this.$message.success('任务创建成功!  任务ID:' + ret.id + '  3秒钟后跳转到任务列表页面!')
              this.on_submit_loading = false
              this.submit_disable = true
              setTimeout(() => this.$router.push({name: 'tableBase'}), 3000)
            })
            .catch(() => {
              this.on_submit_loading = false
            })
        })
      }
    },
    components: {
      panelTitle
    }
  }
</script>
