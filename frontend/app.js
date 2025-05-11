let socket = null;
let currentUser = "";

function signup() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    fetch("http://localhost:8080/signup", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
    })
    .then(response => {
        if (response.ok) {
            alert("Signup successful!");
        } else {
            alert("Signup failed!");
        }
    });
}

function login() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    fetch("http://localhost:8080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            alert("Login successful!");
            // optionally show chat box here
        } else {
            alert("Login failed!");
        }
    });
}


function connectWebSocket() {
    socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
        socket.send(JSON.stringify(currentUser));
    };

    socket.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        const messagesDiv = document.getElementById("messages");
        const msgEl = document.createElement("div");
        msgEl.classList.add("message");
        msgEl.innerText = `${msg.sender} ➜ ${msg.receiver}: ${msg.content}`;
        messagesDiv.appendChild(msgEl);
    };
}

function sendMessage() {
    const receiver = document.getElementById("receiver").value;
    const content = document.getElementById("message").value;

    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ receiver, content }));
        document.getElementById("message").value = "";
    } else {
        alert("WebSocket not connected");
    }
}
