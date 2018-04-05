<template>
  <div class="panel">
    <panel-title :title="$route.meta.title"></panel-title>
    <div class="panel-body">
      {{msg}}
    </div>
  </div>
</template>
<script type="text/javascript">
  import {panelTitle} from 'components'

  export default{
    data(){
      return {
        msg: 'gospider: 一个提供web任务管理界面的，让使用者只需关心页面规则的爬虫框架。'
      }
    },
    components: {
      panelTitle
    }
  }
</script>
