import request from '@/utils/request'

// 数据列表
export function fetchTaskList(params) {
  return request({
    url: '/api/tasks',
    method: 'get',
    params
  })
}

// 根据id查询数据
export function getTask(id) {
  return request({
    url: '/api/tasks/' + id,
    method: 'get'
  })
}

// 查询所有rules
export function getRules() {
  return request({
    url: '/api/rules',
    method: 'get'
  })
}

// 根据id停止任务
export function stopTask(id) {
  return request({
    url: '/api/tasks/' + id + '/stop',
    method: 'put'
  })
}
// 根据id启动非定时任务
export function startTask(id) {
  return request({
    url: '/api/tasks/' + id + '/start',
    method: 'put'
  })
}
// 根据id重启定时任务
export function restartTask(id) {
  return request({
    url: '/api/tasks/' + id + '/restart',
    method: 'put'
  })
}
// 添加数据
export function saveTask(data) {
  return request({
    url: '/api/tasks',
    method: 'post',
    data
  })
}

// 修改数据
export function updateTask(id, data) {
  return request({
    url: '/api/tasks/' + id,
    method: 'put',
    data
  })
}

