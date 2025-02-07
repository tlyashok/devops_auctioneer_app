document.addEventListener('DOMContentLoaded', function () {
    // Переключение между формами регистрации и авторизации
    function toggleForms(formType) {
        const registerForm = document.getElementById('registerForm');
        const loginForm = document.getElementById('loginForm');
        const showRegisterBtn = document.getElementById('showRegister');
        const showLoginBtn = document.getElementById('showLogin');
        const goToRegisterBtn = document.getElementById('goToRegister');
        const goToLoginBtn = document.getElementById('goToLogin');

        if (formType === 'register') {
            registerForm.style.display = 'block';
            loginForm.style.display = 'none';
            showRegisterBtn.style.display = 'none';
            showLoginBtn.style.display = 'none';
            goToRegisterBtn.style.display = 'none';
            goToLoginBtn.style.display = 'block';
        } else {
            registerForm.style.display = 'none';
            loginForm.style.display = 'block';
            showRegisterBtn.style.display = 'none';
            showLoginBtn.style.display = 'none';
            goToRegisterBtn.style.display = 'block';
            goToLoginBtn.style.display = 'none';
        }
    }

    // Функция для авторизации
    function loginUser(username, password) {
        const requestBody = { username, password };

        fetch('/users/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody),
        })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    if (data.error === 'invalid_credentials') {
                        alert('Неверные данные для входа. Проверьте логин и пароль.');
                    } else {
                        alert('Ошибка авторизации. Пожалуйста, попробуйте снова.');
                    }
                } else if (data.token) {
                    alert('Успешная авторизация!');
                    localStorage.setItem('authToken', data.token); // Сохранение токена
                    const decodedToken = jwt_decode(data.token);
                    localStorage.setItem('userId', decodedToken["user_id"]);
                    window.location.href = '/auction.html'; // Перенаправление на страницу аукционов
                } else {
                    alert('Неверные данные');
                }
            })
            .catch(error => {
                alert('Ошибка авторизации: ' + error);
            });
    }


    // Переключение форм по кнопке
    document.getElementById('showRegister').addEventListener('click', function() {
        toggleForms('register');
    });
    document.getElementById('showLogin').addEventListener('click', function() {
        toggleForms('login');
    });

    // Переход на форму регистрации
    document.getElementById('goToRegister').addEventListener('click', function() {
        toggleForms('register');
    });

    // Переход на форму авторизации
    document.getElementById('goToLogin').addEventListener('click', function() {
        toggleForms('login');
    });

    // Обработчик формы регистрации
    const registerForm = document.getElementById('register');
    if (registerForm) {
        registerForm.addEventListener('submit', function (e) {
            e.preventDefault();
            const username = document.getElementById('registerUsername').value;
            const password = document.getElementById('registerPassword').value;
            const balance = parseFloat(document.getElementById('registerBalance').value);

            // Проверка на корректность введенных данных
            if (!username || !password || isNaN(balance) || balance < 0) {
                alert('Пожалуйста, заполните все поля корректно.');
                return;
            }

            const requestBody = {
                username,
                password,
                balance
            };

            fetch('/users/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestBody),
            })
                .then(response => {
                    if (!response.ok) {
                        return response.json().then(errorData => {
                            // Отклоняем промис с ошибкой
                            return Promise.reject(new Error(errorData.message || 'Произошла ошибка!'));
                        });
                    }
                    return response.json();
                })
                .then(data => {
                    alert('Пользователь зарегистрирован!');
                    // После успешной регистрации сразу авторизуем пользователя
                    loginUser(username, password);
                })
                .catch(error => {
                    alert(error);
                });
        });
    }

    // Обработчик формы авторизации
    const loginForm = document.getElementById('login');
    if (loginForm) {
        loginForm.addEventListener('submit', function (e) {
            e.preventDefault();
            const username = document.getElementById('loginUsername').value;
            const password = document.getElementById('loginPassword').value;

            // Проверка на пустые поля
            if (!username || !password) {
                alert('Пожалуйста, заполните все поля.');
                return;
            }

            loginUser(username, password);
        });
    }
});
