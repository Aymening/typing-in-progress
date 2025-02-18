// part yassine
let register = document.getElementsByClassName("register");
let log_in = document.getElementsByClassName("log-in");
let to_register = document.getElementById("to_register");
let to_log_in = document.getElementById("to_log_in");
const body = document.querySelector("body");
const topImg = document.getElementById("topImg")
const createPost = document.getElementById("createPost")

if (to_log_in) {
  to_log_in.addEventListener("click", () => {
    switchtologin();
  });
}

if (to_register) {
  to_register.addEventListener("click", () => {
    log_in[0].style.display = "none";
    log_in[1].style.display = "none";
    register[0].style.display = "flex";
    register[1].style.display = "flex";
  });
}

function switchtologin() {
  register[0].style.display = "none";
  register[1].style.display = "none";
  log_in[0].style.display = "flex";
  log_in[1].style.display = "flex";
}
// end

function showCheckboxes(c) {
  document.getElementById(c).classList.toggle("none");
}

function toggleCateg() {
  document.getElementById("uparrow").classList.toggle("none");
  document.getElementById("downarrow").classList.toggle("none");
  document.getElementById("filterCategContent").classList.toggle("none");
}

function toggleMenu(c, bc, none) {
  if (none === "none") {
    document.getElementById(c).classList.remove("none");
    document.getElementById(bc).classList.remove("none");
  } else {
    document.getElementById(c).classList.toggle("none");
    document.getElementById(bc).classList.toggle("none");
  }
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

//---------------------------------- alert ----------------------
function createAlert(passIntAlert) {
  let alertElernt = document.createElement("alert");
  alertElernt.innerHTML = passIntAlert;
  body.appendChild(alertElernt);
}

function removeAlert() {
  const newAlert = document.querySelector("alert");
  newAlert?.remove();
}

const alert = (clas, message) => {
  // auto remove
  removeAlert();
  setTimeout(() => {
    const newAlert = document.querySelector("alert");
    if (newAlert) {
      removeAlert();
    }
  }, 3000);

  return `<div class="alert">
      <div class="alert-container">
        <div class="alert-content ${clas}">
          <div class="alert-header">
            <span class="alert-message">${message}</span>
            <span class="btnAlert btn-close" onclick="removeAlert()">x</span>
          </div>
        </div>
      </div>
    </div>`;
};
//-------------------------------- mode --------------------------------------------
const root = document.querySelector(':root');
const setVariables = vars => Object.entries(vars).forEach(v => root.style.setProperty(v[0], v[1]));
const lightMode = {
  '--background': '#f7f7f7',
  '--primary-color': '#2c3e50',
  '--secondary-color': '#6b7272',
  '--secondary-color-hover': '#7f8c8d',
  '--third-color': '#e4e4e4',
  '--success-green': '#27ae60',
  '--notification-color': '#3498db',
  '--notification-color-hover': '#2980b9',
  '--warning-color': '#e67e22',
  '--warning-color-hover': '#d35400',
};
const darkMode = {
  '--background': '#1a1f25',
  '--primary-color': '#e0e6ed',
  '--secondary-color': '#a0a9b5',
  '--secondary-color-hover': '#8b94a0',
  '--third-color': '#2b323c',
  '--success-green': '#4caf50',
  '--notification-color': '#2196f3',
  '--notification-color-hover': '#1976d2',
  '--warning-color': '#f44336',
  '--warning-color-hover': '#d32f2f',
};
function setModeColor(btn) {
  let mode = btn.dataset.mode
  if(mode === "dark"){
    btn.classList = "fa-solid fa-sun"
    btn.dataset.mode = "light"
    setVariables(lightMode)
  }else{
    btn.classList = "fa-regular fa-moon"
    btn.dataset.mode = "dark"
    setVariables(darkMode)
  }
}

// ----------------------------- errror --------------------------------

function dangerError(status) {
  let err = "";
  if (status === 500) {
    err = errTemplate(500, "Internal Server Error");
  } else if (status === 400) {
    err = errTemplate(400, "Bad Request ");
  }
  if (err != "") document.body.innerHTML = err;
}

const errTemplate = (code, message) => {
  return `<div class="body-error">
    <h1>Error ${code}</h1>
    <p>${message}</p>
  </div>`;
};

//-------------------
function returnToHome() {
  document.getElementById("centerHome").innerHTML = "";
  display_post(1, "filterPost");
}

//--------------------------------------------- handel path
let msgInput = document.getElementById("messageInput");
function navigate(path) {
  if (path === "/") {
    categ();
    checkLogedAndaddPopup();
    displayCategory();
  } else {
    document.body.innerHTML = errTemplate(404, "page not found");
  }
}

window.addEventListener("popstate", () => {
  navigate(window.location.pathname);
});

navigate(window.location.pathname);


// check is user log in
async function checkLogedAndaddPopup() {
  const { state, username } = await window.checkIfLoggedIn()
  let pname = document.querySelector('.username')
  if (state) {
    returnToHome()
    pname.textContent = username
    await lastUsersChat()
    await getusers(true)
    msessagesWS()
    topImg.classList.remove("none")
    createPost.classList.remove("none")
  } else {
    pname.textContent = ""
    document.getElementById("centerHome").innerHTML = ""
    document.getElementById("users_contenar").innerHTML = ""
    topImg.classList.add("none")
    createPost.classList.add("none")
    CloseChat()
    toggleMenu('popBackgroundLogin', 'popCreateLogin', 'none')
  }
}
