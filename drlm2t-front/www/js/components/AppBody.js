import InfrastructuresTable from '/js/components/InfrastructuresTable.js'

var AppBody = {
  components: {
    'infrastructures-table': InfrastructuresTable,
  },
  props: ['context'],
  template: `
    <div v-if="context.loggedin">
      <infrastructures-table></infrastructures-table>
    </div>

    <div v-else class="text-center"> 
      <main class="form-signin">
      <form>
        <img class="mb-4" src="static/drlm_long.png" alt="" width="250" height="131">
        <h5 class="h5 mb-3 fw-normal">Please sign in</h5>
        <label for="inputUser" class="visually-hidden">Username</label>
        <input type="text" id="inputUser" class="form-control" placeholder="Username" required autofocus>
        <label for="inputPassword" class="visually-hidden">Password</label>
        <input type="password" id="inputPassword" class="form-control" placeholder="Password" required>
        <button class="w-100 btn btn-lg btn-primary" type="submit" v-on:click="appLoggin">Sign in</button>
        <div class="alert alert-success alert-dismissible" id="success" style="display:none;">
          <a href="#" class="close" data-dismiss="alert" aria-label="close">×</a>
        </div>
        <div class="alert alert-danger alert-dismissible" id="error" style="display:none;">
          <a href="#" class="close" data-dismiss="alert" aria-label="close">×</a>
        </div>
      </form>
      </main>
    </div>
  `,
  methods: {
    appLoggin: function(event){
      var formData = {
        'username'  : $('input[id=inputUser]').val(),
        'password'  : $('input[id=inputPassword]').val()
      };

      const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData)
      }

      fetch('/signin', requestOptions)
      .then(response => {
        if(response.ok) {
          $("#error").hide();
          $("#success").show();
          $('#success').html('Loggin OK!');
          this.context.loggedin = true;
        } else {
          $("#success").hide();
          $("#error").show();
          $('#error').html('Invalid Username or Password!');
          this.context.loggedin = false;
        }
      });

      event.preventDefault();
    },
  }
}

export default AppBody;
