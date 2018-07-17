import fetch from 'common/fetch'

//数据列表
export function list(params) {
  return fetch({
    url: '/api/tasks',
    method: 'get',
    params
  })
}

//根据id查询数据
export function get(id) {
  return fetch({
    url: '/api/tasks/' + id,
    method: 'get'
  })
}

// 查询所有rules
export function getRules() {
  return fetch({
    url: '/api/rules',
    method: 'get'
  })
}

//根据id停止任务
export function stop(id) {
  return fetch({
    url: '/api/tasks/' + id + '/stop',
    method: 'put'
  })
}
//添加或修改数据
export function save(data) {
  return fetch({
    url: '/api/tasks',
    method: 'post',
    data
  })
}

