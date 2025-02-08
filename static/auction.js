document.addEventListener('DOMContentLoaded', function () {
    const auctionList = document.getElementById('auctionList');
    const bidForm = document.getElementById('bidForm');
    const placeBidForm = document.getElementById('placeBidForm');
    const auctionForm = document.getElementById('auctionForm');
    let selectedAuctionId = null;
    let currentMaxBid = 0;

    document.getElementById('logoutButton').addEventListener('click', function () {
        localStorage.removeItem('authToken');
        window.location.href = '/';
    });

    function sanitizeInput(str) {
        return DOMPurify.sanitize(str);
    }

    function loadUserData() {
        fetch('/users/me', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        })
            .then(response => response.json())
            .then(user => {
                document.getElementById('username').textContent = sanitizeInput(user.username);
                document.getElementById('userBalance').textContent = sanitizeInput(user.balance.toString());
            })
            .catch(error => {
                console.error('Ошибка при загрузке данных пользователя:', error);
            });
    }

    function loadAuctions() {
        fetch('/auctions', {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        })
            .then(response => response.json())
            .then(auctions => {
                auctionList.innerHTML = '';
                auctions.forEach(auction => {
                    const listItem = document.createElement('li');
                    listItem.classList.add('p-4', 'border', 'rounded', 'flex', 'justify-between', 'items-center');

                    const startDate = auction.start_time ? new Date(auction.start_time).toLocaleString() : 'Не указана';
                    const endDate = auction.end_time ? new Date(auction.end_time).toLocaleString() : 'Не указана';

                    listItem.innerHTML = `
                            <div class="flex-1">
                                <h1 class="text-xl font-semibold">${sanitizeInput(auction.title || 'Без названия')}</h1>
                                <p class="text-gray-700">Статус: ${sanitizeInput(auction.status || 'Не указан')}</p>
                                <p class="text-gray-700">Начальная цена: ${sanitizeInput(auction.starting_price?.toString() || 'Не указана')} ₽</p>
                                <p class="text-gray-700">Макс. ставка: ${sanitizeInput(auction.max_bid?.toString() || auction.starting_price?.toString() || 'Не указана')} ₽</p>
                                <p class="text-gray-700">Создатель: ${sanitizeInput(auction.creator_username || 'Не указан')}</p>
                                <p class="text-gray-700">Победитель: ${sanitizeInput(auction.winner_username || 'Не выбран')}</p>
                                <p class="text-gray-700">Дата начала: ${startDate}</p>
                                <p class="text-gray-700">Дата окончания: ${endDate}</p>
                            </div>
                            <button class="bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600" data-id="${auction.id}">Сделать ставку</button>
                        `;

                    auctionList.appendChild(listItem);

                    listItem.querySelector('button').addEventListener('click', function () {
                        selectedAuctionId = auction.id;
                        currentMaxBid = auction.max_bid || auction.starting_price || 0;
                        bidForm.classList.remove('hidden');
                    });
                });
            })
            .catch(error => console.error('Ошибка при загрузке аукционов:', error));
    }

    placeBidForm.addEventListener('submit', function (e) {
        e.preventDefault();
        const bidAmount = parseFloat(document.getElementById('bidAmount').value);

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
                if ("ID" in data) {
                    alert('Ставка успешно сделана!');
                    bidForm.classList.add('hidden');
                    loadAuctions();
                } else {
                    throw data;
                }
            })
            .catch(error => {
                alert(error.message);
            });
    });

    auctionForm.addEventListener('submit', function (e) {
        e.preventDefault();

        const itemName = sanitizeInput(document.getElementById('itemName').value);
        const description = sanitizeInput(document.getElementById('description').value);
        const duration = parseInt(document.getElementById('duration').value, 10);
        const startingPrice = parseInt(document.getElementById('startingPrice').value, 10);

        const auctionData = {
            title: itemName,
            description: description,
            duration: duration,
            starting_price: startingPrice
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
                loadAuctions();
            })
            .catch(error => {
                alert('Ошибка при создании аукциона');
            });
    });

    loadUserData();
    loadAuctions();
});