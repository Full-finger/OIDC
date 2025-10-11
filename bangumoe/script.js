// script.js

document.addEventListener('DOMContentLoaded', () => {
    // --- DOM 元素获取 ---
    const navLinks = document.querySelectorAll('.nav-link');
    const contentSections = document.querySelectorAll('.content-section');
    const loginBtn = document.getElementById('login-btn');
    const userStatus = document.getElementById('user-status');
    
    // 弹窗元素
    const authModal = document.getElementById('auth-modal');
    const closeAuthModalBtn = document.getElementById('close-auth-modal');
    const tabBtns = document.querySelectorAll('.tab-btn');
    const authForms = document.querySelectorAll('.auth-form');
    
    const collectionModal = document.getElementById('collection-modal');
    const closeCollectionModalBtn = document.getElementById('close-collection-modal');
    const collectionForm = document.getElementById('collection-form');
    const collectionModalTitle = document.getElementById('collection-modal-title');
    
    // 收藏页面元素
    const syncCollectionsBtn = document.getElementById('sync-collections-btn');
    const addCollectionBtn = document.getElementById('add-collection-btn');
    const collectionTypeFilter = document.getElementById('collection-type-filter');
    const ratingFilter = document.getElementById('rating-filter');
    const collectionsContainer = document.getElementById('collections-container');
    
    // Bangumi绑定页面元素
    const bangumiBindingStatus = document.getElementById('bangumi-binding-status');
    const bangumiBindingActions = document.getElementById('bangumi-binding-actions');
    const syncBangumiDataBtn = document.getElementById('sync-bangumi-data-btn');
    const syncResult = document.getElementById('sync-result');
    
    // 表单元素
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    
    // --- 状态管理 ---
    let currentUser = null;
    let currentCollection = null;
    let bangumiBinding = null;
    let collections = [];
    
    // --- 工具函数 ---
    function showSection(sectionId) {
        // 隐藏所有内容区
        contentSections.forEach(section => {
            section.classList.remove('active');
        });
        
        // 显示指定内容区
        document.getElementById(`${sectionId}-section`).classList.add('active');
        
        // 更新导航链接状态
        navLinks.forEach(link => {
            link.classList.remove('active');
            if (link.dataset.tab === sectionId) {
                link.classList.add('active');
            }
        });
    }
    
    function showModal(modal) {
        modal.style.display = 'block';
    }
    
    function hideModal(modal) {
        modal.style.display = 'none';
    }
    
    function showMessage(message, type = 'info') {
        // 创建消息元素
        const messageEl = document.createElement('div');
        messageEl.className = `message message-${type}`;
        messageEl.textContent = message;
        
        // 添加到页面
        document.body.appendChild(messageEl);
        
        // 3秒后自动移除
        setTimeout(() => {
            document.body.removeChild(messageEl);
        }, 3000);
    }
    
    // --- API 调用函数 ---
    const api = {
        // 用户认证相关
        async login(username, password) {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({
                        user: { id: 1, username: username },
                        token: 'fake-jwt-token'
                    });
                }, 500);
            });
        },
        
        async register(username, email, password) {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({
                        user: { id: 1, username: username, email: email },
                        token: 'fake-jwt-token'
                    });
                }, 500);
            });
        },
        
        // 收藏相关
        async getCollections(type = '', rating = '') {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    const mockCollections = [
                        {
                            id: 1,
                            anime_id: 101,
                            type: 'watching',
                            rating: 9,
                            comment: '非常好看，推荐！',
                            anime: {
                                id: 101,
                                title: '我的英雄学院',
                                episode_count: 12
                            }
                        },
                        {
                            id: 2,
                            anime_id: 102,
                            type: 'completed',
                            rating: 8,
                            comment: '剧情不错，但结局有点仓促',
                            anime: {
                                id: 102,
                                title: '鬼灭之刃',
                                episode_count: 24
                            }
                        },
                        {
                            id: 3,
                            anime_id: 103,
                            type: 'wish',
                            rating: null,
                            comment: '',
                            anime: {
                                id: 103,
                                title: '进击的巨人 最终季',
                                episode_count: 16
                            }
                        }
                    ];
                    
                    // 根据筛选条件过滤
                    let filtered = mockCollections;
                    if (type) {
                        filtered = filtered.filter(c => c.type === type);
                    }
                    if (rating) {
                        filtered = filtered.filter(c => c.rating == rating);
                    }
                    
                    resolve(filtered);
                }, 300);
            });
        },
        
        async upsertCollection(collectionData) {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({
                        id: collectionData.id || Date.now(),
                        ...collectionData
                    });
                }, 300);
            });
        },
        
        async deleteCollection(collectionId) {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({ success: true });
                }, 300);
            });
        },
        
        // Bangumi绑定相关
        async getBangumiBinding() {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    // 模拟已绑定状态
                    resolve({
                        id: 1,
                        user_id: 1,
                        bangumi_user_id: 12345,
                        token_expires_at: '2025-12-31 23:59:59',
                        created_at: '2025-01-01 00:00:00',
                        updated_at: '2025-01-01 00:00:00'
                    });
                }, 300);
            });
        },
        
        async bindBangumiAccount() {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({
                        message: 'Bangumi账户绑定成功',
                        bangumi_user_id: 12345,
                        username: 'test_user',
                        nickname: '测试用户'
                    });
                }, 500);
            });
        },
        
        async unbindBangumiAccount() {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({ message: 'Bangumi账户解绑成功' });
                }, 300);
            });
        },
        
        async syncBangumiData() {
            // 模拟API调用
            return new Promise((resolve) => {
                setTimeout(() => {
                    resolve({
                        new_animes: 3,
                        updated_collections: 5,
                        total_collections: 8
                    });
                }, 1000);
            });
        }
    };
    
    // --- UI 渲染函数 ---
    function renderCollections(collections) {
        collectionsContainer.innerHTML = '';
        
        if (collections.length === 0) {
            collectionsContainer.innerHTML = '<p>暂无收藏记录</p>';
            return;
        }
        
        collections.forEach(collection => {
            const card = document.createElement('div');
            card.className = 'collection-card';
            card.innerHTML = `
                <div class="collection-header">
                    <h3>${collection.anime.title}</h3>
                    <span class="collection-type ${collection.type}">${getCollectionTypeText(collection.type)}</span>
                </div>
                <div class="collection-body">
                    <p>集数: ${collection.anime.episode_count}</p>
                    ${collection.rating ? `<p>评分: <span class="rating-display">${'★'.repeat(collection.rating)}${'☆'.repeat(10 - collection.rating)}</span> (${collection.rating}/10)</p>` : ''}
                    ${collection.comment ? `<p>评论: ${collection.comment}</p>` : ''}
                </div>
                <div class="collection-footer">
                    <button class="btn btn-secondary edit-collection-btn" data-id="${collection.id}">
                        <i class="fas fa-edit"></i> 编辑
                    </button>
                    <button class="btn btn-danger delete-collection-btn" data-id="${collection.id}">
                        <i class="fas fa-trash"></i> 删除
                    </button>
                </div>
            `;
            collectionsContainer.appendChild(card);
        });
    }
    
    function getCollectionTypeText(type) {
        const typeMap = {
            'watching': '在看',
            'completed': '看过',
            'wish': '想看',
            'on_hold': '搁置',
            'dropped': '抛弃'
        };
        return typeMap[type] || type;
    }
    
    function renderBangumiBindingStatus() {
        if (bangumiBinding) {
            bangumiBindingStatus.innerHTML = `
                <p>✅ 已绑定到 Bangumi 账户</p>
                <p>Bangumi 用户ID: ${bangumiBinding.bangumi_user_id}</p>
                <p>绑定时间: ${new Date(bangumiBinding.created_at).toLocaleString()}</p>
            `;
            
            bangumiBindingActions.innerHTML = `
                <button id="unbind-bangumi-btn" class="btn btn-danger">
                    <i class="fas fa-unlink"></i> 解绑 Bangumi 账户
                </button>
            `;
            
            document.getElementById('unbind-bangumi-btn').addEventListener('click', handleUnbindBangumi);
        } else {
            bangumiBindingStatus.innerHTML = `
                <p>❌ 未绑定到 Bangumi 账户</p>
            `;
            
            bangumiBindingActions.innerHTML = `
                <button id="bind-bangumi-btn" class="btn btn-primary">
                    <i class="fas fa-link"></i> 绑定 Bangumi 账户
                </button>
            `;
            
            document.getElementById('bind-bangumi-btn').addEventListener('click', handleBindBangumi);
        }
    }
    
    // --- 事件处理函数 ---
    function handleNavClick(e) {
        e.preventDefault();
        const sectionId = e.target.dataset.tab;
        showSection(sectionId);
        
        // 根据不同页面执行初始化操作
        switch (sectionId) {
            case 'collections':
                loadCollections();
                break;
            case 'bangumi':
                loadBangumiBinding();
                break;
        }
    }
    
    function handleLoginClick() {
        showModal(authModal);
    }
    
    function handleCloseAuthModal() {
        hideModal(authModal);
    }
    
    function handleTabClick(e) {
        const targetTab = e.target.dataset.tab;
        
        // 更新按钮状态
        tabBtns.forEach(btn => btn.classList.remove('active'));
        e.target.classList.add('active');
        
        // 更新表单显示
        authForms.forEach(form => {
            if (form.id === `${targetTab}-form`) {
                form.classList.add('active');
            } else {
                form.classList.remove('active');
            }
        });
    }
    
    async function handleLoginFormSubmit(e) {
        e.preventDefault();
        const username = document.getElementById('login-username').value;
        const password = document.getElementById('login-password').value;
        
        try {
            const result = await api.login(username, password);
            currentUser = result.user;
            userStatus.textContent = `欢迎, ${currentUser.username}`;
            hideModal(authModal);
            showMessage('登录成功', 'success');
            
            // 重置表单
            loginForm.reset();
        } catch (error) {
            showMessage('登录失败: ' + error.message, 'error');
        }
    }
    
    async function handleRegisterFormSubmit(e) {
        e.preventDefault();
        const username = document.getElementById('register-username').value;
        const email = document.getElementById('register-email').value;
        const password = document.getElementById('register-password').value;
        const confirmPassword = document.getElementById('register-confirm-password').value;
        
        if (password !== confirmPassword) {
            showMessage('两次输入的密码不一致', 'error');
            return;
        }
        
        try {
            const result = await api.register(username, email, password);
            currentUser = result.user;
            userStatus.textContent = `欢迎, ${currentUser.username}`;
            hideModal(authModal);
            showMessage('注册成功', 'success');
            
            // 重置表单
            registerForm.reset();
        } catch (error) {
            showMessage('注册失败: ' + error.message, 'error');
        }
    }
    
    async function loadCollections() {
        try {
            const type = collectionTypeFilter.value;
            const rating = ratingFilter.value;
            collections = await api.getCollections(type, rating);
            renderCollections(collections);
        } catch (error) {
            showMessage('加载收藏列表失败: ' + error.message, 'error');
        }
    }
    
    function handleAddCollection() {
        currentCollection = null;
        collectionModalTitle.textContent = '添加收藏';
        collectionForm.reset();
        showModal(collectionModal);
    }
    
    function handleEditCollection(e) {
        const collectionId = parseInt(e.target.closest('.edit-collection-btn').dataset.id);
        const collection = collections.find(c => c.id === collectionId);
        
        if (collection) {
            currentCollection = collection;
            collectionModalTitle.textContent = '编辑收藏';
            document.getElementById('anime-id').value = collection.anime_id;
            document.getElementById('collection-type').value = collection.type;
            document.getElementById('collection-rating').value = collection.rating || '';
            document.getElementById('collection-comment').value = collection.comment || '';
            showModal(collectionModal);
        }
    }
    
    async function handleCollectionFormSubmit(e) {
        e.preventDefault();
        
        const collectionData = {
            anime_id: parseInt(document.getElementById('anime-id').value),
            type: document.getElementById('collection-type').value,
            rating: document.getElementById('collection-rating').value || null,
            comment: document.getElementById('collection-comment').value
        };
        
        if (currentCollection) {
            collectionData.id = currentCollection.id;
        }
        
        try {
            await api.upsertCollection(collectionData);
            hideModal(collectionModal);
            showMessage(currentCollection ? '收藏更新成功' : '收藏添加成功', 'success');
            loadCollections(); // 重新加载收藏列表
        } catch (error) {
            showMessage((currentCollection ? '更新' : '添加') + '收藏失败: ' + error.message, 'error');
        }
    }
    
    async function handleDeleteCollection(e) {
        const collectionId = parseInt(e.target.closest('.delete-collection-btn').dataset.id);
        
        if (confirm('确定要删除这个收藏吗？')) {
            try {
                await api.deleteCollection(collectionId);
                showMessage('收藏删除成功', 'success');
                loadCollections(); // 重新加载收藏列表
            } catch (error) {
                showMessage('删除收藏失败: ' + error.message, 'error');
            }
        }
    }
    
    async function loadBangumiBinding() {
        try {
            bangumiBinding = await api.getBangumiBinding();
            renderBangumiBindingStatus();
        } catch (error) {
            bangumiBinding = null;
            renderBangumiBindingStatus();
        }
    }
    
    async function handleBindBangumi() {
        try {
            const result = await api.bindBangumiAccount();
            showMessage(result.message, 'success');
            loadBangumiBinding();
        } catch (error) {
            showMessage('绑定失败: ' + error.message, 'error');
        }
    }
    
    async function handleUnbindBangumi() {
        if (confirm('确定要解绑 Bangumi 账户吗？')) {
            try {
                const result = await api.unbindBangumiAccount();
                showMessage(result.message, 'success');
                bangumiBinding = null;
                renderBangumiBindingStatus();
            } catch (error) {
                showMessage('解绑失败: ' + error.message, 'error');
            }
        }
    }
    
    async function handleSyncBangumiData() {
        syncResult.className = 'sync-result';
        syncResult.textContent = '正在同步数据...';
        syncResult.classList.add('info');
        syncResult.style.display = 'block';
        
        try {
            const result = await api.syncBangumiData();
            syncResult.className = 'sync-result success';
            syncResult.innerHTML = `
                <p>✅ 数据同步完成</p>
                <p>新增番剧: ${result.new_animes}</p>
                <p>更新收藏: ${result.updated_collections}</p>
                <p>总收藏数: ${result.total_collections}</p>
            `;
        } catch (error) {
            syncResult.className = 'sync-result error';
            syncResult.textContent = '❌ 数据同步失败: ' + error.message;
        }
    }
    
    // --- 事件监听 ---
    // 导航链接点击
    navLinks.forEach(link => {
        link.addEventListener('click', handleNavClick);
    });
    
    // 登录按钮点击
    loginBtn.addEventListener('click', handleLoginClick);
    
    // 关闭弹窗按钮
    closeAuthModalBtn.addEventListener('click', handleCloseAuthModal);
    closeCollectionModalBtn.addEventListener('click', () => hideModal(collectionModal));
    
    // 点击弹窗背景关闭弹窗
    window.addEventListener('click', (e) => {
        if (e.target === authModal) {
            hideModal(authModal);
        }
        if (e.target === collectionModal) {
            hideModal(collectionModal);
        }
    });
    
    // 标签页切换
    tabBtns.forEach(btn => {
        btn.addEventListener('click', handleTabClick);
    });
    
    // 表单提交
    loginForm.addEventListener('submit', handleLoginFormSubmit);
    registerForm.addEventListener('submit', handleRegisterFormSubmit);
    collectionForm.addEventListener('submit', handleCollectionFormSubmit);
    
    // 收藏页面事件
    syncCollectionsBtn.addEventListener('click', loadCollections);
    addCollectionBtn.addEventListener('click', handleAddCollection);
    collectionTypeFilter.addEventListener('change', loadCollections);
    ratingFilter.addEventListener('change', loadCollections);
    
    // 收藏卡片事件委托
    collectionsContainer.addEventListener('click', (e) => {
        if (e.target.closest('.edit-collection-btn')) {
            handleEditCollection(e);
        } else if (e.target.closest('.delete-collection-btn')) {
            handleDeleteCollection(e);
        }
    });
    
    // Bangumi绑定页面事件
    syncBangumiDataBtn.addEventListener('click', handleSyncBangumiData);
    
    // --- 初始化 ---
    showSection('home');
});