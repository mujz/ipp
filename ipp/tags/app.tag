<app>
  <style>
:scope {
  width: 40%;
  max-width: 500px;
}

@media (max-width: 800px) {
  :scope {
    width: 60%;
  }
}

@media (max-width: 550px) {
  :scope {
    width: 90%;
  }
}

#logout {
  position: absolute;
  top: 0;
  right: 16px;
}
  </style>
  <button if={ Token } id="logout" class="material-button" type="button" onclick={ logout }>Logout</button>
  <login if={ !Token }/>
  <dash if={ Token }/>

  <script>
token = GetParameterByName("token")
if(token) {
  StoreToken(token);
  this.update();
  window.location.href = window.location.origin;
}

logout(e) {
  ClearToken();
  this.update();
}
  </script>
</app>
