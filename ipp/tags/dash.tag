<dash>
  <style>
h1 {
  text-align: center;
}

material-input {
  width: 100%;
}

#update-form {
  margin-top: 30px;
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

#number-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
  </style>

  <h1>User Dashboard</h1>
  <div id="number-container">
    <p>Current number = { number }</p>
    <form onsubmit={ incrementNumber }>
      <button class="material-button primary"  type="submit">☝️</button>
    </form>
  </div>

  <form id="update-form" onsubmit={ submit }>
    <material-input name="number" label="Number" type="number"></material-input>
    <button class="material-button" type="submit">Update</button>
  </form>
  <span class="error" show={ error }>{ error }</span>

  <script>
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

incrementNumber(e) {
  e.preventDefault()
  xhr("/next", "GET", null, bindNumber, bindError)
}

function bindNumber(res) {
  self.number = res.data.attributes.value;
  self.error = "";
  self.update();
}

submit(e) {
  e.preventDefault();
  updateNumber(e.target.number.value);
}

function bindError(res) {
  self.error = res.errors[0].detail;
  self.update();
  self.parent.update();
}
  </script>
</dash>
