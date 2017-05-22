<login>
  <style>
login {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: "Roboto";
}

.container {
  width: 100%
}

.center-vertical {
  display: flex;
  flex-direction: column;
  align-items: center;
}

material-input {
  width: 100%;
}

.error {
  margin-bottom: 10px;
  align-self: flex-start;
}

.prompt {
  align-self: flex-start;
  margin-bottom: 10px;
}
  </style>

  <h1>i++</h1>
  <div class="container">
    <form class="center-vertical" onsubmit={ submit }>
      <material-input name="email" label="Email" type="email" required />
      <material-input name="password" label="Password" type="password" required />

      <span class="error" show={ loginError }> { loginError } </span>

      <span show={ !login } class="prompt">
        Already have an account? Try <a onclick={ isLogin } href="#login">logging in</a> instead.
      </span>
      <button class="material-button primary" show={ !login } type="submit">Sign up</button>

      <span show={ login } class="prompt">
        Need an account? <a onclick={ isLogin } href="#signup">Sign up</a> now.
      </span>
      <button class="material-button primary" show={ login } type="submit">Login</button>
    </form>

    <div class="center-vertical">
      <p>or</p>
      <a class="material-button facebook" href={ API_URL + "/login/facebook" }>Login with Facebook</a>
    </div>
  </div>

  <script>
self = this;
loginOrSignup();

var err = GetParameterByName("error")
if(err) {
  this.loginError = err;
  self.update();
}

window.onhashchange = loginOrSignup;

submit(e) {
  e.preventDefault();
  auth(self.login ? "/login" : "/signup", e.target.email.value, e.target.password.value);
}

function loginOrSignup() {
  if (location.hash === "#login") {
    self.login = true;
  } else {
    self.login = false;
  }
  self.update();
}

function auth(path, email, password) {
  var payload = JSON.stringify({
    data: {
      type: "user",
      attributes: {
        email: email,
        password: password
      }
    }
  });

  xhr(path, "POST", payload, storeToken, bindLoginError);
}

function storeToken(res) {
  Token = res.data.id;
  StoreToken(Token);
  self.parent.update();
}

function bindLoginError(res) {
  this.loginError = res.errors[0].detail;
  self.update();
}
  </script>
</login>
