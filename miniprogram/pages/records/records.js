Page({
  data: {
    records: [],
    filteredRecords: [],
    searchKeyword: '',
    filterCategory: '全部',
    categories: ['全部', '日常', '工作', '学习', '健康', '财务', '其他']
  },

  onLoad() {
    this.loadRecords();
  },

  onShow() {
    this.loadRecords();
  },

  loadRecords() {
    // 从本地存储获取记录
    const records = wx.getStorageSync('records') || [];
    const userInfo = wx.getStorageSync('userInfo');
    
    // 如果未登录，不显示提示，只设置空数组
    if (!userInfo || !userInfo.openId) {
      this.setData({
        records: [],
        filteredRecords: []
      });
      return;
    }
    
    // 过滤当前用户的记录并按时间倒序排列
    const userRecords = records
      .filter(record => record.openId === userInfo.openId)
      .sort((a, b) => new Date(b.createTime) - new Date(a.createTime));
    
    this.setData({
      records: userRecords,
      filteredRecords: userRecords
    });
  },

  onSearchInput(e) {
    const keyword = e.detail.value;
    this.setData({
      searchKeyword: keyword
    });
    this.filterRecords();
  },

  onFilterChange(e) {
    const selectedCategory = this.data.categories[e.detail.value];
    this.setData({
      filterCategory: selectedCategory
    });
    this.filterRecords();
  },

  filterRecords() {
    const { records, searchKeyword, filterCategory } = this.data;
    
    let filtered = records;
    
    // 根据关键词过滤
    if (searchKeyword) {
      filtered = filtered.filter(record => 
        record.title.toLowerCase().includes(searchKeyword.toLowerCase()) ||
        record.content.toLowerCase().includes(searchKeyword.toLowerCase())
      );
    }
    
    // 根据分类过滤
    if (filterCategory && filterCategory !== '全部') {
      filtered = filtered.filter(record => record.category === filterCategory);
    }
    
    this.setData({
      filteredRecords: filtered
    });
  },

  onRecordTap(e) {
    const recordId = e.currentTarget.dataset.id;
    wx.navigateTo({
      url: `/pages/record-detail/record-detail?id=${recordId}`
    });
  },

  onRefresh() {
    this.loadRecords();
  }
})