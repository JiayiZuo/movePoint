const api = require('../../services/api')

const app = getApp()

Page({
  data: {
    loading: false,
    form: {
      username: '',
      email: '',
      password: '',
      birth_date: ''
    }
  },

  onInput(e) {
    this.setData({ [`form.${e.currentTarget.dataset.field}`]: e.detail.value })
  },

  onDateChange(e) {
    this.setData({ 'form.birth_date': e.detail.value })
  },

  async submit() {
    const form = this.data.form
    if (!form.username || !form.email || !form.password) {
      wx.showToast({ title: '请填写用户名、邮箱和密码', icon: 'none' })
      return
    }

    this.setData({ loading: true })
    try {
      const res = await api.register(form)
      const user = {
        id: res.user_id,
        username: res.username,
        email: res.email
      }
      wx.setStorageSync('token', res.token)
      wx.setStorageSync('user', user)
      app.globalData.user = user
      wx.switchTab({ url: '/pages/dashboard/dashboard' })
    } finally {
      this.setData({ loading: false })
    }
  }
})
