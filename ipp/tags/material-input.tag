<material-input>
  <style>
/* form starting stylings ------------------------------- */
form {
  margin-top:30px;
  font-family:"Roboto";
}

.group {
  position:relative;
  margin-bottom:45px;
}

input {
  width: 100%;
  font-size:18px;
  padding:10px 10px 10px 5px;
  display:block;
  border:none;
  border-bottom:1px solid #757575;
}

input:focus {
  outline:none;
}

/* LABEL ======================================= */
label          {
  color:#999;
  font-size:18px;
  font-weight:normal;
  position:absolute;
  pointer-events:none;
  left:5px;
  top:10px;
  transition:0.2s ease all;
  -moz-transition:0.2s ease all;
  -webkit-transition:0.2s ease all;
}

/* active state */
input:focus ~ label, input:valid ~ label    {
  top:-20px;
  font-size:14px;
  color:#5264AE;
}

/* BOTTOM BARS ================================= */
.bar {
  width: 100%;
  padding: 0px 10px 0px 5px;
  position:relative;
  display:block;
}
.bar:before, .bar:after {
  content:'';
  height:2px;
  width:0;
  bottom:1px;
  position:absolute;
  background:#5264AE;
  transition:0.2s ease all;
  -moz-transition:0.2s ease all;
  -webkit-transition:0.2s ease all;
}
.bar:before {
  left:50%;
}
.bar:after {
  right:50%;
}

/* active state */
input:focus ~ .bar:before, input:focus ~ .bar:after {
  width:50%;
}

/* HIGHLIGHTER ================================== */
input .highlight {
  position:absolute;
  height:60%;
  width:33%;
  top:25%;
  left:0;
  pointer-events:none;
  opacity:0.5;
}

/* active state */
input:focus ~ .highlight {
  -webkit-animation:inputHighlighter 0.3s ease;
  -moz-animation:inputHighlighter 0.3s ease;
  animation:inputHighlighter 0.3s ease;
}

    /* ANIMATIONS ================ */
    @-webkit-keyframes inputHighlighter {
      from { background:#5264AE; }
      to  { width:0; background:transparent; }
    }
    @-moz-keyframes inputHighlighter {
      from { background:#5264AE; }
      to  { width:0; background:transparent; }
    }
    @keyframes inputHighlighter {
      from { background:#5264AE; }
      to  { width:0; background:transparent; }
    }
  </style>


    <div class="group">
      <input name={ opts.name } type={ opts.type } onkeyup={ opts.onkeyup } required>
      <span class="highlight"></span>
      <span class="bar"></span>
      <label>{ opts.label }</label>
    </div>
</material-input>
