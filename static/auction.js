document.addEventListener('DOMContentLoaded', function () {
    const auctionList = document.getElementById('auctionList');
    const bidForm = document.getElementById('bidForm');
    const placeBidForm = document.getElementById('placeBidForm');
    const auctionForm = document.getElementById('auctionForm');
    let selectedAuctionId = null;
    let currentMaxBid = 0;

    document.getElementById('logoutButton').addEventListener('click', function () {
        localStorage.removeItem('authToken'); // Удаляем токен авторизации
        window.location.href = '/'; // Перенаправляем на страницу входа (замени на актуальный URL)
    });

    // Загрузка данных пользователя
    function loadUserData() {
        fetch('/users/me', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        })
            .then(response => response.json())
            .then(user => {
                document.getElementById('username').textContent = user.username;
                document.getElementById('userBalance').textContent = user.balance;
            })
            .catch(error => {
                console.error('Ошибка при загрузке данных пользователя:', error);
            });
    }

    // Загрузка аукционов
    function loadAuctions() {
        fetch('/auctions', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        })
            .then(response => response.json())
            .then(auctions => {
                auctionList.innerHTML = ''; // Очищаем список аукционов

                auctions.forEach(auction => {
                    const listItem = document.createElement('li');
                    listItem.classList.add('p-4', 'border', 'rounded', 'flex', 'justify-between', 'items-center');

                    // Форматирование дат
                    const startDate = auction.start_time ? new Date(auction.start_time).toLocaleString() : 'Не указана';
                    const endDate = auction.end_time ? new Date(auction.end_time).toLocaleString() : 'Не указана';

                    listItem.innerHTML = `
                <div class="flex-1">
                    <h1 class="text-xl font-semibold">${auction.title || 'Без названия'}</h1>
                    <p class="text-gray-700">Статус лота: ${auction.status || 'Не указан'}</p>
                    <p class="text-gray-700">Начальная цена: ${auction.starting_price || 'Не указана'} ₽</p>
                    <p class="text-gray-700">Текущая максимальная ставка: ${auction.max_bid || auction.starting_price || 'Не указана'} ₽</p>
                    <p class="text-gray-700">Создатель: ${auction.creator_username || 'Не указан'}</p>
                    <p class="text-gray-700">Победитель: ${auction.winner_username || 'Не выбран'}</p>
                    <p class="text-gray-700">Дата начала: ${startDate}</p>
                    <p class="text-gray-700">Дата окончания: ${endDate}</p>
                </div>
                <button class="bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600" data-id="${auction.id}">Сделать ставку</button>
            `;

                    auctionList.appendChild(listItem);

                    // Обработчик клика по кнопке "Сделать ставку"
                    listItem.querySelector('button').addEventListener('click', function () {
                        selectedAuctionId = auction.id;
                        currentMaxBid = auction.max_bid || auction.starting_price || 0;

                        // Показываем форму ставки
                        bidForm.classList.remove('hidden');
                    });
                });
            })
            .catch(error => console.error('Ошибка при загрузке аукционов:', error));
    }

    // Обработчик ставки
    placeBidForm.addEventListener('submit', function (e) {
        e.preventDefault();
        const bidAmount = parseFloat(document.getElementById('bidAmount').value);

        // Проверка, что ставка больше текущей максимальной ставки
        if (isNaN(bidAmount) || bidAmount <= currentMaxBid) {
            alert('Ставка должна быть больше текущей максимальной ставки.');
            return;
        }

        const bidData = {
            auction_id: selectedAuctionId,
            amount: bidAmount
        };

        fetch('/auction/bid', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            },
            body: JSON.stringify(bidData)
        })
            .then(response => response.json())
            .then(data => {
                if ("ID" in data) { // Если ставка успешно создана, скрыть форму и обновить аукционы
                    alert('Ставка успешно сделана!');
                    bidForm.classList.add('hidden');
                    loadAuctions(); // Перезагружаем список аукционов
                } else {
                    throw data;
                }
            })
            .catch(error => {
                alert(error.message);
            });
    });

    // Обработчик формы создания аукциона
    auctionForm.addEventListener('submit', function (e) {
        e.preventDefault();

        const itemName = document.getElementById('itemName').value;
        const description = document.getElementById('description').value;
        const duration = document.getElementById('duration').value;
        const startingPrice = document.getElementById('startingPrice').value;

        const auctionData = {
            title: itemName,
            description: description,
            duration: parseInt(duration),
            starting_price: parseInt(startingPrice)
        };

        fetch('/auctions', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            },
            body: JSON.stringify(auctionData),
        })
            .then(response => response.json())
            .then(data => {
                alert('Аукцион успешно создан!');
                loadAuctions(); // Перезагружаем список аукционов
            })
            .catch(error => {
                alert('Ошибка при создании аукциона');
            });
    });

    loadUserData();  // Загружаем информацию о пользователе
    loadAuctions();  // Загружаем список аукционов
});
