// main.js

if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/service-worker.js')
    .then(function(registration) {
        console.log('Service Worker 注册成功:', registration);
    })
    .catch(function(error) {
        console.log('Service Worker 注册失败:', error);
    });
}

document.getElementById('subscribe').addEventListener('click', function() {
    Notification.requestPermission().then(function(permission) {
        if (permission === 'granted') {
            navigator.serviceWorker.ready.then(function(registration) {
                registration.pushManager.subscribe({
                    userVisibleOnly: true,
                    applicationServerKey: 'BFl0cadQjaCPWg6EACPAfBgsDyBDorLNGZM-IEUiz4qCUOMtV611r9A-5RwkQo6yq82NOA7NH93YDKmbrDMh4Bg'
                }).then(function(subscription) {
                    console.log('用户订阅成功:', subscription);
                    // 将订阅对象转换为 JSON 字符串
                    const subscriptionJson = JSON.stringify(subscription);
                    console.log('订阅对象 JSON 字符串:', subscriptionJson);
                    // 更新页面上的内容
                    document.getElementById('subscription').textContent = subscriptionJson;
                    // 将订阅对象发送到服务器
                    fetch('http://localhost:8081/subscribe', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: subscriptionJson
                    }).then(function(response) {
                        if (response.ok) {
                            console.log('订阅对象已发送到服务器');
                        } else {
                            console.log('发送订阅对象到服务器失败');
                        }
                    }).catch(function(error) {
                        console.log('发送订阅对象到服务器时出错:', error);
                    });
                }).catch(function(error) {
                    console.log('用户订阅失败:', error);
                });
            });
        } else {
            console.log('用户拒绝接收通知');
        }
    });
});

document.getElementById('pushMessage').addEventListener('click', function() {
    fetch('http://localhost:8081/pushMessage', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            title: 'Push API',
            body: '欢迎使用 Push API 示例',
        })
    })
    
});