//Start of Getting Users
fetch('http://localhost:1234/api/v1/user')
  .then(response => response.json())
  .then(data => console.log(data));
//End of Getting Users

//Start of Adding New User
fetch('http://localhost:1234/api/v1/user/', {
  method: 'POST',
  mode: 'cors',
  headers: new Headers({
    'Content-Type': 'application/json'
  }), body: JSON.stringify({ username: "john", password: "wawa" })
}).then(response => response.json()).then(data => console.log(data.username))
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
    'password': 'wawa'
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



