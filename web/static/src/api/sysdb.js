import fetch from 'common/fetch'

//数据列表
export function list(params) {
  return fetch({
    url: '/api/sysdbs',
    method: 'get',
    params
  })
}

//创建导出数据库记录
export function create(data) {
  return fetch({
    url: '/api/sysdbs',
    method: 'post',
    data
  })
}

