<template>
  <footer class="footer" v-once>Copyright &copy; {{getYear}} gospider</footer>
</template>
<script type="text/javascript">
  export default {
    computed: {
      getYear(){
        return new Date().getFullYear()
      }
    }
  }
</script>
