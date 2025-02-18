async function displayCategory() {
	const checkbx = document.getElementById('checkboxes');
	if (!checkbx) {
		console.error("Checkbox container not found!");
		return;
	}

	try {
		const response = await fetch('/api/categ');
		if (!response.ok) {
			dangerError(response.status)
			throw new Error('Failed to fetch categories');
		}

		const data = await response.json();
		if (data === null) {
			createAlert(alert("alert-caution", "nout found any category"))
			return
		}
		let checkboxHTML = '';
		for (const val of data) {
			checkboxHTML += `<label>
								<input type="checkbox" name="${val.category}" value="${val.category}" /> ${val.category}
							  </label>`;
		}
		checkbx.innerHTML = checkboxHTML;
	} catch (err) {
		createAlert(alert("alert-danger", "feald desplay category"))
		console.error('Error:', err);
	}
}


//   send request to api creat post and check response
	document.getElementById('CreatePostForm')?.addEventListener('submit', async (event) => {
		event.preventDefault()

		const title = document.getElementById("TitleCreatePost");
		const content = document.getElementById("ContentCreatePost");
		const checkboxes = document.querySelectorAll('#checkboxes input[type="checkbox"]');
		const image = document.getElementById("Image");

		let category = []
		checkboxes.forEach((item) => {
			if (item.checked) {
				category.push(item.value)
			}
		});
		if (title.value === "" && content.value === "") {
			createAlert(alert("alert-caution", "you can't create empty post"))
			return
		}

		const formData = new FormData();

		formData.append('title', title.value);
		formData.append('content', content.value);
		formData.append('category', JSON.stringify(category));
		formData.append('image', image.files[0]);

		try {
			const respFetch = await fetch("/api/create-post", {
				method: 'POST',
				body: formData
			})
			let response = await respFetch.json()

			// check error status code
			if (!respFetch.ok) {
				dangerError(respFetch.status)
				throw new Error(response)
			}

			// if all work like good 
			title.value = "";
			content.value = "";
			image.value = "";
			checkboxes.forEach((item) => item.checked = false);
			document.getElementById('checkboxes').classList.add("none")
			toggleMenu('popBackground', 'popCreate');
			createAlert(alert("alert-success", response));

		} catch (err) {
			createAlert(alert("alert-danger", "feld to creat post"))
			console.error(err)
		}
	})

async function logout() {
	try {
		window.state = null;
		window.username = "";

		const resp = await fetch("/api/logout");

		const valResp = await resp.json();
		if (!resp.ok) {
			dangerError(resp.status)
			throw new Error('Failed to fetch categories');
		}

		createAlert(alert("alert-success", valResp));
		closeWS()
		isMoerUsers = true
		
		setTimeout(() => {
			checkLogedAndaddPopup()
		}, 500);

	} catch (err) {
		createAlert(alert("alert-danger", "feald to logout"))
		console.error(err);
	}
}
