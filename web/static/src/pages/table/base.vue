<template>
  <div class="panel">
    <panel-title :title="$route.meta.title">
      <el-button @click.stop="on_refresh" size="small">
        <i class="fa fa-refresh"></i>
      </el-button>
      <router-link :to="{name: 'tableAdd'}" tag="span">
        <el-button type="primary" icon="plus" size="small">添加任务</el-button>
      </router-link>
    </panel-title>
    <div class="panel-body">
      <el-table
        :data="tableData"
        stripe
        style="width: 100%">
        <el-table-column
          prop="id"
          label="ID"
          width="80">
        </el-table-column>
        <el-table-column
          prop="task_name"
          label="任务名"
          width="200">
        </el-table-column>
        <el-table-column
          prop="status"
          label="状态"
          width="120">
        </el-table-column>
        <el-table-column
          prop="counts"
          label="运行次数"
          width="100">
        </el-table-column>
        <el-table-column
          prop="isCron"
          label="定时任务"
          width="100">
        </el-table-column>
        <el-table-column
          prop="created_at"
          label="创建时间"
          width="220">
        </el-table-column>
        <el-table-column
          label="操作">
          <template scope="props">
            <el-button v-if="props.row.optionbutton & 0b10000" type="info" size="small">详情</el-button>
            <!-- <router-link :to="{name: 'tableUpdate', params: {id: props.row.id}}" tag="span"> -->
            <el-button v-if="props.row.optionbutton & 0b01000" type="info" size="small" icon="edit">修改</el-button>
            <!-- </router-link> -->
            <el-button v-if="props.row.optionbutton & 0b00100" type="danger" size="small" @click="stop(props.row)">停止</el-button>
            <el-button v-if="props.row.optionbutton & 0b00010" type="info" size="small" @click="start(props.row)">启动</el-button>
            <el-button v-if="props.row.optionbutton & 0b00001" type="info" size="small" @click="restart(props.row)">重启</el-button>
          </template>
        </el-table-column>
      </el-table>
      <bottom-tool-bar>
        <div slot="page">
          <el-pagination
            @current-change="handleCurrentChange"
            :current-page="currentPage"
            :page-size="10"
            layout="total, prev, pager, next"
            :total="total">
          </el-pagination>
        </div>
      </bottom-tool-bar>
    </div>
  </div>
</template>
<script type="text/javascript">
  import {panelTitle, bottomToolBar} from 'components'
  export default{
    data(){
      return {
        tableData: [],
        //当前页码
        currentPage: 1,
        //数据总条目
        total: 0,
        //每页显示多少条数据
        size: 10,
        //请求时的loading效果
        load_data: true,
        //批量选择数组
        batch_select: []
      }
    },
    components: {
      panelTitle,
      bottomToolBar
    },
    created(){
      this.get_table_data()
    },
    methods: {
      //刷新
      on_refresh(){
        this.get_table_data()
      },
      //获取数据
      get_table_data(){
        this.load_data = true;
        this.$fetch.api_table.list({
          offset: (this.currentPage-1)*this.size,
          size: this.size
        }).then((ret) => {
            this.tableData = ret.data;
            for (let v of this.tableData) {
              v.isCron = '否';
              if (v.cron_spec) {
                v.isCron = '是';
              }
              //操作按钮，用5位2进制数表示，每位控制一个按钮是否显示
              // ----0----0----0----0----0----
              // ----|----|----|----|----|----
              // ---详情--修改-停止--启动-重启---
              switch(v.status)
              {
                case "未知状态":
                  v.optionbutton = 0b10000;
                  break;
                case "运行中":
                  v.optionbutton = 0b10100;
                  break;
                case "暂停":
                  v.optionbutton = 0b10100;
                  break;
                case "停止":
                  v.optionbutton = 0b11010;
                  if (v.cron_spec){v.optionbutton = 0b11001;}
                  break;
                case "异常退出":
                  v.optionbutton = 0b11010;
                  if (v.cron_spec){v.optionbutton = 0b11100;}
                  break;
                case "完成":
                  v.optionbutton = 0b11010;
                  if (v.cron_spec){v.optionbutton = 0b11100;}
                  break;
                case "运行超时":
                  v.optionbutton = 0b11010;
                  if (v.cron_spec){v.optionbutton = 0b11100;}
                  break;
                default:
                  v.optionbutton = 0b10000;
                  break;
              }
            }

            this.total = ret.total;
            this.load_data = false
        }).catch(() => {
          this.load_data = false
        })
      },
      //单个删除
      stop(item){
        this.$confirm('此操作将停止该任务, 是否继续?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        .then(() => {
          this.load_data = true;
          this.$fetch.api_table.stop(item.id)
            .then(() => {
              this.get_table_data();
              this.$message.success('操作成功!')
            })
            .catch((error) => {
              this.$message.error('停止任务出错!')
            })
        })
      },
      //非定时任务启动
      start(item){
        this.$confirm('是否启动该任务?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        .then(() => {
          this.load_data = true;
          this.$fetch.api_table.start(item.id)
            .then(() => {
              this.get_table_data();
              this.$message.success('操作成功!')
            })
            .catch((error) => {
              this.$message.error('启动任务出错!')
            })
        })
      },
      //定时任务重启
      restart(item){
        this.$confirm('是否重启该定时任务?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        .then(() => {
          this.load_data = true;
          this.$fetch.api_table.restart(item.id)
            .then(() => {
              this.get_table_data();
              this.$message.success('操作成功!')
            })
            .catch((error) => {
              this.$message.error('重启任务出错!')
            })
        })
      },
      //页码选择
      handleCurrentChange(val) {
        this.currentPage = val;
        this.get_table_data()
      }
    }
  }
</script>
