window.state = null;
window.username = "";

window.checkIfLoggedIn = async function () {
    // if (window.state !== null) {
    //     return { state: window.state, username: window.username };
    // }

    try {
        let response = await fetch(`/api/is-logged-in`, {
            method: 'GET',
            headers: {
                "Content-Type": "application/json",
            }
        });

        if (!response.ok) {
            dangerError(response.status)
            return null;
        }

        let user = await response.json();
        if (user !== null) {
            window.state = user.loggedIn;
            window.username = user.username;
        }else{
            window.state = false
        }

        return { state: window.state, username: window.username };
    } catch (error) {
        console.error("Fetch failed:", error);
        createAlert(alert("alert-danger", "feald to check is loged"))
        return null;
    }
};
