// script.js (替换或合并现有代码)

document.addEventListener('DOMContentLoaded', () => {

    // --- 模拟数据 ---
    const bangumiList = [
        { id: 1, title: '葬送的芙莉莲', cover: 'https://picsum.photos/seed/frieren/300/400.jpg', isFavorited: false },
        { id: 2, title: '我推的孩子', cover: 'https://picsum.photos/seed/oshi/300/450.jpg', isFavorited: true },
        { id: 3, title: '咒术回战 第二季', cover: 'https://picsum.photos/seed/jjk2/300/420.jpg', isFavorited: false },
        { id: 4, title: '无职转生 ~到了异世界就拿出真本事~', cover: 'https://picsum.photos/seed/mushoku/300/440.jpg', isFavorited: false },
        { id: 5, title: '间谍过家家 第二季', cover: 'https://picsum.photos/seed/spy2/300/400.jpg', isFavorited: true },
        { id: 6, title: '鬼灭之刃 锻刀村篇', cover: 'https://picsum.photos/seed/katana/300/430.jpg', isFavorited: false },
        { id: 7, title: '进击的巨人 最终季', cover: 'https://picsum.photos/seed/aot/300/410.jpg', isFavorited: true },
        { id: 8, title: '赛博朋克：边缘行者', cover: 'https://picsum.photos/seed/cyber/300/400.jpg', isFavorited: false },
        { id: 9, title: '莉可丽丝', cover: 'https://picsum.photos/seed/lycoris/300/450.jpg', isFavorited: false },
        { id: 10, title: '孤独摇滚！', cover: 'https://picsum.photos/seed/bocchi/300/420.jpg', isFavorited: true },
        { id: 11, title: '电锯人', cover: 'https://picsum.photos/seed/chainsaw/300/400.jpg', isFavorited: false },
        { id: 12, title: '辉夜大小姐想让我告白 第三季', cover: 'https://picsum.photos/seed/kaguya3/300/440.jpg', isFavorited: false },
    ];

    // --- DOM 元素获取 ---
    const waterfallContainer = document.getElementById('waterfall-container');
    const loginBtn = document.getElementById('login-btn');

    // 评价弹窗元素
    const ratingModal = document.getElementById('rating-modal');
    const closeRatingModalBtn = document.getElementById('close-rating-modal');
    const ratingBangumiTitle = document.getElementById('rating-bangumi-title');
    const submitRatingBtn = document.getElementById('submit-rating');

    // 登录/注册弹窗元素
    const authModal = document.getElementById('auth-modal');
    const closeAuthModalBtn = document.getElementById('close-auth-modal');
    const tabBtns = document.querySelectorAll('.tab-btn');
    const authForms = document.querySelectorAll('.auth-form');

    // --- 函数定义 ---

    // 渲染瀑布流
    function renderBangumiList() {
        waterfallContainer.innerHTML = ''; // 清空容器
        bangumiList.forEach(bangumi => {
            const card = document.createElement('div');
            card.className = 'bangumi-card';
            card.innerHTML = `
                <img src="${bangumi.cover}" alt="${bangumi.title}">
                <div class="bangumi-info">
                    <h3>${bangumi.title}</h3>
                    <div class="card-actions">
                        <span class="favorite-btn ${bangumi.isFavorited ? 'favorited' : ''}" data-id="${bangumi.id}">
                            <i class="fas fa-heart"></i>
                        </span>
                        <span class="rating-btn" data-id="${bangumi.id}">
                            <i class="fas fa-star"></i> 评价
                        </span>
                    </div>
                </div>
            `;
            waterfallContainer.appendChild(card);
        });
    }

    // --- 事件监听 ---

    // 瀑布流容器事件委托 (处理收藏和评价按钮点击)
    waterfallContainer.addEventListener('click', (e) => {
        const target = e.target;
        const card = target.closest('.bangumi-card');
        if (!card) return;

        const id = parseInt(card.dataset.id || target.closest('[data-id]').dataset.id);
        const bangumi = bangumiList.find(b => b.id === id);

        // 收藏按钮点击
        if (target.closest('.favorite-btn')) {
            bangumi.isFavorited = !bangumi.isFavorited;
            renderBangumiList(); // 重新渲染以更新心形图标
            // TODO: 调用后端API更新收藏状态
            // api.updateFavoriteStatus(id, bangumi.isFavorited);
        }

        // 评价按钮点击
        if (target.closest('.rating-btn')) {
            ratingBangumiTitle.textContent = bangumi.title;
            ratingModal.style.display = 'block';
            // 可以在这里预加载已有的评价数据
        }
    });

    // 登录按钮点击
    loginBtn.addEventListener('click', () => {
        authModal.style.display = 'block';
    });

    // 关闭评价弹窗
    closeRatingModalBtn.addEventListener('click', () => {
        ratingModal.style.display = 'none';
    });

    // 关闭登录/注册弹窗
    closeAuthModalBtn.addEventListener('click', () => {
        authModal.style.display = 'none';
    });

    // 点击弹窗背景关闭弹窗
    window.addEventListener('click', (e) => {
        if (e.target === ratingModal) {
            ratingModal.style.display = 'none';
        }
        if (e.target === authModal) {
            authModal.style.display = 'none';
        }
    });

    // 登录/注册标签页切换
    tabBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const targetTab = btn.dataset.tab;
            
            // 更新按钮状态
            tabBtns.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');

            // 更新表单显示
            authForms.forEach(form => {
                if (form.id === `${targetTab}-form`) {
                    form.classList.add('active');
                } else {
                    form.classList.remove('active');
                }
            });
        });
    });

    // 提交评价
    submitRatingBtn.addEventListener('click', () => {
        const rating = document.querySelector('input[name="rating"]:checked');
        const comment = document.getElementById('rating-comment').value;
        const isPublic = document.getElementById('rating-public').checked;

        if (!rating) {
            alert('请先选择评分！');
            return;
        }

        const ratingData = {
            bangumiTitle: ratingBangumiTitle.textContent,
            score: rating.value,
            comment: comment,
            isPublic: isPublic
        };

        console.log('提交的评价数据:', ratingData);
        alert('评价提交成功！(仅前端模拟)');
        
        // TODO: 调用后端API提交评价
        // api.submitRating(ratingData).then(response => { ... });

        // 关闭弹窗并重置表单
        ratingModal.style.display = 'none';
        document.getElementById('rating-comment').value = '';
        document.getElementById('rating-public').checked = false;
        if(rating) rating.checked = false;
    });

    // --- 初始化 ---
    renderBangumiList();

});

