<template>
  <div class="panel">
    <panel-title :title="$route.meta.title">
      <el-button @click.stop="refresh" size="small">
        <i class="fa fa-refresh"></i>
      </el-button>
      <router-link :to="{name: 'sysDBAdd'}" tag="span">
        <el-button type="primary" icon="plus" size="small">添加数据库</el-button>
      </router-link>
    </panel-title>
    <div class="panel-body">
      <el-table
        :data="dbData"
        stripe
        style="width: 100%">
        <el-table-column
          prop="id"
          label="ID"
          width="100">
        </el-table-column>
        <el-table-column
          prop="show_name"
          label="显示名称"
          width="200">
        </el-table-column>
        <el-table-column
          prop="host"
          label="主机地址"
          width="200">
        </el-table-column>
        <el-table-column
          prop="port"
          label="端口"
          width="100">
        </el-table-column>
        <el-table-column
          prop="user"
          label="用户名"
          width="100">
        </el-table-column>
        <el-table-column
          prop="db_name"
          label="数据库名"
          width="150">
        </el-table-column>
        <el-table-column
          label="操作">
          <template scope="props">
            <el-button type="info" size="small" icon="edit">修改</el-button>
            <el-button type="danger" size="small" icon="delete" @click="delete_data(props.row)">删除</el-button>
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
        dbData: [],
        //当前页码
        currentPage: 1,
        //数据总条目
        total: 0,
        //每页显示多少条数据
        size: 10,
        //请求时的loading效果
        load_data: true
      }
    },
    components: {
      panelTitle,
      bottomToolBar
    },
    created(){
      this.getSysDBData()
    },
    methods: {
      //刷新
      refresh(){
        this.getSysDBData()
      },
      // 获取数据
      getSysDBData() {
        console.log('get sysdb data...')
        this.load_data = true
        this.$fetch.api_sysdb.list({
          offset: (this.currentPage-1)*this.size,
          size: this.size
        }).then((ret) => {
          this.dbData = ret.data
          this.total = ret.total
          this.load_data = false
        }).catch(() => {
          this.load_data = false
        })
      },
      //单个删除
      delete(item) {
        console.log('delete row...')
      },
      //页码选择
      handleCurrentChange(val) {
        this.currentPage = val
        this.getSysDBData()
      }

    }
  }
</script>
