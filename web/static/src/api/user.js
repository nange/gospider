import fetch from 'common/fetch'
import {port_user} from 'common/port_uri'

//登录
export function login(data) {
  return new Promise((resolve, reject) => {
    if (data.username === 'admin' && data.password === 'admin') {
      resolve({
        data: {
          name: 'admin',
          avatar: ''
        },
        msg: '登陆成功!'
      })
    } else {
      reject()
    }
  })
  
  // return fetch({
  //   url: port_user.login,
  //   method: 'post',
  //   data
  // })
}
//登出
export function logout() {
  return new Promise((resolve, reject) => {
    resolve({
      msg: '登出成功!'
    })
  })
  // return fetch({
  //   url: port_user.logout,
  //   method: 'post'
  // })
}
