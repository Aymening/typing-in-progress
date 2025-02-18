document.getElementById("submit-login")?.addEventListener('click', async (event) => {
	event.preventDefault()
	try {
		let email = document.getElementById("emailL").value
		let password = document.getElementById("passwordL").value


		let response = await fetch('/api/log-in', {
			method: 'POST',
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ emailname: email, password: password })
		})
		let data = await response.json()
		if (response.ok) {
			createAlert(alert("alert-success", data))
			setTimeout(() => {
				checkLogedAndaddPopup()
				toggleMenu('popBackgroundLogin', 'popCreateLogin')
			}, 500)
		} else if (response.status === 401) {
			createAlert(alert("alert-caution", data))
		} else {
			dangerError(response.status)
			return
		}
	} catch (error) {
		createAlert(alert("alert-danger", "feald to login"))
	}
})

function getSelectedGender() {
	const selected = document.querySelector('.mydict input[name="radio"]:checked');
	return selected ? selected.value : null;
}

document.getElementById("submit-register")?.addEventListener('click', async (event) => {
	event.preventDefault()
	try {
		let username = document.getElementById("usernameS").value
		let age = document.getElementById("ageS").value
		let firstname = document.getElementById("firstnameS").value
		let lastname = document.getElementById("lastnameS").value
		let email = document.getElementById("emailS").value
		let password = document.getElementById("paswordS").value
		let gender = getSelectedGender();
		
		let response = await fetch('/api/sign-up', {
			method: 'POST',
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				nickname: username,
				age: parseInt(age),
				firstName: firstname,
				lastName: lastname,
				email: email,
				pasword: password,
				gender: gender,
			})
		})

		let data = await response.json()
		
		if (response.ok) {
			createAlert(alert("alert-success", data))
			setTimeout(() => {
				switchtologin()
			}, 500)
		} else {
			if (response.status === 422) {
				createAlert(alert("alert-caution", data))
			} else if (response.status === 409) {
				createAlert(alert("alert-caution", data))
			} else if (response.status === 401) {
				createAlert(alert("alert-caution", data))
			} else {
				dangerError(response.status)
				return
			}
		}
	} catch (error) {
		createAlert(alert("alert-danger", "feald to signup"))
	}
})


