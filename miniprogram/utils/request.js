const app = getApp()

function getToken() {
  return wx.getStorageSync('token')
}

function request(options) {
  const token = getToken()
  const header = Object.assign(
    {
      'content-type': 'application/json'
    },
    options.header || {}
  )

  if (token) {
    header.Authorization = `Bearer ${token}`
  }

  return new Promise((resolve, reject) => {
    wx.request({
      url: `${app.globalData.baseURL}${options.url}`,
      method: options.method || 'GET',
      data: options.data || {},
      header,
      success(res) {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data)
          return
        }

        if (res.statusCode === 401) {
          wx.removeStorageSync('token')
          wx.removeStorageSync('user')
          app.globalData.user = null
          wx.reLaunch({ url: '/pages/login/login' })
        }

        const message = (res.data && (res.data.error || res.data.message)) || '请求失败'
        wx.showToast({ title: message, icon: 'none' })
        reject(new Error(message))
      },
      fail(err) {
        wx.showToast({ title: '无法连接服务器', icon: 'none' })
        reject(err)
      }
    })
  })
}

module.exports = {
  request
}
