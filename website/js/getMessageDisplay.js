let isMoreMessages = true;
let lastUserClick = "";
let userOpenChat = null
let lastSelect = null;


const chat = document.getElementById("chat");
const sendBtn = document.getElementById("sendBtn");
const container = document.querySelector(".container");
const chatBtnInput = document.getElementById("chatbtninput");

async function getMessages(user, newuser) {
    if (!isMoreMessages) return;

    let lastMsg = "";
    const msgNameText = document.querySelectorAll(".msgNameText");
    const sendTo = sendBtn.getAttribute("sendto");

    if (!user) {
        user = sendTo;
    }

    if (msgNameText.length > 0 && !newuser) {
        lastMsg = msgNameText[0].children[2].dataset.date;
    }
    try {
        const response = await fetch("/api/getMessages", {
            method: "POST",
            body: JSON.stringify({ chatuser: user, lastMessage: lastMsg }),
        });

        if (response.status === 412) {
            isMoreMessages = false;
            return null;
        }

        if (!response.ok) {
            dangerError(response.status);
            throw new Error("Failed to fetch messages");
        }

        const data = await response.json();

        if (!data || data.length === 0) {
            isMoreMessages = false;
            return null;
        }
        return data;
    } catch (err) {
        createAlert("Failed to get messages");
        return null;
    }
}

function chattemplate(msg, right = "") {
    return `
        <div class="msgNameText ${right}">
            <p>${msg.send}</p>
            <h3>${msg.Message}</h3>
            <p data-date="${msg.timeSend}">${new Date(msg.timeSend).toLocaleString()}</p>
        </div>
    `;
}
function Addchattemplate(msg, right = "") {
    return `
        <div class="msgNameText ${right}">
            <p>${msg.sender}</p>
            <h3>${msg.text}</h3>
            <p data-date="${msg.timestamp}">${new Date(msg.timestamp).toLocaleString()}</p>
        </div>
    `;
}

function addMessageInChat(msg) {
    chat.innerHTML += Addchattemplate(msg)
    chat.scrollTop = chat.scrollHeight - chat.clientHeight;
}

function addNotifNumber(neckname) {
    let lbl = document.getElementById(`label_${neckname}`)
    lbl.textContent++
}

function addMyMessage(textmsg) {
    let tmpMsg = chattemplate({ send: window.username, Message: textmsg, timeSend: new Date() }, "directionRight")
    chat.innerHTML += tmpMsg
    chat.scrollTop = chat.scrollHeight - chat.clientHeight;
}

function CloseChat() {
    document.getElementById(lastSelect)?.classList.remove('selectUser')
    lastSelect = null
    lastUserClick = ""
    userOpenChat = null
    sendBtn.setAttribute("sendto", "");
    container.classList.remove("openchat");
    chat.classList.add("none");
    chatBtnInput.classList.add("none");
}

function openchat(data, isNewUser = false) {
    if (isNewUser) {
        chat.innerHTML = "";
    }
    container.classList.add("openchat");
    chat.classList.remove("none");
    chatBtnInput.classList.remove("none");

    if (data === "first chat" || data === "no more message") {
        chat.innerHTML = end(data) + chat.innerHTML;
        return;
    }

    let innerChat = "";
    for (let msg of data.reverse()) {
        const isCurrentUser = msg.send === window.username;
        innerChat += chattemplate(msg, isCurrentUser ? "directionRight" : "");
    }

    chat.insertAdjacentHTML("afterbegin", innerChat);
    chat.scrollTop = chat.scrollHeight - chat.clientHeight;
}

async function clickusers(userElement) {    
    if (lastSelect) {
        document.getElementById(lastSelect).classList.remove('selectUser')
    }
    userElement.classList.add("selectUser")
    lastSelect = userElement.id

    const name = userElement.querySelector("h3").textContent;
    userElement.querySelector("label").textContent = ""
    userOpenChat = name

    if (name === lastUserClick) return;
    lastUserClick = name;

    sendBtn.setAttribute("sendto", name);
    isMoreMessages = true;

    let data = await getMessages(name, true);
    if (!data) data = "first chat";


    openchat(data, true);
}

const chatScroll = debounce(async () => {
    if (!isMoreMessages) return;

    const { scrollTop, scrollHeight, clientHeight } = chat;

    if (scrollTop <= 50 && (scrollHeight !== clientHeight)) {
        const scrollOffset = chat.scrollHeight - chat.scrollTop;

        let data = await getMessages();
        if (!data) data = "no more message";
        openchat(data);

        chat.scrollTop = chat.scrollHeight - scrollOffset;
    }
}, 200);

chat.addEventListener("scroll", chatScroll);