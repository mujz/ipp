riot.tag2('app', '<button if="{Token}" id="logout" class="material-button" type="button" onclick="{logout}">Logout</button> <login if="{!Token}"></login> <dash if="{Token}"></dash>', 'app,[data-is="app"]{ width: 40%; max-width: 500px; } @media (max-width: 800px) { app,[data-is="app"]{ width: 60%; } } @media (max-width: 550px) { app,[data-is="app"]{ width: 90%; } } app #logout,[data-is="app"] #logout{ position: absolute; top: 0; right: 16px; }', '', function(opts) {
token = GetParameterByName("token")
if(token) {
  StoreToken(token);
  this.update();
  window.location.href = window.location.origin;
}

this.logout = function(e) {
  ClearToken();
  this.update();
}.bind(this)
});

riot.tag2('dash', '<h1>User Dashboard</h1> <div id="number-container"> <p>Current number = {number}</p> <form onsubmit="{incrementNumber}"> <button class="material-button primary" type="submit">☝️</button> </form> </div> <form id="update-form" onsubmit="{submit}"> <material-input name="number" label="Number" type="number"></material-input> <button class="material-button" type="submit">Update</button> </form> <span class="error" show="{error}">{error}</span>', 'dash h1,[data-is="dash"] h1{ text-align: center; } dash material-input,[data-is="dash"] material-input{ width: 100%; } dash #update-form,[data-is="dash"] #update-form{ margin-top: 30px; display: flex; justify-content: space-between; align-items: baseline; } dash #number-container,[data-is="dash"] #number-container{ display: flex; justify-content: space-between; align-items: center; }', '', function(opts) {
var self = this;

xhr("/current", "GET", null, bindNumber, bindError);

function updateNumber(newVal) {
  var payload = '{' +
    '"data": {' +
      '"type": "number",' +
      '"attributes": {' +
        '"value": ' + newVal +
      '}' +
    '}' +
  '}';
  xhr("/current", "PUT", payload, bindNumber, bindError);
}

this.incrementNumber = function(e) {
  e.preventDefault()
  xhr("/next", "GET", null, bindNumber, bindError)
}.bind(this)

function bindNumber(res) {
  self.number = res.data.attributes.value;
  self.error = "";
  self.update();
}

this.submit = function(e) {
  e.preventDefault();
  updateNumber(e.target.number.value);
}.bind(this)

function bindError(res) {
  self.error = res.errors[0].detail;
  self.update();
  self.parent.update();
}
});

riot.tag2('login', '<h1>i++</h1> <div class="container"> <form class="center-vertical" onsubmit="{submit}"> <material-input name="email" label="Email" required type="email"></material-input> <material-input name="password" label="Password" type="password" required></material-input> <span class="error" show="{loginError}"> {loginError} </span> <span show="{!login}" class="prompt"> Already have an account? Try <a onclick="{isLogin}" href="#login">logging in</a> instead. </span> <button class="material-button primary" show="{!login}" type="submit">Sign up</button> <span show="{login}" class="prompt"> Need an account? <a onclick="{isLogin}" href="#signup">Sign up</a> now. </span> <button class="material-button primary" show="{login}" type="submit">Login</button> </form> <div class="center-vertical"> <p>or</p> <a class="material-button facebook" href="{API_URL + ⁗/login/facebook⁗}">Login with Facebook</a> </div> </div>', 'login { display: flex; flex-direction: column; align-items: center; font-family: "Roboto"; } login .container,[data-is="login"] .container{ width: 100% } login .center-vertical,[data-is="login"] .center-vertical{ display: flex; flex-direction: column; align-items: center; } login material-input,[data-is="login"] material-input{ width: 100%; } login .error,[data-is="login"] .error{ margin-bottom: 10px; align-self: flex-start; } login .prompt,[data-is="login"] .prompt{ align-self: flex-start; margin-bottom: 10px; }', '', function(opts) {
self = this;
loginOrSignup();

var err = GetParameterByName("error")
if(err) {
  this.loginError = err;
  self.update();
}

window.onhashchange = loginOrSignup;

this.submit = function(e) {
  e.preventDefault();
  auth(self.login ? "/login" : "/signup", e.target.email.value, e.target.password.value);
}.bind(this)

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
});

riot.tag2('material-input', '<div class="group"> <input name="{opts.name}" type="{opts.type}" onkeyup="{opts.onkeyup}" required> <span class="highlight"></span> <span class="bar"></span> <label>{opts.label}</label> </div>', 'material-input form,[data-is="material-input"] form{ margin-top:30px; font-family:"Roboto"; } material-input .group,[data-is="material-input"] .group{ position:relative; margin-bottom:45px; } material-input input,[data-is="material-input"] input{ width: 100%; font-size:18px; padding:10px 10px 10px 5px; display:block; border:none; border-bottom:1px solid #757575; } material-input input:focus,[data-is="material-input"] input:focus{ outline:none; } material-input label,[data-is="material-input"] label{ color:#999; font-size:18px; font-weight:normal; position:absolute; pointer-events:none; left:5px; top:10px; transition:0.2s ease all; -moz-transition:0.2s ease all; -webkit-transition:0.2s ease all; } material-input input:focus ~ label,[data-is="material-input"] input:focus ~ label,material-input input:valid ~ label,[data-is="material-input"] input:valid ~ label{ top:-20px; font-size:14px; color:#5264AE; } material-input .bar,[data-is="material-input"] .bar{ width: 100%; padding: 0px 10px 0px 5px; position:relative; display:block; } material-input .bar:before,[data-is="material-input"] .bar:before,material-input .bar:after,[data-is="material-input"] .bar:after{ content:\'\'; height:2px; width:0; bottom:1px; position:absolute; background:#5264AE; transition:0.2s ease all; -moz-transition:0.2s ease all; -webkit-transition:0.2s ease all; } material-input .bar:before,[data-is="material-input"] .bar:before{ left:50%; } material-input .bar:after,[data-is="material-input"] .bar:after{ right:50%; } material-input input:focus ~ .bar:before,[data-is="material-input"] input:focus ~ .bar:before,material-input input:focus ~ .bar:after,[data-is="material-input"] input:focus ~ .bar:after{ width:50%; } material-input input .highlight,[data-is="material-input"] input .highlight{ position:absolute; height:60%; width:33%; top:25%; left:0; pointer-events:none; opacity:0.5; } material-input input:focus ~ .highlight,[data-is="material-input"] input:focus ~ .highlight{ -webkit-animation:inputHighlighter 0.3s ease; -moz-animation:inputHighlighter 0.3s ease; animation:inputHighlighter 0.3s ease; } @-webkit-keyframes inputHighlighter { from { background:#5264AE; } to { width:0; background:transparent; } } @-moz-keyframes inputHighlighter { from { background:#5264AE; } to { width:0; background:transparent; } } @keyframes inputHighlighter { from { background:#5264AE; } to { width:0; background:transparent; } }', '', function(opts) {
});
