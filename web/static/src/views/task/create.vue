<template>
  <div class="app-container">
    <el-row v-loading="loadData" border fit highlight-current-row>
      <el-col :span="12">
        <el-form ref="form" :model="form" :rules="rules" label-position="right" label-width="150px">
          <el-form-item :label="$t('task.name')" prop="task_name">
            <el-input v-model="form.task_name" placeholder="请输入内容" :disabled="routeID"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.rule')" prop="task_rule_name">
            <el-select v-model="form.task_rule_name" placeholder="请选择" :disabled="routeID">
              <el-option
                v-for="item in ruleOpts"
                :key="item"
                :label="item"
                :value="item">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('task.desc')">
            <el-input type="textarea" placeholder="请输入内容" :rows="2" v-model="form.task_desc">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('task.cron')">
            <el-input v-model="form.cron_spec" placeholder="兼容crontab语法"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.proxy')">
            <el-input v-model="form.proxy_urls" placeholder="兼容socks5,http,https代理,列表以逗号分割"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.agent')">
            <el-input v-model="form.opt_user_agent" placeholder="User Agent"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.maxDepth')">
            <el-input-number v-model="form.opt_max_depth" :controls="false"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('task.allowDomains')">
            <el-input placeholder="默认空,不限制" v-model="form.opt_allowed_domains"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.urlFilter')">
            <el-input placeholder="默认空,不限制,可指定正则表达式" v-model="form.opt_url_filters"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.maxBody')">
            <el-input-number v-model="form.opt_max_body_size" :controls="false" class="fl"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('task.requestTimeout')">
            <el-input-number v-model="form.opt_request_timeout" :controls="false" class="fl"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('task.outType')" prop="output_type">
            <el-select v-model="form.output_type" placeholder="请选择" @change="outputTypeChange">
              <el-option key="mysql" label="MYSQL" value="mysql"></el-option>
              <el-option key="csv" label="CSV" value="csv"></el-option>
            </el-select>
            <el-select v-model="form.output_sysdb_id" placeholder="请选择" v-if="showSysDB">
              <el-option
                v-for="item in sysDBs"
                :key="item.id"
                :label="item.show_name"
                :value="item.id">
              </el-option>
            </el-select>
            <el-checkbox v-model="form.auto_migrate" v-if="showSysDB">{{$t('task.autoMigrate')}}</el-checkbox>
          </el-form-item>

          <el-form-item :label="$t('task.limitEn')">
            <el-checkbox v-model="form.limit_enable"></el-checkbox>
          </el-form-item>
          <el-form-item :label="$t('task.limitDomainGlob')">
            <el-input v-model="form.limit_domain_glob" placeholder="默认*,匹配所有域名"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.limitDelay')">
            <el-input-number v-model="form.limit_delay" :controls="false"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('task.limitRandomDelay')">
            <el-input-number v-model="form.limit_random_delay" :controls="false"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('task.limitPara')">
            <el-input-number v-model="form.limit_parallelism" :controls="false"></el-input-number>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="on_submit_form" :loading="on_submit_loading" :disabled="submit_disable">{{$t('task.add')}}</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { getTask, getRules, updateTask, saveTask } from '@/api/task'
import { fetchExportDbList } from '@/api/exportdb'
import waves from '@/directive/waves'
export default {
  name: 'taskCreate',
  directives: {
    waves
  },
  data() {
    return {
      form: {
        opt_user_agent: navigator.userAgent,
        limit_enable: true,
        auto_migrate: true,
        limit_parallelism: 1,
        opt_request_timeout: 10
      },
      showSysDB: false,
      ruleOpts: [],
      sysDBs: [],
      routeID: this.$route.params.id,
      loadData: false,
      on_submit_loading: false,
      submit_disable: false,
      rules: {
        task_name: [{ required: true, message: '任务名不能为空', trigger: 'blur' }],
        task_rule_name: [{ required: true, message: '请选择规则名称', trigger: 'change' }],
        output_type: [{ required: true, message: '请选择规导出类型', trigger: 'change' }]
      }
    }
  },
  created() {
    this.getRules()
    this.routeID && this.getTaskRuleList()
    this.getSysDBList()
  },
  methods: {
    // 获取数据
    getTaskRuleList() {
      this.loadData = true
      getTask(this.routeID).then(response => {
        const data = response.data
        this.form = data
        this.loadData = false
        this.outputTypeChange(data.output_type)
        setTimeout(() => {
          this.loadData = false
        }, 1.5 * 1000)
      })
    },
    // 获取导出数据库列表
    getSysDBList() {
      this.loadData = true
      fetchExportDbList({
        offset: 0,
        size: -1
      }).then(response => {
        const data = response.data
        this.sysDBs = data.data
        setTimeout(() => {
        }, 1.5 * 1000)
      })
      this.loadData = false
    },

    // 获取所有rules
    getRules() {
      getRules().then(response => {
        const data = response.data
        this.ruleOpts = data.data
        setTimeout(() => {
        }, 1.5 * 1000)
      })
    },

    // 改变output type
    outputTypeChange(output) {
      if (output === 'mysql') {
        this.showSysDB = true
      } else {
        this.showSysDB = false
      }
    },
    // 提交
    on_submit_form() {
      this.$refs.form.validate((valid) => {
        if (!valid) return false
        this.on_submit_loading = true
        if (this.routeID) {
          updateTask(this.routeID, this.form).then((response) => {
            const data = response.data
            this.$message.success('任务修改成功!  任务ID:' + data.id + '  3秒钟后跳转到任务列表页面!')
            this.on_submit_loading = false
            this.submit_disable = true
            setTimeout(() => this.$router.push({ name: 'taskList' }), 3000)
          }).catch(() => {
            this.on_submit_loading = false
          })
        } else {
          saveTask(this.form).then((response) => {
            const data = response.data
            this.$message.success('任务创建成功!  任务ID:' + data.id + '  3秒钟后跳转到任务列表页面!')
            this.on_submit_loading = false
            this.submit_disable = true
            setTimeout(() => this.$router.push({ name: 'taskList' }), 3000)
          }).catch(() => {
            this.on_submit_loading = false
          })
        }
      })
    }
  }
}
</script>
