import fetch from 'common/fetch'
import {port_table} from 'common/port_uri'

//数据列表
export function list(params) {
  return fetch({
    url: port_table.sysdbs,
    method: 'get',
    params
  })
}

//创建导出数据库记录
export function create(data) {
  return fetch({
    url: port_table.sysdbs,
    method: 'post',
    data
  })
}

