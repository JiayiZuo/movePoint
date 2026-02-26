// pages/register/register.js
Page({
  data: {
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  },

  onInputUsername(e) {
    this.setData({
      username: e.detail.value
    });
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

  onInputConfirmPassword(e) {
    this.setData({
      confirmPassword: e.detail.value
    });
  },

  onRegister() {
    const { username, email, password, confirmPassword } = this.data;

    // 简单验证
    if (!username || !email || !password || !confirmPassword) {
      wx.showToast({
        title: '请填写所有信息',
        icon: 'none'
      });
      return;
    }

    if (password !== confirmPassword) {
      wx.showToast({
        title: '两次输入密码不一致',
        icon: 'none'
      });
      return;
    }

    if (password.length < 6) {
      wx.showToast({
        title: '密码长度至少6位',
        icon: 'none'
      });
      return;
    }

    // 发起注册请求
    wx.request({
      url: `${getApp().globalData.serverUrl}/register`,
      method: 'POST',
      data: {
        username: username,
        email: email,
        password: password
      },
      header: {
        'Content-Type': 'application/json'
      },
      success: (res) => {
        if (res.statusCode === 200 || res.statusCode === 201) {
          // 注册成功，保存token
          wx.setStorageSync('token', res.data.token);
          
          // 提示注册成功
          wx.showToast({
            title: '注册成功',
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
          // 注册失败
          wx.showToast({
            title: res.data.error || '注册失败',
            icon: 'none'
          });
        }
      },
      fail: (err) => {
        console.error('注册请求失败:', err);
        wx.showToast({
          title: '网络错误，请重试',
          icon: 'none'
        });
      }
    });
  }
});