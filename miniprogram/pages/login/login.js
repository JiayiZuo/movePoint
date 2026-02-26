// pages/login/login.js
Page({
  data: {
    email: '',
    password: ''
  },

  onInputEmail(e) {
    this.setData({
      email: e.detail.value
    });
  },

  onInputPassword(e) {
    this.setData({
      password: e.detail.value
    });
  },

  onLogin() {
    const { email, password } = this.data;

    // 简单验证
    if (!email || !password) {
      wx.showToast({
        title: '请输入邮箱和密码',
        icon: 'none'
      });
      return;
    }

    // 发起登录请求
    wx.request({
      url: `${getApp().globalData.serverUrl}/login`,
      method: 'POST',
      data: {
        email: email,
        password: password
      },
      header: {
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200) {
          // 登录成功，保存token
          wx.setStorageSync('token', res.data.token);
          
          // 提示登录成功
          wx.showToast({
            title: '登录成功',
            icon: 'success'
          });

          // 返回上一页或跳转到首页
          setTimeout(() => {
            wx.navigateBack({
              delta: 1
            });
            // 或者跳转到首页
            // wx.switchTab({
            //   url: '/pages/index/index'
            // });
          }, 1500);
        } else {
          // 登录失败
          wx.showToast({
            title: res.data.error || '登录失败',
            icon: 'none'
          });
        }
      },
      fail: (err) => {
        console.error('登录请求失败:', err);
        wx.showToast({
          title: '网络错误，请重试',
          icon: 'none'
        });
      }
    });
  },

  navigateToRegister() {
    wx.navigateTo({
      url: '/pages/register/register'
    });
  }
});