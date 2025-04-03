async function signup() {
    /*
        Triggers on submit of signup button
        1.) Retrieves email and password from web elements, jsonify
        2.) Sends a post request on the register API, passing the email and password 
        as JSON
        3.) User is alerted of any errors
        4.) If successful, response body:
            {
                User : {
                            id: ${userID},
                            created_at : ${date created},
                            updated_at : ${date updated},
                            email : ${email}
                        }
            }
    */
    const email = document.getElementById("reg-email").value;
    const password = document.getElementById("reg-password").value;
    jsonObj = {
        email: email,
        password: password
    }
    try{
        const res = await fetch("/api/users", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(jsonObj),
        });
        if (!res.ok){
            const data = await res.json();
            throw new Error(`Failed to create user: ${data.error}`);
        }
        console.log("User created!");      
    }catch(error){
        alert(`Error: ${error.message}`);
    }
}

async function login() {
    /*
        Triggers on submit of login button
        1.) Retrieves email and password from web elements, jsonify
        2.) Sends a post request on the register API, passing the email and password 
        as JSON
        3.) User is alerted of any errors
        4.) If successful, response body:
            {
                User : {
                            id: ${userID},
                            created_at : ${date created},
                            updated_at : ${date updated},
                            email : ${email}
                        }
                token : ${accessToken},
                refresh_token : ${refreshToken}
            }
    */
    const email = document.getElementById("reg-email").value;
    const password = document.getElementById("reg-password").value;
    jsonObj = {
        email: email,
        password: password
    }
    try {
        const res = await fetch("/api/login", {
            method: 'POST',
            headers : {'Content-Type': 'application/json'},
            body: json.stringify(jsonObj),
        });
        const data = await res.json();
        if (!res.ok) {
            throw new Error(`Error logging in: ${data.error}`);
        }
        if (data.token){
            sessionStorage.setItem('token', data.token); 
        }else{
            alert('Login failed. Please check your credentials.');
        }
        
    }catch{
        alert(`Error: ${error.message}`);
    }

}       