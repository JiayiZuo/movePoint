App({
  globalData: {
    baseURL: 'http://127.0.0.1:8080/api',
    user: null
  },

  onLaunch() {
    const user = wx.getStorageSync('user')
    if (user) {
      this.globalData.user = user
    }
  }
})
