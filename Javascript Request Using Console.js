//Start of Getting Users
//When you want to how many user do you have in your database
fetch('http://localhost:1234/api/v1/user')
  .then(response => response.json())
  .then(data => console.log(data));
//End of Getting Users

//Start of selecting all from datase where a specific column is specify
function forLoop(data){
	for(var i = 0; i<data.length;++i){
console.log(data[i].username);
	}
}

fetch('http://localhost:1234/api/v1/user')
  .then(response => response.json())
  .then(data =>forLoop(data));
//Start of selecting all from datase where a specific column is specify

//Start of Adding New User
fetch('http://localhost:1234/api/v1/user/', {
  method: 'POST',
  mode: 'cors',
  headers: new Headers({
    'Content-Type': 'application/json'
  }), body: JSON.stringify({ username: "john", password: "wawa", name: "John" })
}).then(response => response.json()).then(data => console.log(data))
//End of Adding New User


//Start of Login Authentication
fetch("/api/v1/user/auth/", {
  method: 'POST',
  mode: 'cors',
  headers: new Headers({
    'Content-Type': 'application/json'
  }),
  body: JSON.stringify({
    'username': 'john',
    'password': 'wawa'
  })
}).then(r => r.json()).then(d => console.log(d.token))
//End of Login Authentication


//Start of No Input Username or Password
fetch("/api/v1/user/auth/", {
  method: 'POST',
  mode: 'cors',
  headers: new Headers({
    'Content-Type': 'application/json'
  }),
  body: JSON.stringify({
    'username': '',
    'password': ''
  })
}).then(r => r.json()).then(d => console.log(d.message))
//End of No Input Username or Password


//Start of Getting IP and Generating URL (must be login first before using this)
fetch("http://localhost:1234/api/v1/ip/google.com", {
  mode: 'cors',
  headers: new Headers({
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <Token>'
  })
}).then(r => r.json()).then(d => console.log(d))
// End of Getting IP and Generating URL


//Start of Update//
fetch("/api/v1/user/", {
  method: 'PUT',
  mode: 'cors',
  headers: new Headers({
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <Token>'
  }),
  body: JSON.stringify({
    'id': 5,
    'username': 'john1',
    'password': 'wawa',
    'name': ''
  })
}).then(r => r.json()).then(d => console.log(d))
//End of Update


//Start of Delete
fetch("/api/v1/user/5", {
  method: 'DELETE',
  mode: 'cors',
  headers: new Headers({
    'Content-Type': 'application/json',
    'Authorization': 'Bearer <Token>'
  })
}).then(r => r.json()).then(d => console.log(d))
//End of Delete

//Update February 13, 2021//