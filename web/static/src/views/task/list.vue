<template>
  <div class="app-container">
    <div class="filter-container fr">
      <el-button class="filter-item" @click="on_refresh" icon="el-icon-refresh" size="small"></el-button>
      <router-link :to="{name: 'createTask'}" tag="span">
        <el-button class="filter-item" @click="handleCreate" type="primary" icon="el-icon-edit" size="small">{{$t('task.add')}}</el-button>
      </router-link>
    </div>
    <el-table :data="tableData" v-loading="load_data" border fit highlight-current-row
  style="width: 100%;">
      <el-table-column prop="id" :label="$t('task.id')" width="80">
      </el-table-column>
      <el-table-column prop="task_name" :label="$t('task.name')" width="200">
      </el-table-column>
      <el-table-column prop="status" :label="$t('task.status')" width="120">
      </el-table-column>
      <el-table-column prop="counts" :label="$t('task.counts')" width="100">
      </el-table-column>
      <el-table-column prop="isCron" :label="$t('task.iscron')" width="100">
      </el-table-column>
      <el-table-column prop="created_at" :label="$t('task.create_at')" width="220">
      </el-table-column>
      <el-table-column :label="$t('task.actions')">
        <template scope="props">
          <el-button v-if="props.row.optionbutton & 0b10000" type="info" size="small" @click="showDesc(props.row)">{{$t('task.info')}}</el-button>
          <router-link :to="{name: 'editTask', params: {id: props.row.id}}" tag="span">
          <el-button v-if="props.row.optionbutton & 0b01000" type="warning" size="small" icon="edit">{{$t('task.edit')}}</el-button>
          </router-link>
          <el-button v-if="props.row.optionbutton & 0b00100" type="danger" size="small" @click="stop(props.row)">{{$t('task.stop')}}</el-button>
          <el-button v-if="props.row.optionbutton & 0b00010" type="success" size="small" @click="start(props.row)">{{$t('task.start')}}</el-button>
          <el-button v-if="props.row.optionbutton & 0b00001" type="success" size="small" @click="restart(props.row)">{{$t('task.restart')}}</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="pagination-container fr">
      <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="currentPage" :page-sizes="[10,20,30, 50]" :page-size="size" layout="total, sizes, prev, pager, next, jumper" :total="total">
      </el-pagination>
    </div>

    <el-dialog title="任务详情" :visible.sync="showDetails" width="60%" center>
      <el-row border fit highlight-current-row>
        <el-form label-position="right" label-width="150px" label-suffix=":" class="task-desc-form">
          <el-col :span="24">
            <el-form-item :label="$t('task.id')">
              <span>{{details.id}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.name')">
              <span>{{details.task_name}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.rule')">
              <span>{{details.task_rule_name}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.desc')">
              <span>{{details.task_desc}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.iscron')">
              <span>{{details.cron_spec ? '是' : '否'}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.cron')">
              <span>{{details.cron_spec}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.proxy')">
              <span>{{details.proxy_urls}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.agent')">
              <span>{{details.opt_user_agent}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.maxDepth')">
              <span>{{details.opt_max_depth}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.allowDomains')">
              <span>{{details.opt_allowed_domains}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.urlFilter')">
              <span>{{details.opt_url_filters}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.maxBody')">
              <span>{{details.opt_max_body_size}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.requestTimeout')">
              <span>{{details.opt_request_timeout}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.outType')">
              <span>{{details.output_type}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.autoMigrate')">
              <span>{{details.auto_migrate}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.limitEn')">
              <span>{{details.limit_enable}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.limitDomainGlob')">
              <span>{{details.limit_domain_glob}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.limitDelay')">
              <span>{{details.limit_delay}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.limitRandomDelay')">
              <span>{{details.limit_random_delay}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.limitPara')">
              <span>{{details.limit_parallelism}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.status')">
              <span>{{details.status}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.counts')">
              <span>{{details.counts}}</span>
            </el-form-item>
            <el-form-item :label="$t('task.create_at')">
              <span>{{details.created_at}}</span>
            </el-form-item>
          </el-col>

        </el-form>
      </el-row>

      <span slot="footer" class="dialog-footer">
        <el-button @click="showDetails = false">关 闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
  import { fetchTaskList, stopTask, startTask, restartTask } from '@/api/task'
  import waves from '@/directive/waves' // 水波纹指令

  export default {
    directives: {
      waves
    },
    data() {
      return {
        tableData: [],
        // 当前页码
        currentPage: 1,
        // 数据总条目
        total: 0,
        // 每页显示多少条数据
        size: 10,
        // 请求时的loading效果
        load_data: true,
        // 批量选择数组
        batch_select: [],
        // 详情
        details: {},
        showDetails: false
      }
    },
    created() {
      this.get_table_data()
    },
    methods: {
      // 刷新
      on_refresh() {
        this.get_table_data()
      },
      // 获取数据
      get_table_data() {
        this.load_data = true
        fetchTaskList({
          offset: (this.currentPage - 1) * this.size,
          size: this.size
        }).then((response) => {
          const data = response.data
          this.tableData = data.data
          for (const v of this.tableData) {
            v.isCron = '否'
            if (v.cron_spec) {
              v.isCron = '是'
            }
            // 操作按钮，用5位2进制数表示，每位控制一个按钮是否显示
            // ----0----0----0----0----0----
            // ----|----|----|----|----|----
            // ---详情--修改-停止--启动-重启---
            switch (v.status) {
              case '未知状态':
                v.optionbutton = 0b10000
                break
              case '运行中':
                v.optionbutton = 0b11100
                break
              case '暂停':
                v.optionbutton = 0b10100
                break
              case '停止':
                v.optionbutton = 0b11010
                if (v.cron_spec) { v.optionbutton = 0b11001 }
                break
              case '异常退出':
                v.optionbutton = 0b11010
                if (v.cron_spec) { v.optionbutton = 0b11100 }
                break
              case '完成':
                v.optionbutton = 0b11010
                if (v.cron_spec) { v.optionbutton = 0b11100 }
                break
              case '运行超时':
                v.optionbutton = 0b11010
                if (v.cron_spec) { v.optionbutton = 0b1110 }
                break
              default:
                v.optionbutton = 0b10000
                break
            }
          }

          this.total = data.total
          this.load_data = false
        }).catch(() => {
          this.load_data = false
        })
      },
      showDesc(row) {
        this.showDetails = true
        this.details = row
      },
      stop(item) {
        this.$confirm('此操作将停止该任务, 是否继续?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.load_data = true
          stopTask(item.id).then(() => {
            this.get_table_data()
            this.$message.success('操作成功!')
          }).catch(() => {
            this.$message.error('停止任务出错!')
          })
        })
      },
      // 非定时任务启动
      start(item) {
        this.$confirm('是否启动该任务?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.load_data = true
          startTask(item.id).then(() => {
            this.get_table_data()
            this.$message.success('操作成功!')
          }).catch(() => {
            this.$message.error('启动任务出错!')
          })
        })
      },
      // 定时任务重启
      restart(item) {
        this.$confirm('是否重启该定时任务?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.load_data = true
          restartTask(item.id).then(() => {
            this.get_table_data()
            this.$message.success('操作成功!')
          }).catch(() => {
            this.$message.error('重启任务出错!')
          })
        })
      },
      handleSizeChange(val) {
        this.size = val
        this.get_table_data()
      },
      // 页码选择
      handleCurrentChange(val) {
        this.currentPage = val
        this.get_table_data()
      }
    }
  }
</script>

<style scoped>
  .task-desc-form .el-form-item {
    margin-bottom: 0;
  }
</style>
