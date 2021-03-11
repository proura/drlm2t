import AppHeader from '/js/components/AppHeader.js'
import AppBody from '/js/components/AppBody.js'
import AppFooter from '/js/components/Appfooter.js';


const app = new Vue({
  el: '#app',
  components: {
    'app-header': AppHeader,
    'app-body': AppBody,
    'app-footer': AppFooter
  },
  data() {
    return {
      context: {
        loggedin: false,
        route: "/",
      }
    }
  },
  created() {
    this.checkLogin()
  },
  methods: {
    checkLogin: function (){
      fetch('/checkSession',{method: 'POST'})
      .then(response => {
        if (response.ok){
          this.context.loggedin = true;
        } else {
          this.context.loggedin = false;
        }
      });
    }
  },
});
