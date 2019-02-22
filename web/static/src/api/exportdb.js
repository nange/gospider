import request from '@/utils/request'

// 数据列表
export function fetchExportDBList(params) {
  return request({
    url: '/api/exportdb',
    method: 'get',
    params
  })
}

// 创建导出数据库记录
export function createExportDB(data) {
  return request({
    url: '/api/exportdb',
    method: 'post',
    data
  })
}

// 删除导出数据库记录
export function deleteExportDB(params) {
  return request({
    url: '/api/exportdb/' + params.id,
    method: 'delete'
  })
}

