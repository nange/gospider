<template>
  <div class="app-container">
    <el-row v-loading="loadData" border fit highlight-current-row>
      <el-col :span="12">
        <el-form ref="form" :model="form" :rules="rules" label-position="right" label-width="150px">
          <el-form-item :label="$t('task.name')" prop="task_name">
            <el-input v-model="form.task_name" placeholder="please enter the content" :disabled="routeID"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.rule')" prop="task_rule_name">
            <el-select v-model="form.task_rule_name" placeholder="please select" :disabled="routeID">
              <el-option
                v-for="item in ruleOpts"
                :key="item"
                :label="item"
                :value="item">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('task.desc')">
            <el-input type="textarea" placeholder="please enter the content" :rows="2" v-model="form.task_desc">
            </el-input>
          </el-form-item>
          <el-form-item :label="$t('task.cron')">
            <el-input v-model="form.cron_spec" placeholder="compatible with crontab"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.proxy')">
            <el-input v-model="form.proxy_urls" placeholder="compatible with socks5,http,https; separated by commas"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.agent')">
            <el-input v-model="form.opt_user_agent" placeholder="user agent"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.maxDepth')">
            <el-input-number v-model="form.opt_max_depth" :controls="false"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('task.allowDomains')">
            <el-input placeholder="default empty, not limited" v-model="form.opt_allowed_domains"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.urlFilter')">
            <el-input placeholder="default empty, not limited, support regex" v-model="form.opt_url_filters"></el-input>
          </el-form-item>
          <el-form-item :label="$t('task.maxBody')">
            <el-input-number v-model="form.opt_max_body_size" :controls="false" class="fl"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('task.requestTimeout')">
            <el-input-number v-model="form.opt_request_timeout" :controls="false" class="fl"></el-input-number>
          </el-form-item>
          <el-form-item :label="$t('task.outType')" prop="output_type">
            <el-select v-model="form.output_type" placeholder="please select" @change="outputTypeChange">
              <el-option key="mysql" label="MYSQL" value="mysql"></el-option>
              <el-option key="csv" label="CSV" value="csv"></el-option>
              <el-option key="stdout" label="STDOUT" value="stdout"></el-option>
            </el-select>
            <el-select v-model="form.output_exportdb_id" placeholder="please select" v-if="showExportDB">
              <el-option
                v-for="item in exportDBList"
                :key="item.id"
                :label="item.show_name"
                :value="item.id">
              </el-option>
            </el-select>
            <el-checkbox v-model="form.auto_migrate" v-if="showExportDB">{{$t('task.autoMigrate')}}</el-checkbox>
          </el-form-item>

          <el-form-item :label="$t('task.limitEn')">
            <el-checkbox v-model="form.limit_enable"></el-checkbox>
          </el-form-item>
          <el-form-item :label="$t('task.limitDomainGlob')">
            <el-input v-model="form.limit_domain_glob" placeholder="default *"></el-input>
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
            <el-button type="primary" @click="on_submit_form" :loading="on_submit_loading" :disabled="submit_disable">{{routeID ? $t('task.update') : $t('task.add')}}</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { getTask, getRules, updateTask, saveTask } from '@/api/task'
import { fetchExportDBList } from '@/api/exportdb'
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
      showExportDB: false,
      ruleOpts: [],
      exportDBList: [],
      routeID: this.$route.params.id,
      loadData: false,
      on_submit_loading: false,
      submit_disable: false,
      rules: {
        task_name: [{ required: true, message: 'task name should not be empty', trigger: 'blur' }],
        task_rule_name: [{ required: true, message: 'please select rule name', trigger: 'change' }]
      }
    }
  },
  created() {
    this.getRules()
    this.routeID && this.getTaskRuleList()
    this.getExportDBList()
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
    getExportDBList() {
      this.loadData = true
      fetchExportDBList({
        offset: 0,
        size: -1
      }).then(response => {
        const data = response.data
        this.exportDBList = data.data
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
        this.showExportDB = true
      } else {
        this.showExportDB = false
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
            this.$message.success('task update success!  taskID:' + data.id + '  2s redirect to task list page!')
            this.on_submit_loading = false
            this.submit_disable = true
            setTimeout(() => this.$router.push({ name: 'taskList' }), 2000)
          }).catch(() => {
            this.on_submit_loading = false
          })
        } else {
          saveTask(this.form).then((response) => {
            const data = response.data
            this.$message.success('task create success!  taskID:' + data.id + '  2s redirect to task list page!')
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
