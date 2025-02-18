async function postLike(PostID, typeLike) {
    const { state, _ } = await window.checkIfLoggedIn();
    if (!state){
        createAlert(alert("alert-danger", "you need to log in"))
        checkLogedAndaddPopup()
        return
    }

    let statusLike = 0
    let iLike = document.getElementById('pid' + PostID + 'l') 
    let idisLike = document.getElementById('pid' + PostID + 'd')
    let spanLike = document.getElementById('pid' + PostID + 'ls')
    let spanDislike = document.getElementById('pid' + PostID + 'ds')
    if (typeLike == "liked") {
        if (iLike.classList.contains('liked')) {
            iLike.classList.remove('liked')
            spanLike.textContent = parseInt(spanLike.textContent) - 1
            statusLike = 0
        } else {
            if (idisLike.classList.contains('disliked')) {
                idisLike.classList.remove('disliked')
                spanDislike.textContent = parseInt(spanDislike.textContent) - 1
            }
            iLike.classList.add('liked')
            spanLike.textContent = parseInt(spanLike.textContent) + 1 
            statusLike = 1
        }
    } else if (typeLike == "disliked") {
        if (idisLike.classList.contains('disliked')) {
            idisLike.classList.remove('disliked')
            spanDislike.textContent = parseInt(spanDislike.textContent) - 1
            statusLike = 0
        } else {
            if (iLike.classList.contains('liked')) {
                iLike.classList.remove('liked')
                spanLike.textContent = parseInt(spanLike.textContent) - 1
            }
            idisLike.classList.add('disliked')
            spanDislike.textContent = parseInt(spanDislike.textContent) + 1
            statusLike = -1
        }
    } else {
        createAlert(alert("alert-danger", "feald to get post like"))
        return
    }
    try {
        const response = await fetch(`/api/reactions?item_id=${PostID}&status_like=${statusLike}`)

        if (response.status == 401) {
            checkLogedAndaddPopup()
            return
        }
        if (!response.ok) {
            console.log("error")
            dangerError(response.status)
            return
        }

    } catch (error) {
        console.error("Fetch failed:", error);
        createAlert(alert("alert-danger", "feald to get post like"))
    }
}