var AppHeader = {
  props: ['context'],
  template: `
    <div v-if="context.loggedin" class="row">
      <div class="col">
        <img class="mb-4" src="static/drlm_long.png" alt="" width="250" height="131">
      </div>
      <div class="col">
          <button type="button" class="btn btn-primary float-end" v-on:click="appLoggout" id="logout-button">Logout</button>
      </div>
    </div>
  `,
  methods: {
    appLoggout: function(event){

      fetch('/logout',{method: 'POST'})
      .then(response => {
        if (response.ok){
          this.context.loggedin = false;
        }
      });

      event.preventDefault();
    }
  }
}

export default AppHeader;
