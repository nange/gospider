<template>
  <div class="app-container">
    <div class="filter-container fr">
      <el-button class="filter-item" @click="refresh" icon="el-icon-refresh" size="small"></el-button>
      <router-link :to="{name: 'addExpDb'}" tag="span">
        <el-button class="filter-item" @click="handleCreate" type="primary" icon="el-icon-edit" size="small">{{$t('exportdb.add')}}</el-button>
      </router-link>
    </div>
    <el-table :data="dbData" v-loading="load_data" border fit highlight-current-row
  style="width: 100%;">
        <el-table-column prop="id" :label="$t('exportdb.id')" width="100">
        </el-table-column>
        <el-table-column prop="show_name" :label="$t('exportdb.showname')" width="200">
        </el-table-column>
        <el-table-column prop="host" :label="$t('exportdb.host')" width="200">
        </el-table-column>
        <el-table-column prop="port" :label="$t('exportdb.port')" width="100">
        </el-table-column>
        <el-table-column prop="user" :label="$t('exportdb.user')" width="100">
        </el-table-column>
        <el-table-column prop="db_name" :label="$t('exportdb.dbname')" width="150">
        </el-table-column>
        <el-table-column :label="$t('exportdb.actions')">
          <template scope="props">
            <el-button type="warning" size="small" icon="edit">{{$t('exportdb.edit')}}</el-button>
            <el-button type="danger" size="small" icon="delete" @click="deletedb(props.row)">{{$t('exportdb.delete')}}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container fr">
      <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="currentPage" :page-sizes="[10,20,30, 50]" :page-size="size" layout="total, sizes, prev, pager, next, jumper" :total="total">
      </el-pagination>
    </div>

  </div>
</template>

<script>
  import { fetchExportDBList, deleteExportDB } from '@/api/exportdb'
  import waves from '@/directive/waves' // 水波纹指令
  export default{
    directives: {
      waves
    },
    data() {
      return {
        dbData: [],
        // 当前页码
        currentPage: 1,
        // 数据总条目
        total: 0,
        // 每页显示多少条数据
        size: 10,
        // 请求时的loading效果
        load_data: true
      }
    },
    created() {
      this.getExportDBData()
    },
    methods: {
      // 刷新
      refresh() {
        this.getExportDBData()
      },
      // 获取数据
      getExportDBData() {
        this.load_data = true
        fetchExportDBList({
          offset: (this.currentPage - 1) * this.size,
          size: this.size
        }).then((response) => {
          const data = response.data
          this.dbData = data.data
          this.total = data.total
          this.load_data = false
        }).catch(() => {
          this.load_data = false
        })
      },
      // 单个删除
      deletedb(item) {
        this.$confirm('The operating will delete this record, continue?', 'NOTICE', {
          confirmButtonText: 'Yes',
          cancelButtonText: 'Cancel',
          type: 'warning'
        }).then(() => {
          deleteExportDB({
            id: item.id
          }).then(() => {
            this.$message.success('delete export db success!')
            this.refresh()
          }).catch(() => {
            this.$message.error('delete export db failed!')
          })
        })
      },
      handleSizeChange(val) {
        this.size = val
        this.getExportDBData()
      },
      // 页码选择
      handleCurrentChange(val) {
        this.currentPage = val
        this.getExportDBData()
      }

    }
  }
</script>
