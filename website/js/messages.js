const usrsContenar = document.getElementById("users_contenar");

let socket;
let isMoerUsers = true
let usersDisplay = []

function msessagesWS() {
  if (socket && socket.readyState === WebSocket.OPEN) return;

  socket = new WebSocket("/ws");

  socket.onopen = () => {
    console.log("successfully conected");
  };
  socket.onclose = () => {
    // try open agein
    setTimeout(msessagesWS, 3000);
  };

  socket.onerror = (err) => {
    // console.log("socket error: ", err);
  };

  socket.onmessage = (msg) => {

    const data = JSON.parse(msg.data);

    if (data.type === "typing") {
      showTypingIndicator(data.sender)
    } else if (data.type === "error") {
      createAlert(alert("alert-danger", data.message));
    } else if (data.type === "user_status") {
      displayOnline(data);
    } else {
      if (data.sender === userOpenChat) {
        addMessageInChat(data)
      } else {
        createAlert(alert("alert-info", `${data.sender}: ${data.text}`));
        addNotifNumber(data.sender)
      }
    }
  };
}
// typing indicators*********************************************
let loop = {};
async function showTypingIndicator(user) {

  const chatContainer = document.getElementById(`contaner_${user}`)
  if (!chatContainer) return;
  
  let typing = chatContainer.querySelector("#typinggg")
  typing.classList.remove("none")

  if (loop[user]) {
    clearTimeout(loop[user]);
  }

  // Set a timeout to hide the typing indicator after 2 seconds of inactivity
  loop[user] = setTimeout(() => {
    typing.classList.add("none");
  }, 1000);

}

// TYPING===============*************************

async function sendMessage(btn) {
  let name = null;
  try {
    const { state } = await window.checkIfLoggedIn();
    if (!state) {
      createAlert(alert("alert-danger", "you need to login!!!"));
      return;
    }

    let messageInput = document.getElementById("messageInput");
    name = btn.getAttribute("sendto");
    if (!name) {
      createAlert(alert("alert-danger", "error select user"));
      return;
    }
    if (messageInput.value.length > 160) {
      createAlert(alert("alert-danger", "message too long max 160"));
      return;
    }
    if (!messageInput.value.trim()) return;

    // Send message
    socket.send(
      JSON.stringify({
        text: messageInput.value,
        receiver: name,
        timestamp: new Date(),
      })
    );

    // Stop typing notification since the message is sent**************************************************************
    sendTypingNotification(name, false);
  } catch (err) {
    createAlert(alert("alert-danger", "failed to send message"));
    console.error(err);
    return;
  }

  addMyMessage(messageInput.value);
  const contanerUser = document.getElementById(`contaner_${name}`);
  const newUser = document.createElement("div");
  newUser.id = `contaner_${name}`;
  newUser.classList.add("messageProfile", "selectUser");
  newUser.setAttribute("onclick", "clickusers(this)");

  newUser.innerHTML = contanerUser.innerHTML;
  usrsContenar.insertBefore(newUser, usrsContenar.firstChild);
  lastSelect = usersContenar.firstChild.id;
  usrsContenar.removeChild(contanerUser);
  messageInput.value = "";
}

function sendTypingNotification(receiver, isTyping) {
  if (!socket || socket.readyState !== WebSocket.OPEN) return;
  socket.send(
    JSON.stringify({
      type: "typing",
      receiver: receiver,
      typing: isTyping,
    })
  );
}

document.getElementById("messageInput").addEventListener("input", function () {
  const receiver = document.querySelector("[sendto]")?.getAttribute("sendto");
  if (!receiver) return;

  //Send "typing" status
  sendTypingNotification(receiver, true);

});


// *******************************************************************************************


function closeWS() {
  socket.close();
}
function end(txt) {
  return `<div class="txt">${txt}</div>`
}

async function getusers(first) {
  if (!isMoerUsers) {
    return
  }

  let lastuser = ""
  let users = document.querySelectorAll('.messageImg')
  if (users.length !== 0, !first) {
    lastuser = users[users.length - 1].id.slice(5);
  }

  try {
    const response = await fetch("/api/getUsers", {
      method: "POST",
      body: JSON.stringify({ lastuser: lastuser }),
    });

    if (!response.ok) {
      dangerError(response.status);
      throw new Error("failed to fetch users");
    }

    const data = await response.json();

    if (!data || (data.length === 1 && data[0].nickname === window.username)) {
      isMoerUsers = false
      usrsContenar.innerHTML += end("no more usres");
      return
    }
    let allUsers = "";
    for (let usr of data) {
      if (usr.nickname !== window.username && !usersDisplay.includes(usr.nickname)) {
        allUsers += usertemplate(usr);
      }
    }

    usrsContenar.innerHTML += allUsers;
  } catch (err) {
    createAlert(alert("alert-danger", "failed to get users"));
    console.error(err);
  }
}

function usertemplate(usr) {
  let clas = "messageImg";
  if (usr.online) {
    clas += " greenOnline";
  }

  return `<div id="contaner_${usr.nickname}" class="messageProfile" onclick="clickusers(this)">
            <img id="user_${usr.nickname}" class="${clas}" src="/public/img/default.avif" alt="">
            <div class="names">           
                <h3>${usr.nickname}</h3>
                <p>${usr.firstname} ${usr.lastname}</p>
            </div>
            <label id="label_${usr.nickname}"></label>
            <div id="typinggg" class="none">
              <div class="loading none">
                <div></div>
                <div></div>
                <div></div>
              </div>
            </div>
        </div>`;
}

function displayOnline(usr) {
  let usrImg = document.getElementById("user_" + usr.nickname);

  if (!usrImg) return;

  if (usr.connected) {
    usrImg.classList.add("greenOnline");
  } else {
    usrImg.classList.remove("greenOnline");
  }
}
function debounce(func, delay) {
  let timeoutId;
  return function (...args) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => func.apply(this, args), delay);
  };
}
//---------------------------------------------------- scrll users 
const usersContenar = document.getElementById("users_contenar");

const usersScroll = debounce(async () => {
  const { scrollTop, scrollHeight, clientHeight } = usersContenar;
  if (scrollTop + clientHeight >= scrollHeight - 10) {
    await getusers();
  }
}, 200);
usersContenar.addEventListener('scroll', usersScroll)



async function lastUsersChat() {
  try {
    const response = await fetch('/api/lastUsersChat', {
      method: 'POST',
      headers: {
        "Content-Type": "application/json",
      },
    })
    if (!response.ok) {
      dangerError(response.status);
      throw new Error("failed to fetch users");
    }

    let data = await response.json()
    if (!data) {
      return
    }

    let allUsers = "";
    for (let usr of data) {
      usersDisplay.push(usr.nickname)
      allUsers += usertemplate(usr);
    }
    usrsContenar.innerHTML = allUsers;
    usrsContenar.innerHTML += end("Other users");

  } catch (error) {
    createAlert(alert("alert-danger", "feald to get last users"))
  }
}
