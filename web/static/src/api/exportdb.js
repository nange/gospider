import fetch from '@/utils/fetch'

// 数据列表
export function fetchExportDbList(params) {
  return fetch({
    url: '/api/sysdbs',
    method: 'get',
    params
  })
}

// 创建导出数据库记录
export function createExportDb(data) {
  return fetch({
    url: '/api/sysdbs',
    method: 'post',
    data
  })
}

