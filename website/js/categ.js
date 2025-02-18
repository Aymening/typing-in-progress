async function categ() {
    try {
        let res = await fetch('/api/categ')

        if (!res.ok) {
            dangerError(res.status)
            return
        }

        let str = await res.json()
        let div = ""
        for (let i = 0; i < str.length; i++) {
            div += `<div class="Fcateg" data-id="${str[i].category}">${str[i].category}</div>`
        }
        document.getElementById("filterCategContent").innerHTML = div
        let catelm = document.querySelectorAll('.Fcateg')

        for (let i = 0; i < catelm.length; i++) {
            catelm[i].addEventListener('click', function () {
                let category = this.dataset.id
                handleCategoryClick(category)
            })
        }

    } catch (error) {
        console.error("error category :", error)
      //  createAlert(alert("alert-danger", "feald to get category"))
    }
}

function handleCategoryClick(catego) {
    document.getElementById('centerHome').innerHTML = ""
    display_post(1, "filterPost", { category: catego })
}

// function showCateg() {
//     let urlParams = new URLSearchParams(window.location.search);
//     let categ = urlParams.get('categ') || "";
//     let date = urlParams.get('date') || "";
//     let like = urlParams.get('like') || "";

//     if (categ !== "") {
//         document.getElementById('nameCateg').innerHTML = categ
//     }

//     if (date == "DESC") {
//         document.getElementById("new").setAttribute('selected', '')
//     } else if (date == "ASC") {
//         document.getElementById("old").setAttribute('selected', '')
//     } else if (like == "DESC") {
//         document.getElementById("mostLiked").setAttribute('selected', '')
//     } else if (like == "ASC") {
//         document.getElementById("leastLiked").setAttribute('selected', '')
//     }
// }

// showCateg();
