// pages/login/login.js
Page({
  data: {
    // 移除邮箱和密码数据
  },

  onWechatLogin(e) {
    // 检查用户是否拒绝授权
    if (e.detail.errMsg !== 'getUserInfo:ok') {
      wx.showToast({
        title: '需要授权才能登录',
        icon: 'none'
      });
      return;
    }

    // 调用wx.login获取登录凭证
    wx.login({
      success: (res) => {
        if (res.code) {
          // 请求后端登录接口
          wx.request({
            url: `${getApp().globalData.serverUrl}/api/auth/login`,
            method: 'POST',
            data: {
              code: res.code
            },
            header: {
              'Content-Type': 'application/json'
            },
            success: (res) => {
              if (res.statusCode === 200) {
                // 登录成功，保存token和用户信息
                wx.setStorageSync('token', res.data.token);
                wx.setStorageSync('userInfo', res.data);

                // 显示登录成功提示
                wx.showToast({
                  title: '登录成功',
                  icon: 'success'
                });

                // 延迟跳转回上一页或首页
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
        } else {
          console.error('获取用户登录态失败:', res.errMsg);
          wx.showToast({
            title: '获取登录信息失败',
            icon: 'none'
          });
        }
      },
      fail: () => {
        wx.showToast({
          title: '微信登录失败',
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