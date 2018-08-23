import axios from 'axios'
import { Message } from 'element-ui'

export default function fetch(options) {
  return new Promise((resolve, reject) => {
    // 创建一个axios实例
    const instance = axios.create({
      // 设置默认根地址
      baseURL: process.env.BASE_API,
      // 设置请求超时设置
      timeout: 10000
    })
    // 请求处理
    instance(options)
      .then((ret) => {
        // 请求成功时,根据业务判断状态
        resolve(ret.data)
      })
      .catch((error) => {
        // 请求失败时,根据业务判断状态
        if (error.response) {
          const resError = error.response
          const resCode = resError.status
          const resMsg = error.message
          Message.error('操作失败！错误原因 ' + resMsg)
          reject({ code: resCode, msg: resMsg })
        }
      })
  })
}
